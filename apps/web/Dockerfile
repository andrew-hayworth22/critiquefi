FROM node:24-alpine AS base
RUN corepack enable
WORKDIR /repo

FROM base AS deps
COPY pnpm-workspace.yaml pnpm-lock.yaml package.json ./
COPY apps/web/package.json ./apps/web/package.json
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    pnpm install --frozen-lockfile

FROM deps AS build
COPY . .
RUN pnpm --filter web build

FROM base AS runtime
LABEL org.opencontainers.image.source=https://github.com/andrew-hayworth22/critiquefi
RUN addgroup -S app && adduser -S app -G app
WORKDIR /app
COPY --from=build /repo/apps/web/build ./build
COPY --from=build /repo/apps/web/package.json ./package.json
ENV PORT=3000 NODE_ENV=production
EXPOSE ${PORT}
USER app
ENTRYPOINT [ "node", "build" ]
