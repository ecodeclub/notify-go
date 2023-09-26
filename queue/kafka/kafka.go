package kafka

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/ecodeclub/notify-go/pkg/logger"
	"github.com/ecodeclub/notify-go/pkg/notifier"
)

type Topic struct {
	Name   string `toml:"name"`
	Weight int    `toml:"weight"`
}

type TopicMapping struct {
	Strategy string  `toml:"strategy"`
	Group    string  `toml:"group"`
	Topics   []Topic `toml:"topics"`
}

type Config struct {
	Hosts         []string                `toml:"host"`
	TopicMappings map[string]TopicMapping `toml:"topic_mappings"`
}

type Kafka struct {
	Config        Config
	topicBalancer map[string]Balancer[Topic]
}

func NewKafka(cfg Config) *Kafka {
	var balancers = map[string]Balancer[Topic]{}

	for channel, channelTopicsCfg := range cfg.TopicMappings {
		// 为channel类型消息建立balancer
		bala := NewBalanceBuilder[Topic](channel, channelTopicsCfg.Topics).Build(channelTopicsCfg.Strategy)
		balancers[channel] = bala
	}

	return &Kafka{Config: cfg, topicBalancer: balancers}
}

func (k *Kafka) Produce(ctx context.Context, c notifier.IChannel, delivery notifier.Delivery) error {
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

func (k *Kafka) Consume(ctx context.Context, c notifier.IChannel) {
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

func (k *Kafka) WrapSaramaHandler(executor notifier.IChannel) sarama.ConsumerGroupHandler {
	return &ConsumeWrapper{
		Executor: executor,
	}
}

type ConsumeWrapper struct {
	Executor notifier.IChannel
}

func (c *ConsumeWrapper) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		delivery := notifier.Delivery{}
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
