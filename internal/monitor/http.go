package monitor

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/capcom6/tgbot-service-monitor/internal/config"
)

type HttpPinger struct {
	Host        string
	Port        uint16
	Scheme      string
	Path        string
	HTTPHeaders map[string][]string

	client *http.Client
}

func NewHttpPinger(cfg config.HTTPGet) *HttpPinger {
	headers := make(map[string][]string, len(cfg.HTTPHeaders))
	for _, v := range cfg.HTTPHeaders {
		headers[v.Name] = append(headers[v.Name], v.Value)
	}

	return &HttpPinger{
		Host:        cfg.Host,
		Port:        cfg.Port,
		Scheme:      cfg.Scheme,
		Path:        cfg.Path,
		HTTPHeaders: headers,
		client:      &http.Client{},
	}
}

func (p *HttpPinger) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s://%s:%d%s", p.Scheme, p.Host, p.Port, p.Path), nil)
	if err != nil {
		return err
	}

	req.Header = p.HTTPHeaders

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}

	return nil
}
