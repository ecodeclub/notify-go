package task

import "context"

type Task struct {
	TaskInfo Detail
	Executor
}

type Executor interface {
	Name() string
	Execute(ctx context.Context, taskDetail Detail) error
}

type Detail struct {
	TaskId      int64          `json:"task_id"`
	SendChannel string         `json:"send_channel"` // 消息渠道，比如是短信、邮件、推送等
	MsgContent  MessageContent `json:"msg_content"`
	MsgReceiver Receiver       `json:"msg_receiver"`
}

type Receiver struct {
	Id       string
	UserName string
	Email    string
	Phone    string
	T        int8
}

type MessageContent struct {
	Email EmailContent
	Sms   SmsContent
	Push  PushContent
}

type EmailContent struct {
	From    string
	To      []Receiver
	Subject string
	Cc      []Receiver
	Text    string
	Html    string
}

type SmsContent struct {
}
type PushContent struct {
}
