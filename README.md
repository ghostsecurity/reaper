# Reaper

Reaper is a MITM HTTPS proxy for application security testing by [Ghost Security](https://ghostsecurity.ai).

It intercepts HTTP and HTTPS traffic, logs requests and responses to a local SQLite database, and provides a CLI for searching and inspecting captured traffic. Out-of-scope traffic passes through untouched. Reaper is designed to be easy to use by humans and AI agents alike.

## Install

### From source

```
git clone https://github.com/ghostsecurity/reaper.git
cd reaper
make build
```

### From release

Download the latest binary from [Releases](https://github.com/ghostsecurity/reaper/releases).

### Verifying Release Signatures

All release artifacts are signed with [Sigstore cosign](https://github.com/sigstore/cosign) for supply chain security.

```bash
# Install cosign
brew install cosign  # macOS
# or download from https://github.com/sigstore/cosign/releases

# Verify a release artifact
cosign verify-blob reaper_linux_amd64.tar.gz \
  --bundle reaper_linux_amd64.tar.gz.sigstore.json \
  --certificate-identity-regexp 'https://github.com/ghostsecurity/reaper/.github/workflows/release.yml' \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com'
```

**macOS Security Warning:**

When running the binary on macOS, you may see a Gatekeeper warning. This is because the binary is not signed with an Apple Developer certificate. To bypass:

```bash
# Remove quarantine attribute
xattr -d com.apple.quarantine ./reaper

# Or right-click the binary in Finder and select "Open"
```

The binary is safe to run - verify with cosign signatures above.

## Quick start

Start the proxy with one or more target domains:

```
reaper start --domains example.com
```

Configure your browser or tool to use `http://localhost:8443` as an HTTP proxy. HTTPS traffic to in-scope domains will be intercepted using a generated CA certificate — you'll need to accept the self-signed certificate warning or configure your client to skip TLS verification.

```
# curl
curl -x http://localhost:8443 -k https://api.example.com/users
```

```
# Chrome on macOS
/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome \
  --proxy-server=localhost:8443 \
  --user-data-dir="$HOME/.chrome-reaper" \
  --ignore-certificate-errors
```

```
# Chrome on Linux
google-chrome \
  --proxy-server="localhost:8443" \
  --user-data-dir="$HOME/.chrome-reaper" \
  --ignore-certificate-errors
```

In another terminal, inspect what was captured:

```
# show logs
reaper logs
# get the full request/response of the first connection that we captured
reaper get 1
```

Stop the proxy with `Ctrl+C` or with the `stop` command (if started in daemon mode).

```
reaper stop
```

## How it works

### Scope filtering

Reaper only intercepts traffic to domains or hosts you specify. Everything else passes through as a blind TCP relay with no logging or MITM.

- `--domains example.com` — suffix match. Intercepts `example.com`, `api.example.com`, `deep.sub.example.com`, etc.
- `--hosts api.example.com` — exact match. Only intercepts that specific hostname.

Both flags can be combined and accept multiple comma-separated values.

### Live activity

When running in foreground mode, the proxy prints live connection activity:

```
14:30:05 ⇄ GET https://api.example.com/users 200 142ms
14:30:08 ⇄ POST https://api.example.com/login 200 89ms
14:30:12 ⇄ GET http://example.com/health 200 53ms
```

- `⇄` — intercepted (in scope, request and response captured)
- `=` — passthrough (out of scope, blind relay)

## Commands

### `reaper start`

Start the MITM proxy.

```
reaper start --domains example.com,api.io --hosts special.internal.host
reaper start --domains example.com --port 9090
reaper start --domains example.com -d
```

| Flag | Description |
|------|-------------|
| `--domains` | Domain suffixes to intercept (comma-separated) |
| `--hosts` | Exact hostnames to intercept (comma-separated) |
| `--port` | Proxy listen port (default: `8443`) |
| `-d, --daemon` | Run as a background daemon |

At least one `--domains` or `--hosts` is required.

### `reaper stop`

Stop a running daemon.

```
reaper stop
```

### `reaper logs`

Show recent captured entries as a table.

```
reaper logs
reaper logs -n 100
```

| Flag | Description |
|------|-------------|
| `-n, --number` | Number of entries to show (default: `50`) |

Output columns: ID, METHOD, HOST, PATH, STATUS, MS (duration), REQ (request body size), RES (response body size).

### `reaper search`

Search captured entries with filters.

```
reaper search --method POST
reaper search --domains example.com --status 200
reaper search --host *.api.example.com
reaper search --path /api/v2/*
```

| Flag | Description |
|------|-------------|
| `--method` | Filter by HTTP method (exact match) |
| `--host` | Filter by hostname (supports `*` wildcard) |
| `--domains` | Filter by domain suffix |
| `--path` | Filter by path (prefix match, supports `*` wildcard) |
| `--status` | Filter by HTTP status code |
| `-n, --limit` | Max results (default: `100`) |

### `reaper get <id>`

Print the full raw HTTP request and response for an entry.

```
reaper get 42
```

### `reaper req <id>`

Print just the raw HTTP request.

```
reaper req 42
```

### `reaper res <id>`

Print just the raw HTTP response.

```
reaper res 42
```

### `reaper version`

Print version information.

## Data directory

Reaper stores its data in `~/.reaper/`:

| File | Purpose |
|------|---------|
| `reaper.db` | SQLite database with captured entries |
| `reaper.sock` | Unix socket for IPC between CLI and daemon |
| `reaper.pid` | Daemon PID file |

The CA certificate and key are generated in memory on each start and are not persisted.

## Daemon mode

By default, `reaper start` runs in the foreground with live activity output. Use `-d` to run as a background daemon:

```
reaper start --domains example.com -d
```

The parent process waits for the daemon to be ready, then exits. Use `reaper stop` to shut it down. All CLI commands (`logs`, `search`, `get`, `req`, `res`) communicate with the running daemon over a Unix socket.

## Development

```
make build       # Build the binary
make test        # Run tests
make lint        # Run linters
make dev         # Run with air (hot reload)
make run         # Run directly
```

## License

Apache License 2.0 — see [LICENSE](LICENSE).
