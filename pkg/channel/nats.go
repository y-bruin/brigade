package channel

import (
	common "brigade/internal/gen/common/v1"
	"context"
)

type NatsChannel struct {
}

func NewNatsChannel() ChannelInterface {
	return &NatsChannel{}
}

func (c *NatsChannel) Receive(ctx context.Context) (<-chan *common.Message, error) {
	return nil, nil
}
