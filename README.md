# Stockyard Quartermaster

**Self-hosted home inventory manager.** Single Go binary, embedded SQLite, no external dependencies.
Part of the [Stockyard](https://stockyard.dev) suite of self-hosted developer tools.

Log every item you own — electronics, furniture, appliances — with purchase date, value, and warranty expiry.

## Install

```bash
curl -fsSL https://stockyard.dev/quartermaster/install.sh | sh
```

That downloads the latest release for your platform from
[GitHub releases](https://github.com/stockyard-dev/stockyard-quartermaster/releases/latest)
and drops a single binary on disk. Linux (amd64/arm64), macOS (Intel/Apple
Silicon), and Windows (amd64) are all supported.

Then run it:

```bash
stockyard-quartermaster
```

Dashboard at [http://localhost:10230/ui](http://localhost:10230/ui).
HTTP API at `http://localhost:10230/api`. Data lives in `~/.stockyard/quartermaster/`
by default — override with `-data <path>` or `DATA_DIR=<path>`.

## Personalization

Quartermaster is one of 169 tools in the Stockyard suite. When you install it via
the [Stockyard launcher](https://github.com/stockyard-dev/stockyard-launcher) or
through the AI toolkit builder at [stockyard.dev](https://stockyard.dev), the
launcher writes a `config.json` into the data directory that personalizes the
dashboard for your specific use case — custom field labels, default values,
terminology that matches your business.

Without `config.json`, the tool runs with sensible defaults for the
"self-hosted home inventory manager" category. With it, the same binary serves a dashboard
tailored to your business in 30 seconds, no rebuild required. Read the
config schema at `GET /api/config`.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `10230` | HTTP listen port |
| `DATA_DIR` | `~/.stockyard/quartermaster` | SQLite + config directory |
| `STOCKYARD_LICENSE_KEY` | *(empty)* | Pro license key (see below) |

Command-line flags `-port` and `-data` override the env vars.

## Free vs Pro

| | Free | Pro |
|-|------|-----|
| Use case | Personal, hobby, small teams | Production, paid customers |
| Limits | Per-tool free tier (see [pricing](https://stockyard.dev/pricing/)) | Unlimited |
| Price | $0/mo | $0.99/mo per tool, or $7.99/mo for the full Stockyard suite |

Get a Pro license at [stockyard.dev/pricing/?plan=quartermaster-pro](https://stockyard.dev/pricing/?plan=quartermaster-pro).
14-day free trial, cancel anytime, your data stays on your machine either way.

## Screenshot

See the live dashboard and feature tour at
[stockyard.dev/quartermaster/](https://stockyard.dev/quartermaster/).

## License

BSL 1.1 — free for non-production use, free for production use under the
Stockyard Pro license, converts to Apache 2.0 four years after release.
See `LICENSE` for the full terms.
