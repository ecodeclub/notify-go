package mq

import (
	"context"
	"encoding/json"
	"runtime"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/ecodeclub/notify-go/internal/pkg/logger"
	"github.com/ecodeclub/notify-go/internal/pkg/task"
	"github.com/panjf2000/ants/v2"
)

type Kafka struct {
	KafkaConfig
	pool          *ants.Pool
	topicBalancer map[string]Balancer
}

func NewQueueService(cfg KafkaConfig) IQueueService {
	balas := make(map[string]Balancer, len(cfg.TopicMappings))

	// 每一种类型消息建立balancer
	for channel, topics := range cfg.TopicMappings {
		bala := NewBalanceBuilder(channel, topics.Topics).Build(topics.Strategy)
		balas[channel] = bala
	}

	pool, err := ants.NewPool(runtime.NumCPU())
	if err != nil {
		panic(err)
	}

	return &Kafka{KafkaConfig: cfg, pool: pool, topicBalancer: balas}
}

func (k *Kafka) Produce(ctx context.Context, m task.Message) error {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer(k.Hosts, config)
	if err != nil {
		logger.Panic("[mq] 创建生产者出错", logger.Any("err", err.Error()))
	}
	defer producer.AsyncClose()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-producer.Successes()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-producer.Errors()
	}()

	// 根据channel类型，和路由策略选取发送的topic
	topic, err := k.topicBalancer[m.SendChannel].GetNext()
	if err != nil {
		logger.Panic("[Producer] choose topic fail", logger.String("channel", m.SendChannel),
			logger.String("err", err.Error()))
	}

	data, err := json.Marshal(m)
	if err != nil {

	}

	saramaMsg := &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.ByteEncoder(data)}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	select {
	case producer.Input() <- saramaMsg:
	case <-ctx.Done():
		logger.Warn("[mq] 发送消息超时")
	}
	cancel()
	wg.Wait()

	return nil
}

func (k *Kafka) Consume(ctx context.Context, channel string, executor task.Executor) {
	c, err := k.newConsumer(channel)
	if err != nil {
		logger.Fatal("[kafka] 消费者启动失败", logger.String("err", err.Error()))
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		topics := k.getTopicsByChannel(channel)

		er := c.Consume(ctx, topics, k.WrapSaramaHandler(executor))
		if er != nil {
			logger.Error("Consume err: ", logger.String("err", err.Error()))
		}
	}
}

func (k *Kafka) newConsumer(channel string) (sarama.ConsumerGroup, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Consumer.Return.Errors = true

	groupId := k.getGroupIdByChannel(channel)

	return sarama.NewConsumerGroup(k.Hosts, groupId, saramaCfg)
}

func (k *Kafka) WrapSaramaHandler(executor task.Executor) sarama.ConsumerGroupHandler {
	return &ConsumeWrapper{
		Executor: executor,
		pool:     k.pool,
	}
}

func (k *Kafka) getTopicsByChannel(channel string) []string {
	topicCfg, ok := k.TopicMappings[channel]
	if !ok {
		logger.Panic("找不到该channel的topic：%s", logger.String("channel", channel))
		return nil
	}
	topics := make([]string, 0, len(topicCfg.Topics))

	for _, item := range topicCfg.Topics {
		topics = append(topics, item.Name)
	}
	return topics
}

func (k *Kafka) getGroupIdByChannel(channel string) string {
	topicCfg, ok := k.TopicMappings[channel]
	if !ok {
		logger.Panic("找不到该channel的topic：%s", logger.String("channel", channel))
	}
	return topicCfg.Group
}
