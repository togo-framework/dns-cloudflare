<!-- togo-header -->
<div align="center">
  <picture><source media="(prefers-color-scheme: dark)" srcset=".github/assets/togo-mark-dark.svg" /><img src=".github/assets/togo-mark.svg" alt="ToGO" height="64" /></picture>
  <h1>togo-framework/dns-cloudflare</h1>
  <p>Cloudflare authoritative-DNS driver for togo dns.</p>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1F8A99" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/dns-cloudflare"><img src="https://pkg.go.dev/badge/github.com/togo-framework/dns-cloudflare.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/dns-cloudflare
```
<!-- /togo-header -->

Cloudflare authoritative-DNS driver for togo's [`dns`](https://github.com/togo-framework/dns)
subsystem. Manages records over the Cloudflare v4 API (token auth).

## Config

| Env | Meaning |
|-----|---------|
| `DNS_DRIVER` | set to `cloudflare` |
| `CLOUDFLARE_API_TOKEN` | API token with **DNS:Edit** on the zone (required) |
| `CLOUDFLARE_ZONE_ID` | default zone id (optional — can pass per call) |

```go
svc, _ := dns.FromKernel(k)
svc.UpsertRecord(ctx, "", dns.Record{Type: "A", Name: "app", Content: "203.0.113.10", Proxied: true})
```

`UpsertRecord` is idempotent: it updates an existing record of the same
type+name, otherwise creates one. Proxy hosts / gateway routes return
`dns.ErrUnsupported`.

<!-- togo-sponsors -->
---
<div align="center">
  <h3>Premium sponsors</h3>
  <p><a href="https://id8media.com"><strong>ID8 Media</strong></a> &nbsp;·&nbsp; <a href="https://one-studio.co"><strong>One Studio</strong></a></p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
<!-- /togo-sponsors -->
