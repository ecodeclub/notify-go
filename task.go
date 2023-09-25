package notify_go

import (
	"context"
	"time"
)

type DefaultTask struct {
	Err chan error
	*Notification
}

// TriggerTask 定时触发任务
type TriggerTask struct {
	Err         chan error
	TriggerTime time.Time
	*Notification
}

// CircleTask 循环触发任务
type CircleTask struct {
	Err        chan error
	CircleTime string
	*Notification
}

func NewTriggerTask(notification *Notification, t time.Time) *TriggerTask {
	r := &TriggerTask{
		TriggerTime:  t,
		Notification: notification,
	}

	return r
}

func (tt *TriggerTask) Send(ctx context.Context) *TriggerTask {
	for {
		select {
		case <-ctx.Done():
			tt.Err <- ctx.Err()
			return tt
		case <-time.After(tt.TriggerTime.Sub(time.Now())):
			err := tt.Notification.Send(ctx)
			tt.Err <- err
		default:
		}
	}
}
