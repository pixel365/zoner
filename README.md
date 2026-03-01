# zoner

Experimental EPP server for managing domain zones.

## Status

This project is in **early, slow-paced development**.  
APIs, behavior, and internal architecture may change at any time.

## Environment Variables

`cmd/epp` loads `.env` automatically on startup via `godotenv` (if file exists).
Use [.env.example](./.env.example) as a template for local setup.

### Runtime (`cmd/epp`)

| Variable                      | Required      | Default              | Used by                      | Description                                                                               |
|-------------------------------|---------------|----------------------|------------------------------|-------------------------------------------------------------------------------------------|
| `CONFIG_PATH`                 | Yes           | -                    | config loader                | Path to YAML config file. Service panics at startup if empty or unreadable.               |
| `HEALTH_LISTEN_ADDR`          | No            | `:8081`              | health server                | HTTP listen address for `/livez`, `/readyz`, `/healthz`.                                  |
| `METRICS_ENABLED`             | No            | `false`              | metrics collector            | Enables OpenTelemetry metrics when `true`. If not `true`, noop metrics collector is used. |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Conditionally | -                    | OTEL metrics                 | OTLP gRPC endpoint for metric export (required when `METRICS_ENABLED=true`).              |
| `OTEL_EXPORTER_OTLP_INSECURE` | No            | `false`              | OTEL metrics                 | Use insecure OTLP transport (`true`) instead of TLS.                                      |
| `SERVICE_VERSION`             | No            | `0.0.1`              | OTEL metrics resource        | Value for `service.version` in OTEL resource attributes.                                  |
| `SERVICE_NAME`                | No            | `unknown`            | logger config                | Fallback service name for logs if not set via logger options.                             |
| `CI_BUILD_TAG`                | No            | `0.0.1`              | logger config                | Fallback log `service.version`.                                                           |
| `TARGET_SYSTEM`               | No            | `dev`                | logger config                | Fallback log `service.environment`.                                                       |
| `POD_NAME`                    | No            | `dev`                | logger config                | Fallback log `service.instance.id`.                                                       |
| `HOSTNAME`                    | No            | OS hostname / `GOOS` | logger helpers               | Hostname override used in logger host metadata.                                           |
| `POSTGRES_DB`                 | Yes           | -                    | postgres config              | Database name for PostgreSQL DSN.                                                         |
| `POSTGRES_USER`               | Yes           | -                    | postgres config              | Username for PostgreSQL DSN.                                                              |
| `POSTGRES_PASSWORD`           | Yes           | -                    | postgres config              | Password for PostgreSQL DSN.                                                              |
| `POSTGRES_HOST`               | Yes           | -                    | postgres config              | PostgreSQL host for connections (e.g. `postgres` in docker network).                      |
| `POSTGRES_PORT`               | Yes           | -                    | postgres config, dev compose | PostgreSQL port for app DSN and local `docker-compose.dev.yaml` port mapping.             |
| `POSTGRES_SSL_MODE`           | Yes           | -                    | postgres config              | `sslmode` value in PostgreSQL DSN (`disable`, `require`, etc.).                           |
| `POSTGRES_DATA_PATH`          | No            | `./data/postgres`    | dev compose                  | Local host path used by `docker-compose.dev.yaml` volume for postgres data.               |

### Build/Test

| Variable       | Type             | Required                | Used by                                 | Description                                                 |
|----------------|------------------|-------------------------|-----------------------------------------|-------------------------------------------------------------|
| `SERVICE_NAME` | Docker build arg | Yes (tests image build) | `testing.dockerfile`, integration tests | Chooses which `./cmd/<name>` binary to build in test image. |
