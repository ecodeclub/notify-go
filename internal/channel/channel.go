package channel

import (
	"context"
	"github.com/ecodeclub/notify-go/internal/pkg/types"
	"github.com/ecodeclub/notify-go/internal/queue"
)

type SyncChannel struct {
	types.IChannel
}

type AsyncChannel struct {
	Queue queue.IQueue
	types.IChannel
}

func (s SyncChannel) Execute(ctx context.Context, deli types.Delivery) error {
	err := s.IChannel.Execute(ctx, deli)
	return err
}

func (ac AsyncChannel) Execute(ctx context.Context, deli types.Delivery) error {
	// 提前启动 channel 对应的消费者
	// 发送的时候，发送具体的 sender函数 和 参数
	err := ac.Queue.Produce(ctx, ac.IChannel, deli)
	return err
}

func NewChannel(name string, queue queue.IQueue) types.IChannel {
	var c types.IChannel
	switch name {
	case "email":
		c = NewEmailChannel(EmailConfig{})
	case "sms":
		c = NewSmsChannel(SmsConfig{})
	case "push":
		c = NewPushChannel(PushConfig{})
	default:
	}

	//return SyncChannel{
	//	IChannel: c,
	//}

	return AsyncChannel{
		IChannel: c,
		Queue:    queue,
	}
}
