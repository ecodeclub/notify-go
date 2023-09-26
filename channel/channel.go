package channel

import (
	"context"

	"github.com/ecodeclub/notify-go/pkg/notifier"
	"github.com/ecodeclub/notify-go/queue"
)

type SyncChannel struct {
	notifier.IChannel
}

type AsyncChannel struct {
	Queue queue.IQueue
	notifier.IChannel
}

func (s SyncChannel) Execute(ctx context.Context, deli notifier.Delivery) error {
	err := s.IChannel.Execute(ctx, deli)
	return err
}

func (ac AsyncChannel) Execute(ctx context.Context, deli notifier.Delivery) error {
	// 提前启动 channel 对应的消费者
	// 发送的时候，发送具体的 sender函数 和 参数
	err := ac.Queue.Produce(ctx, ac.IChannel, deli)
	return err
}
