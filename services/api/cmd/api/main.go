package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/andrew-hayworth22/critiquefi/services/api/internal/business/authbus"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/business/filmbus"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/business/sysbus"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/config"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/http/authhttp"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/http/filmhttp"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/http/syshttp"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/mail"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/mail/authmail"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/mail/senders/logmail"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/mail/senders/ses"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/middleware"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/server"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/store/postgres"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/store/postgres/authpg"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/store/postgres/filmpg"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/store/postgres/syspg"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	ctx := context.Background()

	// Build configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// Setup logging
	logger := cfg.NewLogger()
	logger.Info("starting server", "ENV", cfg.Env, "PORT", cfg.Port, "LOG_LEVEL", cfg.LogLevel)

	// Establish database connection
	db, err := postgres.NewDB(ctx, postgres.DBConfig{
		URL:             cfg.DatabaseURL,
		MaxOpenConns:    cfg.MaxOpenDBConns,
		MaxIdleConns:    cfg.MaxIdleDBConns,
		ConnMaxLifetime: cfg.DBConnMaxLifetime,
		ConnMaxIdleTime: cfg.DBConnMaxIdleTime,
	})
	if err != nil {
		logger.Error("failed connecting to database", "error", err)
		return
	}

	// Run migrations
	m, err := migrate.New("file://migrations", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("error connecting to database for migrations: %v", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error("failed running migrations", "error", err)
		return
	}

	// Build mail sender
	var sender mail.Sender
	switch strings.ToLower(cfg.EmailProvider) {
	case "ses":
		sender, err = ses.New(ctx, ses.Config{
			Region:      "",
			FromAddress: "",
			FromName:    "",
		})
		if err != nil {
			logger.Error("failed to create SES email provied", "error", err)
		}
	case "logmail":
		sender = logmail.New(logger)
	default:
		logger.Error("invalid email provider", "provider", cfg.EmailProvider)
		return
	}

	// Build storage packages
	sysStore := syspg.New(db)
	authStore := authpg.New(db)
	filmStore := filmpg.New(db)

	// Build mail packages
	authMailer := authmail.New(sender, cfg.BaseURL)

	// Build business logic packages
	sysBus := sysbus.New(sysStore)
	authBus := authbus.New(authbus.BusConfig{
		Logger:                logger,
		Store:                 authStore,
		Mailer:                authMailer,
		AccessTokenKey:        cfg.JWTSecret,
		AccessTokenTTL:        cfg.AccessTokenTTL,
		RefreshTokenTTL:       cfg.RefreshTokenTTL,
		PasswordResetTokenTTL: cfg.PasswordResetTokenTTL,
	})
	filmBus := filmbus.New(filmStore)

	// Build handler packages
	sysHandler := syshttp.New(sysBus)
	authHandler := authhttp.New(authBus, cfg.RefreshTokenCookieName, cfg.RefreshTokenCookieDomain)
	filmHandler := filmhttp.New(filmBus)

	// Build middleware packages
	authMiddleware := middleware.NewAuthMiddleware(authBus)

	// Build router
	dependencies := server.Dependencies{
		Logger:         logger,
		AuthHandler:    authHandler,
		SysHandler:     sysHandler,
		FilmHandler:    filmHandler,
		AuthMiddleware: authMiddleware,
	}

	srv := server.New(server.Config{
		Addr:         cfg.Port,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}, dependencies)

	// Start server
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("error starting server: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for an interrupt signal to gracefully shut down the server with
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown server after 30 seconds
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("error shutting down server: %v", err)
		os.Exit(1)
	}
}
