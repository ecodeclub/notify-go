package internal

import (
	"context"
	"github.com/ecodeclub/notify-go/internal/pkg/types"
	"github.com/pborman/uuid"
)

type Notification struct {
	types.Delivery
	Channel types.IChannel
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

func NewNotification(c types.IChannel, recvs []types.Receiver, cont types.Content) *Notification {
	no := &Notification{
		Channel: c,
		Delivery: types.Delivery{
			DeliveryID: uuid.NewUUID().String(),
			Receivers:  recvs,
			Content:    cont,
		},
	}
	return no
}
