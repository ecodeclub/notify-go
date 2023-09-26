package sms

import (
	"context"
	"time"

	"github.com/ecodeclub/notify-go/pkg/notifier"
	"github.com/ecodeclub/notify-go/tool"
)

type Config struct{}

type ChannelSmsImpl struct{}

type Content struct{}

func NewSmsChannel(c Config) *ChannelSmsImpl {
	return &ChannelSmsImpl{}
}

func (sc *ChannelSmsImpl) Execute(ctx context.Context, deli notifier.Delivery) error {
	// Mock time cost
	n := tool.RandIntN(700, 800)
	time.Sleep(time.Millisecond * time.Duration(n))
	return nil
}

func (sc *ChannelSmsImpl) Name() string {
	return "sms"
}
