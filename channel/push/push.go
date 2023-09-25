package push

import (
	"context"
	"time"

	"github.com/ecodeclub/notify-go/pkg/notifier"
	"github.com/ecodeclub/notify-go/tool"
)

type Config struct{}

type ChannelPushImpl struct{}

type Content struct{}

func NewPushChannel(c Config) *ChannelPushImpl {
	return &ChannelPushImpl{}
}

func (pc *ChannelPushImpl) Execute(ctx context.Context, deli notifier.Delivery) error {
	// Mock time cost
	n := tool.RandIntN(700, 800)
	time.Sleep(time.Millisecond * time.Duration(n))
	return nil
}

func (pc *ChannelPushImpl) Name() string {
	return "push"
}
