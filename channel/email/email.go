package email

import (
	"context"
	"net/smtp"
	"time"

	"github.com/ecodeclub/notify-go/pkg/notifier"
	"github.com/ecodeclub/notify-go/tool"

	"github.com/jordan-wright/email"
)

type Config struct {
	Addr string
	Auth smtp.Auth
}

type ChannelEmailImpl struct {
	EmailClient *email.Email
	cfg         Config
}

type Content struct {
	From    []string
	Subject []string
	Cc      []string
	Html    []string
	Text    []string
}

func NewEmailChannel(cfg Config) *ChannelEmailImpl {
	return &ChannelEmailImpl{
		EmailClient: email.NewEmail(),
		cfg:         cfg,
	}
}

func (c *ChannelEmailImpl) Execute(ctx context.Context, deli notifier.Delivery) error {
	// Mock time cost
	n := tool.RandIntN(700, 800)
	time.Sleep(time.Millisecond * time.Duration(n))
	return nil
}

func (c *ChannelEmailImpl) Name() string {
	return "email"
}
