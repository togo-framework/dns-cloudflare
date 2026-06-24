// Package cloudflare is a togo dns driver for Cloudflare authoritative DNS.
// It manages A/AAAA/CNAME/TXT/MX records over the Cloudflare v4 API. Reverse-
// proxy hosts and gateway routes are not applicable and return ErrUnsupported.
//
// Install: `togo install togo-framework/dns-cloudflare`, set DNS_DRIVER=cloudflare.
package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/togo-framework/dns"
	"github.com/togo-framework/togo"
)

func init() {
	dns.RegisterDriver("cloudflare", func(k *togo.Kernel) (dns.Provider, error) {
		token := os.Getenv("CLOUDFLARE_API_TOKEN")
		if token == "" {
			return nil, fmt.Errorf("dns-cloudflare: CLOUDFLARE_API_TOKEN not set")
		}
		base := os.Getenv("CLOUDFLARE_API_BASE")
		if base == "" {
			base = "https://api.cloudflare.com/client/v4"
		}
		return &provider{
			token: token,
			zone:  os.Getenv("CLOUDFLARE_ZONE_ID"),
			base:  base,
			hc:    &http.Client{Timeout: 20 * time.Second},
		}, nil
	})
}

type provider struct {
	token, zone, base string
	hc                *http.Client
}

type cfResp struct {
	Success bool              `json:"success"`
	Errors  []json.RawMessage `json:"errors"`
	Result  json.RawMessage   `json:"result"`
}

func (p *provider) zoneID(zone string) string {
	if zone == "" {
		return p.zone
	}
	return zone
}

func (p *provider) do(ctx context.Context, method, path string, body any, out any) error {
	var rdr io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, p.base+path, rdr)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+p.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var r cfResp
	if err := json.Unmarshal(raw, &r); err != nil {
		return fmt.Errorf("dns-cloudflare: decode %s %s: %w (%s)", method, path, err, string(raw))
	}
	if !r.Success {
		return fmt.Errorf("dns-cloudflare: %s %s failed: %s", method, path, string(raw))
	}
	if out != nil && len(r.Result) > 0 {
		return json.Unmarshal(r.Result, out)
	}
	return nil
}

type cfRecord struct {
	ID       string `json:"id,omitempty"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	TTL      int    `json:"ttl,omitempty"`
	Proxied  bool   `json:"proxied"`
	Priority int    `json:"priority,omitempty"`
}

func (p *provider) UpsertRecord(ctx context.Context, zone string, r dns.Record) (string, error) {
	z := p.zoneID(zone)
	if z == "" {
		return "", fmt.Errorf("dns-cloudflare: zone id required (CLOUDFLARE_ZONE_ID or zone arg)")
	}
	ttl := r.TTL
	if ttl == 0 {
		ttl = 1 // automatic
	}
	payload := cfRecord{Type: r.Type, Name: r.Name, Content: r.Content, TTL: ttl, Proxied: r.Proxied, Priority: r.Prio}

	// Find an existing record with the same type+name to update in place.
	existing, _ := p.ListRecords(ctx, z)
	for _, e := range existing {
		if e.Type == r.Type && e.Name == r.Name {
			var got cfRecord
			if err := p.do(ctx, http.MethodPut, fmt.Sprintf("/zones/%s/dns_records/%s", z, e.ID), payload, &got); err != nil {
				return "", err
			}
			return got.ID, nil
		}
	}
	var got cfRecord
	if err := p.do(ctx, http.MethodPost, fmt.Sprintf("/zones/%s/dns_records", z), payload, &got); err != nil {
		return "", err
	}
	return got.ID, nil
}

func (p *provider) DeleteRecord(ctx context.Context, zone, id string) error {
	return p.do(ctx, http.MethodDelete, fmt.Sprintf("/zones/%s/dns_records/%s", p.zoneID(zone), id), nil, nil)
}

func (p *provider) ListRecords(ctx context.Context, zone string) ([]dns.Record, error) {
	z := p.zoneID(zone)
	var got []cfRecord
	if err := p.do(ctx, http.MethodGet, "/zones/"+url.PathEscape(z)+"/dns_records?per_page=100", nil, &got); err != nil {
		return nil, err
	}
	out := make([]dns.Record, len(got))
	for i, g := range got {
		out[i] = dns.Record{ID: g.ID, Type: g.Type, Name: g.Name, Content: g.Content, TTL: g.TTL, Proxied: g.Proxied, Prio: g.Priority}
	}
	return out, nil
}

func (p *provider) UpsertProxyHost(context.Context, dns.ProxyHost) (string, error) {
	return "", dns.ErrUnsupported
}
func (p *provider) DeleteProxyHost(context.Context, string) error { return dns.ErrUnsupported }
func (p *provider) UpsertRoute(context.Context, dns.Route) (string, error) {
	return "", dns.ErrUnsupported
}
func (p *provider) DeleteRoute(context.Context, string) error { return dns.ErrUnsupported }
