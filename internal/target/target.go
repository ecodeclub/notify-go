package target

import (
	"context"
)

type Service struct{}

type ITargetService interface {
	GetTarget(ctx context.Context, targetId int64) []string
}

func NewTargetService() ITargetService {
	return &Service{}
}

func (ts *Service) GetTarget(ctx context.Context, targetId int64) []string {
	targets := []string{""}
	return targets
}
