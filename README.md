# Critiquefi

Critiquefi is a social platform for logging and rating media — movies, TV shows, books, video games, and music.

## Architecture

```
critiquefi/
├── apps/
│   └── web/            SvelteKit web frontend 
├── packages/           Shared Svelte components/libraries (used by web today; will be shared with a future native app)
├── services/
│   └── api/            Go API — talks to Postgres, owns business logic, auth, migrations
└── docker-compose.yaml Local/staging orchestration for all services
```

### Services

- **`services/api`** — Go API server (chi, pgx/sqlx, JWT auth, SQL migrations via golang-migrate). Owns all business logic and is the only thing that talks to Postgres directly.
- **`apps/web`** — SvelteKit app running in Node (`adapter-node`). Renders the UI and proxies API calls server-side, so the browser never talks to the Go API directly.
- **`packages/*`** — Shared Svelte component/library code, managed as pnpm workspace packages. Currently consumed by `apps/web`; will also be consumed by a planned Tauri + Svelte native app (desktop/mobile) so UI code isn't duplicated across web and native.

### Why a monorepo + pnpm workspaces

The pnpm workspace (`pnpm-workspace.yaml`) exists so that Svelte UI code in `packages/*` can be shared between `apps/web` today and the upcoming Tauri-based native app, without publishing packages to a registry.

## Environments

| Environment | Host                              | Postgres                          | Notes |
|-------------|------------------------------------|------------------------------------|-------|
| Local       | Docker Compose on macOS            | Local Postgres container (`local` compose profile) | Full stack runs via `docker compose --profile local up` |
| Staging     | Raspberry Pi running [Coolify](https://coolify.io/) | Local Postgres container (`local` profile) | Coolify deploys directly from `docker-compose.yaml` |
| Production  | AWS EC2 running Coolify            | AWS RDS (Postgres)                 | `local` profile is omitted, so no Postgres container is started — the API connects to RDS via `DB_URL`. Coolify also handles TLS termination/routing on this instance, same as staging |

The `postgres` service in `docker-compose.yaml` is scoped to the `local` Compose profile specifically so that staging and local both get a self-contained Postgres instance, while production (EC2 + RDS) does not spin up a redundant database container — it just points `DB_URL` at the RDS endpoint.

## Running locally

```bash
# install JS dependencies (web + shared packages)
pnpm install

# start the full stack, including a local Postgres container
docker compose --profile local up --build
```

This starts:
- `postgres` — Postgres 18, seeded from `POSTGRES_USER`/`POSTGRES_PASSWORD`/`POSTGRES_DB`
- `api` — Go API on `http://localhost:8080`
- `web` — SvelteKit server on `http://localhost:3000`

Migrations run automatically — `cmd/api/main.go` runs `migrate.Up()` against `DB_URL` on every boot, before the server starts accepting requests. There's no separate migration step to run in any environment, including production deploys.

### Configuration

Each service reads its configuration from environment variables (see `docker-compose.yaml` for the full list and defaults, and `services/api/.env.example` for local API development outside Docker). Notably:

- `DB_URL` — Postgres connection string (local container, or RDS in production)
- `JWT_SECRET`, `ACCESS_TOKEN_TTL`, `REFRESH_TOKEN_TTL` — auth token config
- `CORS_ORIGINS` — origins allowed to call the API directly
- `EMAIL_PROVIDER`, `AWS_*`, `FROM_ADDRESS` — outbound email (SES-backed in non-local environments)

## Deployment

- **Staging**: The same `docker-compose.yaml` used locally is pointed at by Coolify running on a Raspberry Pi, using the `local` profile so it also runs its own Postgres container.
- **Production**: Coolify running on the AWS EC2 instance deploys the `api` and `web` services from `docker-compose.yaml` without the `local` profile, so no Postgres container starts. `DB_URL` is set to the AWS RDS instance instead. Coolify handles TLS/domain routing on EC2 the same way it does on the staging Pi.

### Image registries and versioning

Locally, `api` and `web` build from source (`build:` in `docker-compose.yaml`). Staging and production instead pull prebuilt images via the `API_IMAGE` / `WEB_IMAGE` environment variables (Coolify sets these per-environment; they default to a local tag when unset):

- **Staging** pulls from GitHub Container Registry — `ghcr.io/andrew-hayworth22/critiquefi-{api,web}`. Each service has its own workflow, `.github/workflows/api.yml` and `.github/workflows/web.yml`, path-filtered to that service's files. Both build+push to GHCR (tagged `latest` and `sha-<commit>`) only on push to `main`, and only after that service's own lint/test job passes. The same workflows also run their test/lint job (and, for web, an advisory Playwright suite) on every PR into `main` as required checks — pushing/building an image never happens on a PR, only on `main`.
- **Production** pulls from AWS ECR, to avoid egress/data-transfer costs pulling into EC2 from outside AWS. Production images are never built independently — `.github/workflows/promote.yml` copies the exact GHCR image digest for a commit into ECR (via `docker buildx imagetools create`) when a version tag (`vX.Y.Z`) is pushed. This guarantees production runs the identical artifact that was already running on staging, not a fresh rebuild.

To ship a production release: confirm the commit is healthy on staging, then `git tag vX.Y.Z && git push origin vX.Y.Z` — the tag must point at the exact commit currently verified on staging, not just whatever's newest on `main`. `promote.yml` requires one-time setup — see the comment header in that workflow file for the ECR repos and IAM role it expects.

**Note**: neither GHCR nor ECR receiving a new image automatically means Coolify redeploys it — how Coolify on the Pi/EC2 notices and rolls out a new `latest` image (registry polling vs. a webhook) is configured in Coolify itself, not in this repo.

## Roadmap

- **Tauri + Svelte native app** — a desktop/mobile client sharing UI code with `apps/web` via the `packages/*` workspace packages.
