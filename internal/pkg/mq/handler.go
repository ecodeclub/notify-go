package mq

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/ecodeclub/notify-go/internal/pkg/logger"
	"github.com/ecodeclub/notify-go/internal/pkg/task"
	"github.com/panjf2000/ants/v2"
)

type ConsumeWrapper struct {
	task.Executor
	pool *ants.Pool
}

func (cw *ConsumeWrapper) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {

		taskDetail := task.Message{}

		err := json.Unmarshal(msg.Value, &taskDetail)
		if err != nil {
			logger.Error("[consumer] unmarshal task detail fail", logger.String("err", err.Error()))
		}

		err = cw.pool.Submit(func() {
			_ = cw.Execute(context.TODO(), taskDetail)
		})
		if err != nil {
			logger.Error("[consumer] 执行消息发送失败",
				logger.String("topic", msg.Topic), logger.Int32("partition", msg.Partition),
				logger.Int64("offset", msg.Offset), logger.String("err", err.Error()))
			return err
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}

func (cw *ConsumeWrapper) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (cw *ConsumeWrapper) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
