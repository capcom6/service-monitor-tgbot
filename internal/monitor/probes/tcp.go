package probes

import (
	"context"
	"fmt"
	"net"
)

type TcpSocket struct {
	Config TcpSocketConfig

	dialer *net.Dialer
}

func NewTcpSocket(cfg TcpSocketConfig) *TcpSocket {
	return &TcpSocket{
		Config: cfg,
		dialer: &net.Dialer{},
	}
}

func (p *TcpSocket) Probe(ctx context.Context) error {
	conn, err := p.dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", p.Config.Host, p.Config.Port))
	if err != nil {
		return err
	}

	_ = conn.Close()
	return err
}
