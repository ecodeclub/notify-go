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

package main

import (
	"context"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ecodeclub/ekit/slice"
	notifygo "github.com/ecodeclub/notify-go"
	"github.com/ecodeclub/notify-go/channel"
	"github.com/ecodeclub/notify-go/channel/push"
	"github.com/ecodeclub/notify-go/pkg/notifier"
	"github.com/ecodeclub/notify-go/pkg/ral"
	"github.com/ecodeclub/notify-go/queue"
	"github.com/ecodeclub/notify-go/queue/kafka"
)

var (
	kafkaCfg   kafka.Config
	getui      ral.Resource
	pushConfig push.Config
)

func init() {
	// kafka配置
	_, _ = toml.DecodeFile("./example/base/conf/kafka.toml", &kafkaCfg)

	getui, _ = slice.Find[ral.Resource](
		ral.NewService("./example/base/conf/ral.toml").Resources,
		func(src ral.Resource) bool {
			return src.Name == "getui"
		})

	pushConfig = push.Config{}

	go serve()
	<-time.After(2 * time.Second)
}

func serve() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	qSrv := kafka.NewKafka(kafkaCfg)
	qSrv.Consume(ctx, push.NewPushChannel(pushConfig, ral.NewClient(getui)))
}

func main() {
	// 后续支持通过target服务获取发送对象
	receivers := []notifier.Receiver{
		{
			Email:  "xxx@qq.com",
			Phone:  "+8613111111111",
			UserId: "2a4b74c682210299781c4a1b9b308c5e",
		},
	}

	// 后续支持通过content服务获取发送内容
	msg := notifier.Content{
		Title:     "测试一下一下下titile",
		Data:      []byte("测试一下一下下content"),
		ClickType: "none",
	}

	// 同步push发送
	sendPushSync(receivers, msg)

	// 异步push发送
	q := kafka.NewKafka(kafkaCfg)
	sendPushAsync(q, receivers, msg)

	// 定时任务发送
	triggerSend(q, receivers, msg)
}

func sendPushSync(recvs []notifier.Receiver, msg notifier.Content) {
	// 同步发送
	syncChannel := channel.SyncChannel{IChannel: push.NewPushChannel(pushConfig, ral.NewClient(getui))}

	n := notifygo.NewNotification(syncChannel, recvs, msg)
	err := n.Send(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}

func sendPushAsync(q queue.IQueue, recvs []notifier.Receiver, msg notifier.Content) {
	// 异步发送
	asyncChannel := channel.AsyncChannel{Queue: q, IChannel: push.NewPushChannel(pushConfig, ral.NewClient(getui))}
	n := notifygo.NewNotification(asyncChannel, recvs, msg)
	err := n.Send(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}

func triggerSend(q queue.IQueue, recvs []notifier.Receiver, msg notifier.Content) {
	// 异步发送
	asyncChannel := channel.AsyncChannel{Queue: q, IChannel: push.NewPushChannel(pushConfig, ral.NewClient(getui))}
	n := notifygo.NewNotification(asyncChannel, recvs, msg)

	task := notifygo.NewTriggerTask(n, time.Now().Add(time.Minute))
	task.Send(context.TODO())
	<-task.Err
}
