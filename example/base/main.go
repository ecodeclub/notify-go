package main

import (
	"context"
	"github.com/BurntSushi/toml"
	"log"

	"github.com/ecodeclub/notify-go/internal/content"
	"github.com/ecodeclub/notify-go/internal/pkg/mq"
	"github.com/ecodeclub/notify-go/internal/pkg/task"
	"github.com/ecodeclub/notify-go/internal/send"
	"github.com/ecodeclub/notify-go/internal/send/sender"
	"github.com/ecodeclub/notify-go/internal/store/mysql"
	"github.com/ecodeclub/notify-go/internal/target"
)

var (
	engine   = mysql.NewEngine(mysql.DBConfig{DriverName: "", Dsn: ""})
	kafkaCfg = mq.KafkaConfig{}
)

func serve() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	// 创建邮件发送handler, 并包装成消息消费者任务
	emailExecutor := sender.NewEmailHandler(sender.EmailConfig{Addr: "", Auth: nil})
	t := task.NewTask(emailExecutor)

	// 读取kafka配置
	_, err := toml.DecodeFile("./conf/kafka.toml", &kafkaCfg)
	if err != nil {
		log.Fatal("kafka配置读取失败")
	}

	// 启动邮件发送的消费者
	qSrv := mq.NewQueueService(kafkaCfg)
	qSrv.Consume(ctx, "email", &t)
}

func main() {
	// 邮件消息发送服务
	go serve()

	// 发送邮件
	sendSrv := send.NewSendService(
		mq.NewQueueService(kafkaCfg),
		mysql.NewINotifyGoDAO(engine),
	)

	// 通过target服务获取发送对象
	targetSrv := target.NewTargetService()
	receivers := targetSrv.GetTarget(context.TODO(), 123)

	// 通过content服务获取发送内容
	contentSrv := content.NewContentService(mysql.NewITemplateDAO(engine))
	msg, _ := contentSrv.GetContent(context.TODO(), nil, 123, nil)

	// 执行发送
	_ = sendSrv.Send(context.TODO(), receivers, msg)
}
