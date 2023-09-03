package sender

import (
	"context"
	"time"

	"github.com/ecodeclub/notify-go/internal/pkg/task"
	"github.com/ecodeclub/notify-go/tool"
)

const PushNAME = "push"

type PushConfig struct{}

type PushHandler struct{}

func NewPushHandler(c PushConfig) *PushHandler {
	return &PushHandler{}
}

func (fh *PushHandler) Name() string {
	return PushNAME
}

func (fh *PushHandler) Execute(ctx context.Context, taskDetail task.Detail) error {
	if taskDetail.SendChannel != "push" {
		return nil
	}
	// Mock time cost
	n := tool.RandIntN(700, 800)
	time.Sleep(time.Millisecond * time.Duration(n))
	return nil
}
