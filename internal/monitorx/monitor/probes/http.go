package probes

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type HttpGet struct {
	Config HttpGetConfig

	client *http.Client
}

func NewHttpGet(cfg HttpGetConfig) *HttpGet {
	return &HttpGet{
		Config: cfg,
		client: &http.Client{},
	}
}

func (p *HttpGet) Probe(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s://%s:%d%s", p.Config.Scheme, p.Config.Host, p.Config.Port, p.Config.Path), nil)
	if err != nil {
		return err
	}

	req.Header = p.Config.HTTPHeaders

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}

	return nil
}
