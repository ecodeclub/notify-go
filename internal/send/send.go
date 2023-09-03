package send

import (
	"context"

	"github.com/ecodeclub/notify-go/internal/pkg/logger"
	"github.com/ecodeclub/notify-go/internal/pkg/mq"
	"github.com/ecodeclub/notify-go/internal/pkg/task"
	"github.com/ecodeclub/notify-go/internal/store/mysql"
)

type ISendService interface {
	Send(ctx context.Context, targets []task.Receiver, content string) error
}

type Service struct {
	nDao mysql.INotifyGoDAO
	mq   mq.IQueueService
}

func NewSendService(m mq.IQueueService, dao mysql.INotifyGoDAO) ISendService {
	return &Service{
		nDao: dao,
		mq:   m,
	}
}

func (s *Service) Send(ctx context.Context, targets []task.Receiver, content string) error {
	for _, r := range targets {
		// 写入db
		//err := s.NotifyGoDAO.InsertRecord(ctx, templateId, task.MsgReceiver, task.MsgContent.Content)

		// 发送
		err := s.mq.Produce(ctx, mq.Message{Content: content, Target: r})

		if err != nil {
			logger.Error("[send] 发送消息到消息队列失败", logger.Any("receiver", r),
				logger.String("err", err.Error()))
		}
	}

	return nil
}
