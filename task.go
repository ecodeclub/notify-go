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
	"log/slog"
	"time"

	"github.com/ecodeclub/ekit/slice"
	"github.com/ecodeclub/notify-go/pkg/iterator"
	"github.com/gorhill/cronexpr"
)

// DefaultTask 默认任务
type DefaultTask struct {
	Err        chan error
	HookBefore func()
	HookAfter  func()
	*Notification
}

// TriggerTask 定时触发任务
type TriggerTask struct {
	Err         chan error
	TriggerTime time.Time
	HookBefore  func()
	HookAfter   func()
	*Notification
}

// CircleTask 循环触发任务
type CircleTask struct {
	IterCronTimes iterator.Iterable[time.Time]
	CronExpr      string
	BeginTime     time.Time
	EndTime       time.Time
	HookBefore    func()
	HookAfter     func()
	CircleNum     uint64
	CircleFailNum uint64
	*Notification
}

func NewTriggerTask(notification *Notification, t time.Time) *TriggerTask {
	r := &TriggerTask{
		Err:          make(chan error),
		TriggerTime:  t,
		Notification: notification,
		HookBefore:   func() {},
		HookAfter:    func() {},
	}
	return r
}

// Send 一次性任务
func (tt *TriggerTask) Send(ctx context.Context) {
	tt.HookBefore()
	defer tt.HookAfter()

	timer := time.After(time.Until(tt.TriggerTime))
	select {
	case <-ctx.Done():
		tt.Err <- ctx.Err()
	case <-timer:
		err := tt.Notification.Send(ctx)
		tt.Err <- err
	}
}

func NewCircleTask(notification *Notification, expr string, begin time.Time, end time.Time) *CircleTask {
	ct := &CircleTask{
		BeginTime:    begin,
		EndTime:      end,
		CronExpr:     expr,
		HookBefore:   func() {},
		HookAfter:    func() {},
		Notification: notification,
	}
	ct.fillCronTimes(begin, end)
	return ct
}

func (ct *CircleTask) fillCronTimes(begin, end time.Time) {
	cron := cronexpr.MustParse(ct.CronExpr)
	cronTimes := make([]time.Time, 0)

	// 根据cronexpr规则, 首次执行执行时间为>begin的第一个整点
	// 比如time.Now是12:00:22, 每分钟执行, 则首次执行时整数分钟, 12:01:00
	begin = cron.Next(begin)
	for begin.Before(end) && !begin.IsZero() {
		cronTimes = append(cronTimes, begin)
		begin = cron.Next(begin)
	}

	// 过滤历史时间
	// 比较的时候注意时区问题, time.Parse默认是UTC
	legalCronTimes := slice.FilterDelete[time.Time](cronTimes,
		func(idx int, src time.Time) bool {
			return src.Before(time.Now())
		})
	ct.IterCronTimes = iterator.NewListIter(legalCronTimes)
}

func (ct *CircleTask) Send(ctx context.Context) {
	ct.HookBefore()
	defer ct.HookAfter()

	select {
	case <-ctx.Done():
		return
	default:
		for {
			triggerPoint, done := ct.IterCronTimes.Next()
			if done {
				break
			}
			<-time.After(time.Until(triggerPoint))
			ct.CircleNum++
			err := ct.Notification.Send(context.TODO())
			if err != nil {
				ct.CircleFailNum++
				slog.Error("circle task execute fail", "err", err)
			}
		}
	}
}
