package notify_go

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ecodeclub/ekit/slice"
	"github.com/ecodeclub/notify-go/pkg/notifier"
	"github.com/ecodeclub/notify-go/pkg/notifier/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

/*
测试场景:
1. 调用一次：超时、不超时没有返回error、不超时返回error
2. 并发调用多次：一般不存在这种场景, 同一个TrigerTask只会调用一次, 不同的task间并发安全
*/
func TestTriggerTask_Send(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	channel := mocks.NewMockIChannel(ctrl)

	o1 := channel.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)
	o2 := channel.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(context.DeadlineExceeded)
	o3 := channel.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(errors.New("test error appear"))
	gomock.InOrder(o1, o2, o3)

	type fields struct {
		Err          chan error
		TriggerTime  time.Time
		Notification *Notification
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name: "normal",
			fields: fields{
				Err:          make(chan error),
				TriggerTime:  time.Now().Add(1 * time.Second),
				Notification: NewNotification(channel, []notifier.Receiver{}, notifier.Content{}),
			},
		},
		{
			name: "time out",
			fields: fields{
				Err:          make(chan error),
				TriggerTime:  time.Now().Add(4 * time.Second),
				Notification: NewNotification(channel, []notifier.Receiver{}, notifier.Content{}),
			},
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "error",
			fields: fields{
				Err:          make(chan error),
				TriggerTime:  time.Now().Add(1 * time.Second),
				Notification: NewNotification(channel, []notifier.Receiver{}, notifier.Content{}),
			},
			wantErr: errors.New("test error appear"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tt := &TriggerTask{
				Err:          tc.fields.Err,
				TriggerTime:  tc.fields.TriggerTime,
				Notification: tc.fields.Notification,
				HookBefore:   func() {},
				HookAfter:    func() {},
			}
			ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
			defer cancel()
			go tt.Send(ctx)
			err := <-tc.fields.Err
			assert.Equal(t, err, tc.wantErr)
		})
	}
}

/*
测试场景:
1. 调用一次：超时、不超时没有返回error、不超时返回error
2. 并发调用多次：一般不存在这种场景
*/
func TestCircleTask_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	channel := mocks.NewMockIChannel(ctrl)
	channel.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	tests := []struct {
		name         string
		Notification *Notification
		cronExpr     string
		begin        time.Time
		end          time.Time
		wantFailCnt  uint64
		wantCnt      uint64
	}{
		{
			name:         "base",
			Notification: NewNotification(channel, []notifier.Receiver{}, notifier.Content{}),
			begin:        time.Now().Add(time.Second),
			end:          time.Now().Add(5 * time.Second),
			cronExpr:     "* * * * * * *", // 每秒执行
			wantCnt:      4,
			wantFailCnt:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := NewCircleTask(tt.Notification, tt.cronExpr, tt.begin, tt.end)
			ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
			defer cancel()
			ct.Send(ctx)
			assert.Equal(t, tt.wantCnt, ct.CircleNum)
		})
	}
}

func TestCircleTask_fillCronTimes(t *testing.T) {
	type args struct {
		expr  string
		begin time.Time
		end   time.Time
	}
	tests := []struct {
		name string
		args args
		want []time.Time
	}{
		{
			name: "every minute, 历史时间全部被过滤",
			args: args{
				expr: "* * * * *",
				begin: func() time.Time {
					ti, _ := time.ParseInLocation("2006-01-02 15:04:05", "2023-10-26 15:31:00", time.Local)
					return ti
				}(),
				end: func() time.Time {
					ti, _ := time.ParseInLocation("2006-01-02 15:04:05", "2023-10-26 15:33:00", time.Local)
					return ti
				}(),
			},
			want: make([]time.Time, 0),
		},
		{
			name: "every minute 开始时间晚于结束时间",
			args: args{
				expr: "* * * * *",
				begin: func() time.Time {
					ti, _ := time.ParseInLocation("2006-01-02 15:04:05", "2023-10-26 15:33:00", time.Local)
					return ti
				}(),
				end: func() time.Time {
					ti, _ := time.ParseInLocation("2006-01-02 15:04:05", "2023-10-26 15:31:00", time.Local)
					return ti
				}(),
			},
			want: make([]time.Time, 0),
		},
		{
			name: "every second",
			args: args{
				expr:  "* * * * * * *",
				begin: func() time.Time { return time.Now().Add(100 * time.Second) }(),
				end:   func() time.Time { return time.Now().Add(105 * time.Second) }(),
			},
			want: func() []time.Time {
				res := make([]time.Time, 0)
				begin := time.Now().Add(100 * time.Second)
				end := time.Now().Add(105 * time.Second)
				begin = begin.Add(time.Second)
				for begin.Before(end) {
					res = append(res, begin)
					begin = begin.Add(time.Second)
				}
				return res
			}(),
		},
		{
			name: "every minte",
			args: args{
				expr: "* * * * *",
				begin: func() time.Time {
					return time.Now().Add(100 * time.Minute)
				}(),
				end: func() time.Time {
					return time.Now().Add(104 * time.Minute)
				}(),
			},
			want: func() []time.Time {
				res := make([]time.Time, 0)
				begin := time.Now().Add(100 * time.Minute)
				end := time.Now().Add(104 * time.Minute)
				begin = begin.Add(time.Minute)
				for begin.Before(end) {
					res = append(res, begin.Truncate(time.Minute)) // 秒数需要置0
					begin = begin.Add(time.Minute)
				}
				return res
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := &CircleTask{
				CronExpr: tt.args.expr,
			}

			ct.fillCronTimes(tt.args.begin, tt.args.end)
			cronTimes := make([]time.Time, 0)
			for {
				ti, done := ct.IterCronTimes.Next()
				if done {
					break
				}
				cronTimes = append(cronTimes, ti)
			}
			assert.Equal(t, toStr(tt.want), toStr(cronTimes))
		})
	}
}

func toStr(arr []time.Time) []string {
	return slice.Map[time.Time](arr, func(idx int, src time.Time) string {
		return src.Format("2006-01-02 15:04:05")
	})
}
