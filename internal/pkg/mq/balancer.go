package mq

import (
	"fmt"
	"sync/atomic"
)

type Balancer interface {
	Name() string
	GetNext() (string, error)
}

type RRBalance struct {
	cnt    int64
	topics []Topic
}

func (rr *RRBalance) Name() string {
	return "round-robin"
}

func (rr *RRBalance) GetNext() (string, error) {
	if len(rr.topics) == 0 {
		return "", fmt.Errorf("没有可选择的topic")
	}
	cnt := atomic.AddInt64(&rr.cnt, 1)
	index := cnt % int64(len(rr.topics))
	return rr.topics[index].Name, nil
}

type BalanceBuilder struct {
	name   string
	topics []Topic
}

func NewBalanceBuilder(name string, topics []Topic) *BalanceBuilder {
	return &BalanceBuilder{
		name:   name,
		topics: topics,
	}
}

func (bb *BalanceBuilder) Build(name string) Balancer {
	switch name {
	case "round-robin":
		return &RRBalance{cnt: -1, topics: bb.topics}
	default:
		return nil
	}
}
