package probes

import (
	"context"
	"fmt"
	"net"
)

type TCPSocket struct {
	Config TCPSocketConfig

	dialer *net.Dialer
}

func NewTCPSocket(cfg TCPSocketConfig) *TCPSocket {
	return &TCPSocket{
		Config: cfg,
		dialer: new(net.Dialer),
	}
}

func (p *TCPSocket) Probe(ctx context.Context) error {
	conn, err := p.dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", p.Config.Host, p.Config.Port))
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	_ = conn.Close()
	return nil
}
