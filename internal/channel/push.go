package channel

import (
	"context"
	"github.com/ecodeclub/notify-go/internal/pkg/types"
	"time"

	"github.com/ecodeclub/notify-go/tool"
)

type PushConfig struct{}

type PushChannel struct{}

func NewPushChannel(c PushConfig) *PushChannel {
	return &PushChannel{}
}

func (pc *PushChannel) Execute(ctx context.Context, deli types.Delivery) error {
	// Mock time cost
	n := tool.RandIntN(700, 800)
	time.Sleep(time.Millisecond * time.Duration(n))
	return nil
}

func (pc *PushChannel) Name() string {
	return "push"
}
