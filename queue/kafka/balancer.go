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
	"fmt"
	"sync/atomic"
)

type Balancer[T any] interface {
	Name() string
	GetNext() (T, error)
}

type RRBalance[T any] struct {
	cnt     int64
	objects []T
}

func (rr *RRBalance[T]) Name() string {
	return "round-robin"
}

func (rr *RRBalance[T]) GetNext() (T, error) {
	var t T
	if len(rr.objects) == 0 {
		return t, fmt.Errorf("没有可选择的对象")
	}
	cnt := atomic.AddInt64(&rr.cnt, 1)
	index := cnt % int64(len(rr.objects))
	return rr.objects[index], nil
}

type BalanceBuilder[T any] struct {
	name    string
	objects []T
}

func NewBalanceBuilder[T any](name string, objects []T) *BalanceBuilder[T] {
	return &BalanceBuilder[T]{
		name:    name,
		objects: objects,
	}
}

func (bb *BalanceBuilder[T]) Build(name string) Balancer[T] {
	switch name {
	case "round-robin":
		return &RRBalance[T]{cnt: -1, objects: bb.objects}
	default:
		return nil
	}
}
