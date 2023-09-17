package handler

import (
	"context"
	"time"

	"github.com/ecodeclub/notify-go/tool"
)

type SmsConfig struct{}

type SmsHandler struct{}

func NewSmsHandler(c SmsConfig) *SmsHandler {
	return &SmsHandler{}
}

func (fh *SmsHandler) Name() string {
	return "sms"
}

func (fh *SmsHandler) Execute(ctx context.Context, delivery Delivery) error {
	// Mock time cost
	n := tool.RandIntN(700, 800)
	time.Sleep(time.Millisecond * time.Duration(n))
	return nil
}
