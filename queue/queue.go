package queue

import (
	"context"

	"github.com/ecodeclub/notify-go/pkg/notifier"
)

type IQueue interface {
	Produce(ctx context.Context, channel notifier.IChannel, delivery notifier.Delivery) error
	Consume(ctx context.Context, channel notifier.IChannel)
}
