package target

import (
	"context"

	"github.com/ecodeclub/notify-go/pkg/notifier"
)

type Service struct{}

type ITargetService interface {
	GetTarget(ctx context.Context, targetId int64) []notifier.Receiver
}

func NewTargetService() ITargetService {
	return &Service{}
}

func (ts *Service) GetTarget(ctx context.Context, targetId int64) []notifier.Receiver {
	var targets []notifier.Receiver
	return targets
}
