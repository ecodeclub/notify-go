package target

import (
	"context"

	"github.com/ecodeclub/notify-go/internal/pkg/task"
)

type Service struct{}

type ITargetService interface {
	GetTarget(ctx context.Context, targetId int64) []task.Receiver
}

func NewTargetService() ITargetService {
	return &Service{}
}

func (ts *Service) GetTarget(ctx context.Context, targetId int64) []task.Receiver {
	targets := []task.Receiver{
		{Id: 111},
		{Email: "ch_hakun@163.com"},
		{Phone: "+8618800187099"},
	}
	return targets
}
