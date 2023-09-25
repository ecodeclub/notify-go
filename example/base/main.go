package main

import (
	"context"
	"log"
	"time"

	notify_go "github.com/ecodeclub/notify-go"
	"github.com/ecodeclub/notify-go/channel"
	"github.com/ecodeclub/notify-go/channel/email"
	"github.com/ecodeclub/notify-go/queue/kafka"

	"github.com/BurntSushi/toml"

	"github.com/ecodeclub/notify-go/content"
	"github.com/ecodeclub/notify-go/store/mysql"
	"github.com/ecodeclub/notify-go/target"
)

var (
	engine   = mysql.NewEngine(mysql.DBConfig{DriverName: "", Dsn: ""})
	kafkaCfg = kafka.Config{}
)

func serve() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	// 读取kafka配置
	_, err := toml.DecodeFile("./conf/kafka.toml", &kafkaCfg)
	if err != nil {
		log.Fatal("kafka配置读取失败")
	}

	// 启动邮件发送的消费者
	qSrv := kafka.NewKafka(kafkaCfg)
	qSrv.Consume(ctx, email.NewEmailChannel(email.Config{}))
}

func main() {
	// 邮件消息发送服务
	go serve()

	// 通过target服务获取发送对象
	targetSrv := target.NewTargetService()
	receivers := targetSrv.GetTarget(context.TODO(), 123)

	// 通过content服务获取发送内容
	contentSrv := content.NewContentService(mysql.NewITemplateDAO(engine))
	msg, _ := contentSrv.GetContent(context.TODO(), receivers, 123, nil)

	// 创建异步邮件发送队列
	q := kafka.NewKafka(kafkaCfg)
	asyncChannel := channel.AsyncChannel{Queue: q, IChannel: email.NewEmailChannel(email.Config{})}

	// 执行普通发送
	n := notify_go.NewNotification(asyncChannel, receivers, msg)
	_ = n.Send(context.TODO())

	// 定时任务发送
	task := notify_go.NewTriggerTask(n, time.Now().Add(time.Minute))
	task.Send(context.TODO())
	<-task.Err
}
