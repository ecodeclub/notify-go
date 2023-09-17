package types

import "encoding/json"

type Content interface {
	Name() string
	Data() []byte
}

type EmailMessage struct {
	From    []string
	To      []string
	Subject []string
	Cc      []string
	Content string
}

type SMSMessage struct {
	Phone   string
	Content string
}

type PushMessage struct {
	UserId  string
	Content string
}

func (e *EmailMessage) Name() string {
	return "EMAIL"
}

func (e *EmailMessage) Data() []byte {
	data, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return data
}

func (S *SMSMessage) Name() string {
	//TODO implement me
	panic("implement me")
}

func (S *SMSMessage) Data() []byte {
	//TODO implement me
	panic("implement me")
}

func (p *PushMessage) Name() string {
	//TODO implement me
	panic("implement me")
}

func (p *PushMessage) Data() []byte {
	//TODO implement me
	panic("implement me")
}
