package task

type Message struct {
	Id          int64          `json:"id"`
	SendChannel string         `json:"send_channel"` // 消息渠道，比如是短信、邮件、推送等
	MsgContent  MessageContent `json:"msg_content"`
	MsgReceiver Receiver       `json:"msg_receiver"`
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
