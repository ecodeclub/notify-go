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

package notify_go

import (
	"context"

	"github.com/ecodeclub/notify-go/pkg/notifier"
	"github.com/pborman/uuid"
)

type Notification struct {
	notifier.Delivery
	Channel notifier.IChannel
}

type ChannelFunc func(ctx context.Context, no *Notification) error

type Middleware func(channelFunc ChannelFunc) ChannelFunc

func (no *Notification) Send(ctx context.Context, mls ...Middleware) error {
	var root ChannelFunc = func(ctx context.Context, no *Notification) error {
		return no.Channel.Execute(ctx, no.Delivery)
	}

	for i := len(mls) - 1; i > 0; i-- {
		root = mls[i](root)
	}

	return root(ctx, no)
}

func NewNotification(c notifier.IChannel, recvs []notifier.Receiver, content notifier.Content) *Notification {
	no := &Notification{
		Channel: c,
		Delivery: notifier.Delivery{
			DeliveryID: uuid.NewUUID().String(),
			Receivers:  recvs,
			Content:    content,
		},
	}
	return no
}
