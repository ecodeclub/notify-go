package internal

import (
	"context"
	"github.com/ecodeclub/notify-go/internal/pkg/types"
	"github.com/ecodeclub/notify-go/internal/queue"
)

type Channel struct {
	Queue queue.IQueue
}

func (c Channel) Send(ctx context.Context, receivers []types.Receiver, content types.Content) error {
	err := c.Queue.Produce(ctx, content.Data())
	return err
}
