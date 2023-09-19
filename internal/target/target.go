package target

import (
	"context"
	"github.com/ecodeclub/notify-go/internal/pkg/types"
)

type Service struct{}

type ITargetService interface {
	GetTarget(ctx context.Context, targetId int64) []types.Receiver
}

func NewTargetService() ITargetService {
	return &Service{}
}

func (ts *Service) GetTarget(ctx context.Context, targetId int64) []types.Receiver {
	var targets []types.Receiver
	return targets
}
