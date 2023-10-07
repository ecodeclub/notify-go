package email

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net"
	"net/smtp"

	"github.com/ecodeclub/notify-go/pkg/notifier"
	"github.com/jordan-wright/email"
	"github.com/pkg/errors"
)

type Config struct {
	smtpHostAddr  string
	smtpAuth      smtp.Auth
	senderAddress string // e.g. Hooko <xxx.xx@xx>
}

type ChannelEmailImpl struct {
	email *email.Email
	cfg   Config
}

type Content struct {
	Subject string
	Cc      []string
	Html    string
	Text    string
}

func NewEmailChannel(cfg Config) *ChannelEmailImpl {
	return &ChannelEmailImpl{
		email: email.NewEmail(),
		cfg:   cfg,
	}
}

func (c *ChannelEmailImpl) Execute(ctx context.Context, deli notifier.Delivery) error {
	var (
		msgContent Content
		tos        []string
	)
	err := json.Unmarshal(deli.Content, &msgContent)
	if err != nil {
		return err
	}

	for _, r := range deli.Receivers {
		tos = append(tos, r.Email)
	}

	c.email.To = tos
	c.email.From = c.cfg.senderAddress
	// TODO cc不是抄送, 而是append到to内
	c.email.Cc = msgContent.Cc
	c.email.Subject = msgContent.Subject
	c.email.HTML = []byte(msgContent.Html)
	c.email.Text = []byte(msgContent.Text)

	ch := make(chan struct{})
	go func() {
		defer func() {
			close(ch)
		}()
		host, _, _ := net.SplitHostPort(c.cfg.smtpHostAddr)
		// TODO 如果SendWithTLS执行时间太长, 有goroutine泄露问题
		// 需要改造SendWithTLS 为 SendWithTLSContext()
		err = c.email.SendWithTLS(c.cfg.smtpHostAddr, c.cfg.smtpAuth, &tls.Config{ServerName: host})
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ch:
			if err != nil {
				err = errors.Wrap(err, "failed to send mail")
			}
			return err
		}
	}
}

func (c *ChannelEmailImpl) Name() string {
	return "email"
}
