package channel

import (
	common "brigade/internal/gen/common/v1"
	"context"
)

func NewChannel(channelName string) ChannelInterface {
	switch channelName {
	case "nats":
		return NewNatsChannel()
	}
	return nil
}

type ChannelInterface interface {
	Receive(ctx context.Context) (<-chan *common.Message, error)
}
