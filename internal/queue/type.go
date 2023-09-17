package queue

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
