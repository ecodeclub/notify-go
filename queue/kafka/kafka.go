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

package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/ecodeclub/notify-go/pkg/log"

	"github.com/IBM/sarama"
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
	logger := log.FromContext(ctx)
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer(k.Config.Hosts, config)
	if err != nil {
		slog.Error("[mq] 创建生产者出错", "err", err)
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
		logger.Error("[Producer] choose topic fail", "channel", c.Name(), "err", err)
	}

	// 序列化data
	data, _ := json.Marshal(delivery)
	saramaMsg := &sarama.ProducerMessage{Topic: topic.Name, Key: nil, Value: sarama.ByteEncoder(data)}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	select {
	case producer.Input() <- saramaMsg:
		logger.Info("[mq] 发送消息成功")
	case <-ctx.Done():
		logger.Warn("[mq] 发送消息超时")
	}
	cancel()
	wg.Wait()

	return nil
}

func (k *Kafka) Consume(ctx context.Context, c notifier.IChannel) {
	logger := log.FromContext(ctx)
	consumer, err := k.newConsumer(c.Name())
	if err != nil {
		logger.Error("[kafka] 消费者启动失败", "err", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		topics := k.getTopicsByChannel(c.Name())

		er := consumer.Consume(ctx, topics, k.WrapSaramaHandler(ctx, c))
		if er != nil {
			logger.Error("Consume err: ", "err", err)
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
		return nil
	}
	topics := make([]string, 0, len(topicCfg.Topics))

	for _, item := range topicCfg.Topics {
		topics = append(topics, item.Name)
	}
	return topics
}

func (k *Kafka) getGroupIdByChannel(channel string) string {
	return k.Config.TopicMappings[channel].Group
}

func (k *Kafka) WrapSaramaHandler(ctx context.Context, executor notifier.IChannel) sarama.ConsumerGroupHandler {
	logger := log.FromContext(ctx)
	return &ConsumeWrapper{
		logger:   logger,
		Executor: executor,
	}
}

type ConsumeWrapper struct {
	logger   *slog.Logger
	Executor notifier.IChannel
}

func (c *ConsumeWrapper) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		delivery := notifier.Delivery{}
		err := json.Unmarshal(msg.Value, &delivery)
		if err != nil {
			c.logger.Error("[consumer] unmarshal task detail fail", "err", err)
		}
		err = c.Executor.Execute(context.TODO(), delivery)
		if err != nil {
			c.logger.Error("[consumer] 执行消息发送失败",
				"topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset, "err", err)
			return err
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

func (c *ConsumeWrapper) Setup(session sarama.ConsumerGroupSession) error { return nil }

func (c *ConsumeWrapper) Cleanup(session sarama.ConsumerGroupSession) error { return nil }
