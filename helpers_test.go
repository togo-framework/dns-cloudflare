package cloudflare

import "github.com/togo-framework/dns"

func Record() dns.Record       { return dns.Record{Type: "A", Name: "app", Content: "1.2.3.4"} }
func proxyHost() dns.ProxyHost { return dns.ProxyHost{Domain: "x", Upstream: "http://y"} }
