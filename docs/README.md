# dns-cloudflare — docs

**Cloudflare DNS.** Authoritative DNS records via the Cloudflare v4 API.

## Install

```bash
togo install togo-framework/dns-cloudflare
```

Registers on the [`dns`](https://github.com/togo-framework/dns) base; select it with **dns.provider in togo.yaml (or DNS_DRIVER)**, then use **`togo proxy`**.

## Interface

`Provider` — `UpsertRecord`/`DeleteRecord`/`ListRecords`, `UpsertProxyHost`/`DeleteProxyHost`, `UpsertRoute`/`DeleteRoute`.

## Configuration

| Env var | Description |
|---|---|
| `CLOUDFLARE_API_TOKEN` | Cloudflare API token with DNS edit permission (required). |
| `CLOUDFLARE_ZONE_ID` | Cloudflare zone id the records belong to (required). |
| `CLOUDFLARE_API_BASE` | Override the Cloudflare API base URL (testing). Optional. |

## Usage & notes

Idempotent `UpsertRecord` (list→PUT-or-POST) for A/CNAME/TXT, plus list/delete. Proxy/route ops return `ErrUnsupported` (use npm/caddy/kong for those).

## Example

```bash
togo proxy:host:add app.example.com http://localhost:3000 --provider cloudflare --dry-run
```

## Links

- [Cloudflare API](https://developers.cloudflare.com/api/)
- [Marketplace](https://to-go.dev/marketplace)
- [Source](https://github.com/togo-framework/dns-cloudflare)
