package mq

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/ecodeclub/notify-go/internal/pkg/logger"
	"github.com/panjf2000/ants/v2"
)

type KafkaConfig struct {
	Host   []string
	Topics []string
}

type Kafka struct {
	KafkaConfig
	pool *ants.Pool
}

func NewQueueService() IQueueService {
	return &Kafka{}
}

func (k *Kafka) Produce(ctx context.Context, msg Message) error {
	//TODO implement me
	panic("implement me")
}

func (k *Kafka) Consume(ctx context.Context, channel string, task QueueTask) {
	c, err := k.newConsumer(channel)
	if err != nil {
		logger.Fatal("[kafka] 消费者启动失败", logger.String("err", err.Error()))
	}

	handler := k.WrapSaramaHandler(task)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		topics := k.getTopicsByChannel(channel)

		err := c.Consume(ctx, topics, handler)
		if err != nil {
			logger.Error("Consume err: ", logger.String("err", err.Error()))
		}
	}
}

func (k *Kafka) newConsumer(channel string) (sarama.ConsumerGroup, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Consumer.Return.Errors = true

	groupId := k.getGroupIdByChannel(channel)

	return sarama.NewConsumerGroup(k.Host, groupId, saramaCfg)
}

func (k *Kafka) WrapSaramaHandler(task QueueTask) sarama.ConsumerGroupHandler {
	return &ConsumeWrapper{
		QueueTask: task,
		pool:      k.pool,
	}
}

func (k *Kafka) getTopicsByChannel(channel string) []string {
	return nil
}

func (k *Kafka) getGroupIdByChannel(channel string) string {
	return ""
}
