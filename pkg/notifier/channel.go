package notifier

import "context"

type Delivery struct {
	DeliveryID string
	Receivers  []Receiver
	Content    Content
}

type IChannel interface {
	Name() string
	Execute(ctx context.Context, deli Delivery) error
}
