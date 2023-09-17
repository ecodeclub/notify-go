package handler

import (
	"context"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailConfig struct {
	Addr string
	Auth smtp.Auth
}

type EmailHandler struct {
	EmailClient *email.Email
	cfg         EmailConfig
}

func NewEmailHandler(cfg EmailConfig) *EmailHandler {
	return &EmailHandler{
		EmailClient: email.NewEmail(),
		cfg:         cfg,
	}
}

func (eh *EmailHandler) Name() string {
	return "email"
}

func (eh *EmailHandler) Execute(ctx context.Context, delivery Delivery) error {
	return nil
}

//func (eh *EmailHandler) Execute(ctx context.Context, task *internal.Task) (err error) {
//	if task.SendChannel != "email" {
//		return nil
//	}
//  Mock time cost
//	n := tool.RandIntN(700, 800)
//	time.Sleep(time.Millisecond * time.Duration(n))
//	log.Printf("[email] %+v\n, cost: %d ms", task, n)
//	return nil
//}
