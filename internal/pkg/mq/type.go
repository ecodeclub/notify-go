package mq

import (
	"context"

	"github.com/ecodeclub/notify-go/internal/pkg/task"
)

type Message struct {
	Content string
	Target  task.Receiver
}

type IQueueService interface {
	Produce(ctx context.Context, msg Message) error
	Consume(ctx context.Context, channel string, task QueueTask)
}

type QueueTask interface {
	Run(ctx context.Context, detail task.Detail)
}
