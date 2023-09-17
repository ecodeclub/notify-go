package handler

import (
	"context"
	"time"

	"github.com/ecodeclub/notify-go/tool"
)

type PushConfig struct{}

type PushHandler struct{}

func NewPushHandler(c PushConfig) *PushHandler {
	return &PushHandler{}
}

func (fh *PushHandler) Name() string {
	return "push"
}

func (fh *PushHandler) Execute(ctx context.Context, delivery Delivery) error {
	// Mock time cost
	n := tool.RandIntN(700, 800)
	time.Sleep(time.Millisecond * time.Duration(n))
	return nil
}
