package handler

func NewChannelHandler(channel string) Executor {
	var e Executor

	if channel == "email" {
		e = NewEmailHandler(EmailConfig{})
	} else if channel == "sms" {
		e = NewSmsHandler(SmsConfig{})
	} else if channel == "push" {
		e = NewPushHandler(PushConfig{})
	}

	return e
}
