package probes

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type HTTPGet struct {
	Config HTTPGetConfig

	client *http.Client
}

func NewHTTPGet(cfg HTTPGetConfig) *HTTPGet {
	return &HTTPGet{
		Config: cfg,
		client: &http.Client{},
	}
}

func (p *HTTPGet) Probe(ctx context.Context) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s://%s:%d%s", p.Config.Scheme, p.Config.Host, p.Config.Port, p.Config.Path),
		nil,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrConfigError, err)
	}

	req.Header = p.Config.HTTPHeaders

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInfrastructureError, err)
	}

	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("%w: status %d", ErrResponseError, resp.StatusCode)
	}

	return nil
}
