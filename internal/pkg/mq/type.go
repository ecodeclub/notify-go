package mq

import (
	"context"

	"github.com/ecodeclub/notify-go/internal/pkg/task"
)

type Message struct {
	Content string
	Channel string // 消息类型
	Target  task.Receiver
}

type IQueueService interface {
	Produce(ctx context.Context, msg Message) error
	Consume(ctx context.Context, channel string, task QueueTask)
}

type QueueTask interface {
	Run(ctx context.Context, detail task.Detail)
}

type Topic struct {
	Name   string `toml:"name"`
	Weight int    `toml:"weight"`
}

type TopicMapping struct {
	Strategy string  `toml:"strategy"`
	Group    string  `toml:"group"`
	Topics   []Topic `toml:"topics"`
}

type KafkaConfig struct {
	Hosts         []string                `toml:"host"`
	TopicMappings map[string]TopicMapping `toml:"topic_mappings"`
}
