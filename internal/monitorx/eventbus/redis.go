package eventbus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/redis/go-redis/v9"
)

type redisBus struct {
	channelName string

	client *redis.Client
}

func newRedisBus(u *url.URL) (*redisBus, error) {
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
	}, nil
}

func (b *redisBus) Send(ctx context.Context, event interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return b.client.Publish(ctx, b.channelName, string(payload)).Err()
}
