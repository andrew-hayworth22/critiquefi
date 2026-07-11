build:
	docker build \
	  --build-arg VERSION=1.0.0 \
	  --build-arg COMMIT=$(git rev-parse --short HEAD) \
	  -t critiquefi-service:latest \
	  -t critiquefi-service:1.0.0 .

run:
	docker stop critiquefi-service && docker rm critiquefi-service && docker run -d --env-file .env -p 8080:8080 --name critiquefi-service critiquefi-service

build-run: build run

genkey:
	go run ./cmd/genkey/main.go

migrate-up:
	go run ./cmd/migrate -action=up

migrate-down:
	go run ./cmd/migrate -action=down

migrate-fresh:
	go run ./cmd/migrate -action=fresh

test:
	go test ./...

test-store:
	go test ./internal/store/...

test-business:
	go test ./internal/business/...

test-http:
	go test ./internal/http/...

list-deps:
	cd cmd/api; go list -deps -f '{{define "M"}}{{.Path}}@{{.Version}}{{end}}{{with .Module}}{{if not .Main}}{{if .Replace}}{{template "M" .Replace}}{{else}}{{template "M" .}}{{end}}{{end}}{{end}}' | sort -u

profile-coverage:
	go test ./... -coverprofile=profiles/coverage.out
	go tool cover -html=profiles/coverage.out