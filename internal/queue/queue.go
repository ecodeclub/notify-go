package queue

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/ecodeclub/notify-go/internal/pkg/logger"
	"github.com/ecodeclub/notify-go/internal/pkg/types"
	"sync"
	"time"
)

type IQueue interface {
	Produce(ctx context.Context, channel types.IChannel, delivery types.Delivery) error
	Consume(ctx context.Context, channel types.IChannel)
}

type Kafka struct {
	Config        KafkaConfig
	topicBalancer map[string]Balancer[Topic]
}

func NewQueueService(cfg KafkaConfig) IQueue {
	var balancers = map[string]Balancer[Topic]{}

	for channel, channelTopicsCfg := range cfg.TopicMappings {
		// 为channel类型消息建立balancer
		bala := NewBalanceBuilder[Topic](channel, channelTopicsCfg.Topics).Build(channelTopicsCfg.Strategy)
		balancers[channel] = bala
	}

	return &Kafka{Config: cfg, topicBalancer: balancers}
}

func (k *Kafka) Produce(ctx context.Context, c types.IChannel, delivery types.Delivery) error {
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

	// 根据channel类型，和路由策略选取发送的topic
	topic, err := k.topicBalancer[c.Name()].GetNext()
	if err != nil {
		logger.Panic("[Producer] choose topic fail", logger.String("channel", c.Name()),
			logger.String("err", err.Error()))
	}

	// 序列化data
	data, _ := json.Marshal(delivery)
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

func (k *Kafka) Consume(ctx context.Context, c types.IChannel) {
	consumer, err := k.newConsumer(c.Name())
	if err != nil {
		logger.Fatal("[kafka] 消费者启动失败", logger.String("err", err.Error()))
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		topics := k.getTopicsByChannel(c.Name())

		er := consumer.Consume(ctx, topics, k.WrapSaramaHandler(c))
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

func (k *Kafka) WrapSaramaHandler(executor types.IChannel) sarama.ConsumerGroupHandler {
	return &ConsumeWrapper{
		Executor: executor,
	}
}

type ConsumeWrapper struct {
	Executor types.IChannel
}

func (c *ConsumeWrapper) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		delivery := types.Delivery{}
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
