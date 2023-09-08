package task

import (
	"context"
)

type Executor interface {
	Name() string
	Execute(ctx context.Context, msg Message) error
}
