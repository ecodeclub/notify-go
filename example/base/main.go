package main

import (
	"context"
	"github.com/ecodeclub/notify-go/internal"
	"github.com/ecodeclub/notify-go/internal/queue"
	"log"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/ecodeclub/notify-go/internal/content"
	"github.com/ecodeclub/notify-go/internal/store/mysql"
	"github.com/ecodeclub/notify-go/internal/target"
)

var (
	engine   = mysql.NewEngine(mysql.DBConfig{DriverName: "", Dsn: ""})
	kafkaCfg = queue.KafkaConfig{}
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
	qSrv := queue.NewQueueService(kafkaCfg)
	qSrv.Consume(ctx, "email")
}

func main() {
	// 邮件消息发送服务
	go serve()

	// 通过target服务获取发送对象
	targetSrv := target.NewTargetService()
	receivers := targetSrv.GetTarget(context.TODO(), 123)

	// 通过content服务获取发送内容
	contentSrv := content.NewContentService(mysql.NewITemplateDAO(engine))
	msg, _ := contentSrv.GetContent(context.TODO(), receivers[0], 123, nil)

	qSrv := queue.NewQueueService(kafkaCfg)
	channel := internal.Channel{Queue: qSrv}

	// 执行普通发送
	n := internal.NewNotification(channel, receivers, msg)
	_ = n.Send(context.TODO())

	// 定时任务发送
	task := internal.NewTriggerTask(n, time.Now().Add(time.Minute))
	task.Send(context.TODO())
	<-task.Err
}
