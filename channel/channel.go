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

package channel

import (
	"context"

	"github.com/ecodeclub/notify-go/pkg/notifier"
	"github.com/ecodeclub/notify-go/queue"
)

type SyncChannel struct {
	notifier.IChannel
}

type AsyncChannel struct {
	Queue queue.IQueue
	notifier.IChannel
}

func (s SyncChannel) Execute(ctx context.Context, deli notifier.Delivery) error {
	err := s.IChannel.Execute(ctx, deli)
	return err
}

func (ac AsyncChannel) Execute(ctx context.Context, deli notifier.Delivery) error {
	// 提前启动 channel 对应的消费者
	// 发送的时候，发送具体的 sender函数 和 参数
	err := ac.Queue.Produce(ctx, ac.IChannel, deli)
	return err
}
