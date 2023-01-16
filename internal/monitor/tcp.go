package monitor

import (
	"context"
	"fmt"
	"net"

	"github.com/capcom6/tgbot-service-monitor/internal/config"
)

type TcpSocketProbe struct {
	Host string
	Port uint16

	dialer *net.Dialer
}

func NewTcpSocketProbe(cfg config.TCPSocket) *TcpSocketProbe {
	return &TcpSocketProbe{
		Host:   cfg.Host,
		Port:   cfg.Port,
		dialer: &net.Dialer{},
	}
}

func (p *TcpSocketProbe) Probe(ctx context.Context) error {
	conn, err := p.dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", p.Host, p.Port))
	if err != nil {
		return err
	}

	_ = conn.Close()
	return err
}
