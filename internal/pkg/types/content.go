package types

type Content struct {
	Message string
	EmailContent
	SMSContent
	PushContent
}

type EmailContent struct {
	From    []string
	Subject []string
	Cc      []string
}

type SMSContent struct {
}

type PushContent struct {
}
