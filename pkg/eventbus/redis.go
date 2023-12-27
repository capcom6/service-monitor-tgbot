package eventbus

import (
	"context"
	"errors"
	"net/url"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type redisBus struct {
	channelName string

	logger *zap.Logger
	client *redis.Client
}

func newRedisBus(u *url.URL, logger *zap.Logger) (*redisBus, error) {
	channelName := u.Query().Get("channel")
	if channelName == "" {
		return nil, errors.New("channel name is required")
	}

	query := u.Query()
	query.Del("channel")
	u.RawQuery = query.Encode()

	opts, err := redis.ParseURL(u.String())
	if err != nil {
		return nil, err
	}

	return &redisBus{
		channelName: channelName,
		client:      redis.NewClient(opts),
		logger:      logger,
	}, nil
}

func (b *redisBus) Send(ctx context.Context, event string) error {
	return b.client.Publish(ctx, b.channelName, event).Err()
}

func (b *redisBus) Receive(ctx context.Context) (<-chan string, error) {
	ch := make(chan string)

	go func() {
		defer close(ch)

		sub := b.client.Subscribe(ctx, b.channelName)
		defer sub.Close()

		b.logger.Debug("subscribed", zap.String("channel", b.channelName))

		for {
			select {
			case msg, ok := <-sub.Channel():
				if !ok {
					b.logger.Warn("channel closed")
					return
				}

				select {
				case ch <- msg.Payload:
				case <-ctx.Done():
					return
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return ch, nil
}

func (b *redisBus) Close() error {
	return b.client.Close()
}
