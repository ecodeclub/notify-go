package sender

import (
	"github.com/ecodeclub/notify-go/internal/pkg/item"
	"github.com/ecodeclub/notify-go/internal/pkg/task"
)

type HandleManager struct {
	manager *item.Manager
}

func NewHandlerManager(ec EmailConfig, sc SmsConfig, pc PushConfig) *HandleManager {
	return &HandleManager{
		manager: item.NewManager(
			NewEmailHandler(ec),
			NewSmsHandler(sc),
			NewPushHandler(pc),
		),
	}
}

func (hm *HandleManager) Get(key string) (resp task.Executor, err error) {
	if h, err := hm.manager.Get(key); err == nil {
		return h.(task.Executor), nil
	}
	return nil, err
}
