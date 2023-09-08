package item

import (
	"fmt"
)

type Manager struct {
	itemMap map[string]Itemable
}

type Itemable interface {
	Name() string
}

func NewManager(items ...Itemable) *Manager {
	im := new(Manager)
	im.Init(items...)
	return im
}

func (im *Manager) Init(items ...Itemable) {
	for _, item := range items {
		im.Register(item)
	}
}

func (im *Manager) Register(item Itemable) {
	if im.itemMap == nil {
		im.itemMap = make(map[string]Itemable)
	}
	im.itemMap[item.Name()] = item
}

func (im *Manager) Get(key string) (Itemable, error) {
	item, ok := im.itemMap[key]
	if !ok {
		return nil, fmt.Errorf("unknown item %s", key)
	}
	return item, nil
}
