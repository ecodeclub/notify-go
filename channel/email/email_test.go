package email

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/ecodeclub/notify-go/pkg/notifier"
)

func TestSend(t *testing.T) {
	e := NewEmailChannel(Config{
		SenderAddress: "Hooko <hooko_1@cooode.fun>",
		SmtpHostAddr:  "gz-smtp.qcloudmail.com:465",
		SmtpUserName:  "hooko_1@cooode.fun",
		SmtpPwd:       "xxx",
	})
	content := Content{
		Subject: "发送主题-测试",
		//Cc:      []string{"648646891@qq.com"},
		Html: "<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"utf-8\">\n<title>hello world</title>\n</head>\n<body>\n " +
			"<h1>我的第一个标题</h1>\n    <p>我的第一个段落。</p>\n</body>\n</html>",
	}

	cb, _ := json.Marshal(content)
	deli := notifier.Delivery{
		Receivers: []notifier.Receiver{
			{Email: "648646891@qq.com"},
		},
		Content: cb,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := e.Execute(ctx, deli)
	t.Log(err)
}
