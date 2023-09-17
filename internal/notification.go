package internal

import (
	"context"
	"github.com/ecodeclub/notify-go/internal/pkg/types"
)

type Notification struct {
	Receivers []types.Receiver
	Content   types.Content
	Channel   Channel
}

type ChannelFunc func(ctx context.Context, no *Notification) error

type Middleware func(channelFunc ChannelFunc) ChannelFunc

func NewNotification(c Channel, receivers []string, msg string) *Notification {
	var (
		recvs   []types.Receiver
		content types.Content
	)

	// TODO 根据入参构造recvs、content
	return &Notification{
		Receivers: recvs,
		Channel:   c,
		Content:   content,
	}
}

func (no *Notification) Send(ctx context.Context, mls ...Middleware) error {
	var root ChannelFunc = func(ctx context.Context, no *Notification) error {
		return no.Channel.Send(ctx, no.Receivers, no.Content)
	}

	for i := len(mls) - 1; i > 0; i-- {
		root = mls[i](root)
	}

	return root(ctx, no)
}
