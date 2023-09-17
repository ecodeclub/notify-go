package handler

import (
	"context"
	"github.com/ecodeclub/notify-go/internal/pkg/types"
)

type Delivery struct {
	Receivers []types.Receiver
	EmailMsg  types.EmailMessage
	SmsMsg    types.SMSMessage
	PushMsg   types.PushMessage
}

// Executor 实际发送的抽象
type Executor interface {
	Name() string
	Execute(ctx context.Context, delivery Delivery) error
}
