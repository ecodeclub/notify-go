package channel

import (
	"context"
	"github.com/ecodeclub/notify-go/internal/pkg/types"
	"github.com/ecodeclub/notify-go/tool"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
)

type EmailConfig struct {
	Addr string
	Auth smtp.Auth
}

type EmailChannel struct {
	EmailClient *email.Email
	cfg         EmailConfig
}

func NewEmailChannel(cfg EmailConfig) *EmailChannel {
	return &EmailChannel{
		EmailClient: email.NewEmail(),
		cfg:         cfg,
	}
}

func (ec *EmailChannel) Execute(ctx context.Context, deli types.Delivery) error {
	// Mock time cost
	n := tool.RandIntN(700, 800)
	time.Sleep(time.Millisecond * time.Duration(n))
	return nil
}

func (ec *EmailChannel) Name() string {
	return "email"
}
