package notify_go

import (
	"context"

	"github.com/ecodeclub/notify-go/pkg/notifier"
	"github.com/pborman/uuid"
)

type Notification struct {
	notifier.Delivery
	Channel notifier.IChannel
}

type ChannelFunc func(ctx context.Context, no *Notification) error

type Middleware func(channelFunc ChannelFunc) ChannelFunc

func (no *Notification) Send(ctx context.Context, mls ...Middleware) error {
	var root ChannelFunc = func(ctx context.Context, no *Notification) error {
		return no.Channel.Execute(ctx, no.Delivery)
	}

	for i := len(mls) - 1; i > 0; i-- {
		root = mls[i](root)
	}

	return root(ctx, no)
}

func NewNotification(c notifier.IChannel, recvs []notifier.Receiver, content []byte) *Notification {
	no := &Notification{
		Channel: c,
		Delivery: notifier.Delivery{
			DeliveryID: uuid.NewUUID().String(),
			Receivers:  recvs,
			Content:    content,
		},
	}
	return no
}
