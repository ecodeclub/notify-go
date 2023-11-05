// Copyright 2021 ecodeclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"

	"github.com/ecodeclub/notify-go/pkg/log"

	"github.com/ecodeclub/ekit/slice"
	"github.com/ecodeclub/notify-go/pkg/notifier"
	"github.com/jordan-wright/email"
	"github.com/pkg/errors"
)

type Config struct {
	SmtpHostAddr  string `json:"smtp_host_addr"`
	SmtpUserName  string `json:"smtp_user_name"`
	SmtpPwd       string `json:"smtp_pwd"`
	SenderAddress string `json:"sender_address"` // e.g. Hooko <xxx.xx@xx>
}

type ChannelEmailImpl struct {
	email    *email.Email
	smtpAuth smtp.Auth
	smtpHost string
	Config
}

type Content struct {
	Subject string
	Cc      []string
	Html    string
	Text    string
}

func NewEmailChannel(cfg Config) *ChannelEmailImpl {
	host, _, _ := net.SplitHostPort(cfg.SmtpHostAddr)
	return &ChannelEmailImpl{
		email:    email.NewEmail(),
		smtpHost: host,
		smtpAuth: smtp.PlainAuth("", cfg.SmtpUserName, cfg.SmtpPwd, host),
		Config:   cfg,
	}
}

func (c *ChannelEmailImpl) Execute(ctx context.Context, deli notifier.Delivery) error {
	var (
		err    error
		logger = log.FromContext(ctx)
	)
	msgContent := c.initEmailContent(deli.Content)

	c.email.To = slice.Map[notifier.Receiver, string](deli.Receivers, func(idx int, src notifier.Receiver) string {
		return src.Email
	})

	c.email.From = c.SenderAddress
	// TODO cc不是抄送, 而是append到to内
	c.email.Cc = msgContent.Cc
	c.email.Subject = msgContent.Subject
	c.email.HTML = []byte(msgContent.Html)
	c.email.Text = []byte(msgContent.Text)

	logger.Info("email execute", "params", fmt.Sprintf("to[%v], from[%s]", c.email.To, c.email.From))
	ch := make(chan struct{})
	go func() {
		defer func() {
			close(ch)
		}()
		// TODO 如果SendWithTLS执行时间太长, 有goroutine泄露问题
		// 需要改造SendWithTLS 为 SendWithTLSContext()
		err = c.email.SendWithTLS(c.SmtpHostAddr, c.smtpAuth, &tls.Config{ServerName: c.smtpHost})
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()
		logger.Error("email execute err", "err", err.Error())
	case <-ch:
		if err != nil {
			err = errors.Wrap(err, "failed to send mail")
			logger.Error("email execute err", "err", err.Error())
		}
	}

	return err
}

func (c *ChannelEmailImpl) Name() string {
	return "email"
}

func (c *ChannelEmailImpl) initEmailContent(nc notifier.Content) Content {
	cc := Content{
		Subject: nc.Title,
		Html:    string(nc.Data),
	}
	return cc
}
