// Copyright 2021 ecodeclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	Deadline   string
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
		case <-time.After(time.Until(tt.TriggerTime)):
			err := tt.Notification.Send(ctx)
			tt.Err <- err
		}
	}
}
