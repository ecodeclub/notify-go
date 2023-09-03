package task

import (
	"context"
)

// NewTask 要执行的发送Task的消费逻辑
func NewTask(e Executor) Task {
	return Task{
		Executor: e,
	}
}

func (t *Task) Run(ctx context.Context, detail Detail) {
	_ = t.Execute(ctx, detail)
	return
}
