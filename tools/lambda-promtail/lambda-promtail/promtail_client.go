package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/grafana/dskit/backoff"
)

type Client interface {
	sendToPromtail(ctx context.Context, b *batch) error
}

// Implements Client
type promtailClient struct {
	config *promtailClientConfig
	http   *http.Client
	log    *log.Logger
}

type promtailClientConfig struct {
	backoff *backoff.Config
	http    *httpClientConfig
}

type httpClientConfig struct {
	timeout       time.Duration
	skipTlsVerify bool
}

func NewPromtailClient(cfg *promtailClientConfig, log *log.Logger) *promtailClient {
	return &promtailClient{
		config: cfg,
		http:   NewHTTPClient(cfg.http),
		log:    log,
	}
}

func NewHTTPClient(cfg *httpClientConfig) *http.Client {
	//transport := http.DefaultTransport

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	// httpClient = &http.Client{
	// Timeout:   10 * time.Second,
	// Transport: t,
	// }

	if cfg.skipTlsVerify {
		t = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}
	return &http.Client{
		Timeout:   cfg.timeout,
		Transport: t,
	}
}
