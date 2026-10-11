# Minimal single-image deployment: Go serves both the API and the frontend
# production build on one port (:8080). Default database is sqlite, no nginx,
# no separate static server, no external DB required.
#
#   docker build -t go-refine-admin .
#   docker run -p 8080:8080 -v ea-data:/app/data go-refine-admin

# ---- frontend build ----
FROM node:22-alpine AS web
WORKDIR /src/frontend
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile
COPY frontend/ ./
RUN pnpm build

# ---- backend build (pure-Go sqlite, no CGO needed) ----
FROM golang:1.25-alpine AS api
WORKDIR /src/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -o /out/go-refine-admin .

# ---- runtime ----
FROM alpine:3.21
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=api /out/go-refine-admin ./
COPY backend/config/settings.yml ./config/settings.yml
COPY --from=web /src/frontend/dist ./web

ENV PORT=8080 \
    WEB_DIR=/app/web \
    SQLITE_PATH=/app/data/go-refine-admin.db

VOLUME ["/app/data", "/app/uploads"]
EXPOSE 8080
CMD ["./go-refine-admin"]
