package sender

import (
	"context"
	"errors"
	"net/smtp"

	"github.com/ecodeclub/notify-go/internal/pkg/task"
	"github.com/jordan-wright/email"
)

const EmailNAME = "email"

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
	return EmailNAME
}

func (eh *EmailHandler) Execute(ctx context.Context, taskDetail task.Message) error {
	var (
		tos []string
		ccs []string
	)

	if len(taskDetail.MsgContent.Email.To) == 0 {
		return errors.New("notify-go: 收件人为空")
	}

	for _, to := range taskDetail.MsgContent.Email.To {
		tos = append(tos, to.Email)
	}

	for _, cc := range taskDetail.MsgContent.Email.Cc {
		ccs = append(ccs, cc.Email)
	}

	eh.EmailClient.From = taskDetail.MsgContent.Email.From
	eh.EmailClient.To = tos
	eh.EmailClient.Cc = ccs
	eh.EmailClient.Subject = taskDetail.MsgContent.Email.Subject
	eh.EmailClient.HTML = []byte(taskDetail.MsgContent.Email.Html)

	err := eh.EmailClient.Send(eh.cfg.Addr, eh.cfg.Auth)
	return err
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
