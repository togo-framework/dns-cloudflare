package cloudflare

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpsertRecordCreate(t *testing.T) {
	var gotAuth, gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth, gotMethod, gotPath = r.Header.Get("Authorization"), r.Method, r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet { // ListRecords → empty so we create
			json.NewEncoder(w).Encode(map[string]any{"success": true, "result": []any{}})
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"success": true, "result": map[string]any{"id": "rec_123"}})
	}))
	defer srv.Close()

	p := &provider{token: "tkn", zone: "zone1", base: srv.URL, hc: srv.Client()}
	id, err := p.UpsertRecord(context.Background(), "", Record())
	if err != nil {
		t.Fatalf("upsert: %v", err)
	}
	if id != "rec_123" {
		t.Fatalf("id=%q", id)
	}
	if gotAuth != "Bearer tkn" {
		t.Fatalf("auth=%q", gotAuth)
	}
	if gotMethod != http.MethodPost || gotPath != "/zones/zone1/dns_records" {
		t.Fatalf("method=%s path=%s", gotMethod, gotPath)
	}
}

func TestUnsupported(t *testing.T) {
	p := &provider{}
	if _, err := p.UpsertProxyHost(context.Background(), proxyHost()); err == nil {
		t.Fatal("want ErrUnsupported")
	}
}
