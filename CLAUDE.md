# Critiquefi

Social platform for logging and rating media (movies, TV shows, books, video games, music). See `README.md` for full architecture details.

## Monorepo layout

- `services/api` — Go API (chi, pgx/sqlx, JWT auth, golang-migrate). Only thing that talks to Postgres directly. Owns all business logic.
- `apps/web` — SvelteKit app (`adapter-node`). Renders UI, proxies API calls server-side — the browser never calls the Go API directly.
- `packages/*` — Shared Svelte components/libraries, pnpm workspace packages. Used by `apps/web` today; will also be used by a planned Tauri + Svelte native app (desktop/mobile), which is why this is a pnpm workspace monorepo instead of separate repos.

## Environments

- **Local** (Docker Compose on macOS): `docker compose --profile local up --build`. Runs `postgres`, `api`, `web`.
- **Staging**: Raspberry Pi running Coolify, deploying directly from `docker-compose.yaml` with the `local` profile (so it also runs its own Postgres container).
- **Production**: AWS EC2 running Coolify, without the `local` profile — no Postgres container starts. `DB_URL` points at AWS RDS instead. Coolify handles TLS/domain routing on EC2, same as staging.

The `postgres` service is scoped to the `local` Compose profile specifically so prod (EC2 + RDS) doesn't spin up a redundant local database.

## Images and releases

- `api`/`web` build from source locally; staging/prod pull prebuilt images via `API_IMAGE`/`WEB_IMAGE` env vars in `docker-compose.yaml`.
- `.github/workflows/release.yml` builds and pushes to GHCR (`ghcr.io/andrew-hayworth22/critiquefi-{api,web}`) on every push to `main` — staging tracks these directly.
- `.github/workflows/promote.yml` copies an existing GHCR image digest into ECR on a `vX.Y.Z` tag push — production never builds independently, only promotes what already ran on staging.
- Migrations run automatically on `api` boot (`cmd/api/main.go` calls `migrate.Up()`) — there's no separate migration step in any environment.

## Working in this repo

- Install JS deps from the repo root with `pnpm install` (pnpm workspace covers `apps/*` and `packages/*`).
- Go work happens inside `services/api` (its own `go.mod`, `Makefile` — e.g. `make test`).
- Env vars for all services are defined in `docker-compose.yaml`; `services/api/.env.example` covers running the API outside Docker.
