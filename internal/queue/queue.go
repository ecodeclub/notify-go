package queue

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/ecodeclub/notify-go/internal/handler"
	"github.com/ecodeclub/notify-go/internal/pkg/logger"
	"sync"
	"time"
)

type IQueue interface {
	Produce(ctx context.Context, data []byte) error
	Consume(ctx context.Context, channel string)
}

type Kafka struct {
	Config        KafkaConfig
	topicBalancer map[string]Balancer[Topic]
	handlers      map[string]handler.Executor
}

func NewQueueService(cfg KafkaConfig) IQueue {
	var (
		balancers = map[string]Balancer[Topic]{}
		handlers  = map[string]handler.Executor{}
	)

	for channel, channelTopicsCfg := range cfg.TopicMappings {
		// 为channel类型消息建立balancer
		bala := NewBalanceBuilder[Topic](channel, channelTopicsCfg.Topics).Build(channelTopicsCfg.Strategy)
		executor := handler.NewChannelHandler(channel)
		balancers[channel] = bala
		handlers[channel] = executor
	}

	return &Kafka{Config: cfg, topicBalancer: balancers, handlers: handlers}
}

func (k *Kafka) Produce(ctx context.Context, data []byte) error {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer(k.Config.Hosts, config)
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

	// TODO data结构定义
	var channel string

	// 根据channel类型，和路由策略选取发送的topic
	topic, err := k.topicBalancer[channel].GetNext()
	if err != nil {
		logger.Panic("[Producer] choose topic fail", logger.String("channel", channel),
			logger.String("err", err.Error()))
	}

	saramaMsg := &sarama.ProducerMessage{Topic: topic.Name, Key: nil, Value: sarama.ByteEncoder(data)}

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

func (k *Kafka) Consume(ctx context.Context, channel string) {
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

		er := c.Consume(ctx, topics, k.WrapSaramaHandler(k.handlers[channel]))
		if er != nil {
			logger.Error("Consume err: ", logger.String("err", err.Error()))
		}
	}
}

func (k *Kafka) newConsumer(channel string) (sarama.ConsumerGroup, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Consumer.Return.Errors = true

	groupId := k.getGroupIdByChannel(channel)

	return sarama.NewConsumerGroup(k.Config.Hosts, groupId, saramaCfg)
}

func (k *Kafka) getTopicsByChannel(channel string) []string {
	topicCfg, ok := k.Config.TopicMappings[channel]
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
	topicCfg, ok := k.Config.TopicMappings[channel]
	if !ok {
		logger.Panic("找不到该channel的topic：%s", logger.String("channel", channel))
	}
	return topicCfg.Group
}

func (k *Kafka) WrapSaramaHandler(executor handler.Executor) sarama.ConsumerGroupHandler {
	return &ConsumeWrapper{
		Executor: executor,
	}
}

type ConsumeWrapper struct {
	Executor handler.Executor
}

func (c *ConsumeWrapper) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		delivery := handler.Delivery{}
		err := json.Unmarshal(msg.Value, &delivery)
		if err != nil {
			logger.Error("[consumer] unmarshal task detail fail", logger.String("err", err.Error()))
		}

		err = c.Executor.Execute(context.TODO(), delivery)
		if err != nil {
			logger.Error("[consumer] 执行消息发送失败",
				logger.String("topic", msg.Topic), logger.Int32("partition", msg.Partition),
				logger.Int64("offset", msg.Offset), logger.String("err", err.Error()))
			return err
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

func (c *ConsumeWrapper) Setup(session sarama.ConsumerGroupSession) error { return nil }

func (c *ConsumeWrapper) Cleanup(session sarama.ConsumerGroupSession) error { return nil }
