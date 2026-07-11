# Build source code
FROM golang:1.26.2-alpine AS builder
WORKDIR /src

RUN apk add --no-cache git ca-certificates

ARG VERSION=1.0.0
ARG COMMIT=none

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" \
    -o /out/api ./cmd/api

FROM alpine:latest AS runtime

RUN apk add --no-cache ca-certificates tzdata
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=builder /out/api /app/api
COPY /migrations /app/migrations

ENV PORT=8080
EXPOSE ${PORT}

HEALTHCHECK --interval=30s --timeout=3s --retries=3 CMD wget -qO- "http://127.0.0.1:${PORT}/liveness" >/dev/null 2>&1 || exit 1

USER app
ENTRYPOINT ["/app/api"]