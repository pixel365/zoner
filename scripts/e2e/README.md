# E2E Scripts

Manual end-to-end checks against a running local EPP server.

## Prerequisites

- EPP server is running and reachable
- TLS endpoint is available (self-signed cert is fine)
- `openssl`, `perl`, `xxd`, and `timeout` are installed

## Authentication test

Script: `scripts/e2e/epp-auth-test.sh`

Run:

```bash
bash scripts/e2e/epp-auth-test.sh ok
bash scripts/e2e/epp-auth-test.sh invalid
bash scripts/e2e/epp-auth-test.sh empty
```

Custom host/port:

```bash
bash scripts/e2e/epp-auth-test.sh invalid 127.0.0.1 7000
```

Custom credentials:

```bash
EPP_USER=<user> EPP_PASS=<pass> bash scripts/e2e/epp-auth-test.sh ok
```

Expected response codes:

- `1000` for successful login
- `2200` for invalid/empty credentials
- `2400` for internal server error

Notes:

- Script sends proper EPP frame format (4-byte big-endian length + XML payload)
- Script does not suppress connection/TLS errors
