package channel

import (
	"context"
	"github.com/ecodeclub/notify-go/internal/pkg/types"
	"time"

	"github.com/ecodeclub/notify-go/tool"
)

type SmsConfig struct{}

type SmsChannel struct{}

func NewSmsChannel(c SmsConfig) *SmsChannel {
	return &SmsChannel{}
}

func (sc *SmsChannel) Execute(ctx context.Context, deli types.Delivery) error {
	// Mock time cost
	n := tool.RandIntN(700, 800)
	time.Sleep(time.Millisecond * time.Duration(n))
	return nil
}

func (sc *SmsChannel) Name() string {
	return "sms"
}
