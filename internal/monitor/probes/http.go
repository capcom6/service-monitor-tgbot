package probes

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
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

	start := time.Now()
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w after %s: %w", ErrInfrastructureError, time.Since(start), err)
	}

	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= http.StatusBadRequest {
		const maxBodyLen = 512
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, maxBodyLen))
		return fmt.Errorf(
			"%w: status %d, body %q, duration %s",
			ErrResponseError,
			resp.StatusCode,
			string(bodyBytes),
			time.Since(start),
		)
	}

	return nil
}
