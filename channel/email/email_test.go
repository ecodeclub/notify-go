package email

import (
	"testing"
)

func TestSend(t *testing.T) {
	//e := NewEmailChannel(Config{
	//	SenderAddress: "Hooko <hooko_1@cooode.fun>",
	//	SmtpHostAddr:  "gz-smtp.qcloudmail.com:465",
	//	SmtpUserName:  "hooko_1@cooode.fun",
	//	SmtpPwd:       "xxx",
	//})
	//
	//deli := notifier.Delivery{
	//	Receivers: []notifier.Receiver{
	//		{Email: "648646891@qq.com"},
	//	},
	//	Content: notifier.Content{
	//		Title: "发送主题-测试",
	//		Data: []byte("<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"utf-8\">\n<title>hello world</title>\n</head>\n<body>\n " +
	//			"<h1>我的第一个标题</h1>\n    <p>我的第一个段落。</p>\n</body>\n</html>"),
	//	},
	//}
	//ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//defer cancel()
	//err := e.Execute(ctx, deli)
	//t.Log(err)
}
