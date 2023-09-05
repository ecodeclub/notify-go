package task

const (
	email = iota
	sms
	push
)

type Receiver struct {
	Id       string
	UserName string
	Email    string
	Phone    string
	T        int8
}

func (r *Receiver) Type() string {
	switch r.T {
	case email:
		return "email"
	case sms:
		return "sms"
	case push:
		return "push"
	default:
		return "unknown"
	}
}

func (r *Receiver) Value() string {
	switch r.T {
	case email:
		return r.Email
	case sms:
		return r.Phone
	case push:
		return r.Id
	default:
		return "unknown"
	}
}
