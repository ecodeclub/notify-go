package sender

import (
	"context"
	"time"

	"github.com/ecodeclub/notify-go/internal/pkg/task"
	"github.com/ecodeclub/notify-go/tool"
)

const SmsNAME = "sms"

type SmsConfig struct{}

type SmsHandler struct{}

func NewSmsHandler(c SmsConfig) *SmsHandler {
	return &SmsHandler{}
}

func (fh *SmsHandler) Name() string {
	return SmsNAME
}

func (fh *SmsHandler) Execute(ctx context.Context, taskDetail task.Detail) error {
	if taskDetail.SendChannel != "sms" {
		return nil
	}
	// Mock time cost
	n := tool.RandIntN(700, 800)
	time.Sleep(time.Millisecond * time.Duration(n))
	return nil
}
