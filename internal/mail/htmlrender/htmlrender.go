// Package htmlrender handles rendering HTML templates for emails
package htmlrender

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

//go:embed templates/*
var templateFiles embed.FS

const (
	welcomeTemplate       = "welcome.html"
	passwordResetTemplate = "password-reset.html"
)

// WelcomeData represents the data needed to render the welcome template
type WelcomeData struct {
	Name string
}

// RenderWelcomeTemplate renders the welcome template
func RenderWelcomeTemplate(data WelcomeData) (string, error) {
	return render(welcomeTemplate, data)
}

// PasswordResetData represents the data needed to render the password reset template
type PasswordResetData struct {
	Name     string
	ResetURL string
}

// RenderPasswordResetTemplate renders the password reset template
func RenderPasswordResetTemplate(data PasswordResetData) (string, error) {
	return render(passwordResetTemplate, data)
}

// render renders a template with the given data
func render(name string, data any) (string, error) {
	// Prepend templates/ to the template name
	name = fmt.Sprintf("templates/%s", name)

	// Parse the template
	tmpl, err := template.ParseFS(templateFiles, name)
	if err != nil {
		return "", fmt.Errorf("parsing template %s: %w", name, err)
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template %s: %w", name, err)
	}
	return buf.String(), nil
}
