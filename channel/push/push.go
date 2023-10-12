package push

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/ecodeclub/notify-go/pkg/logger"
	"github.com/ecodeclub/notify-go/pkg/notifier"
	"github.com/ecodeclub/notify-go/pkg/ral"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type Config struct {
	AppKey       string `json:"app_key"`
	MasterSecret string `json:"master_secret"`
	AppId        string `json:"app_id"`
}

type ChannelPushImpl struct {
	config Config
	client ral.Client
}

// Content 个推的请求参数
type Content struct {
	RequestID   string      `json:"request_id"`
	Settings    Settings    `json:"settings"`
	Audience    Audience    `json:"audience"`
	PushMessage PushMessage `json:"push_message"`
}

type Settings struct {
	TTL int `json:"ttl"`
}

type Audience struct {
	Cid []string `json:"cid"`
}

type Notification struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	ClickType string `json:"click_type"`
	URL       string `json:"url"`
}

type PushMessage struct {
	Notification Notification `json:"notification"`
}

type Result struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data map[string]string `json:"data"`
}

func NewPushChannel(c Config) *ChannelPushImpl {
	pc := &ChannelPushImpl{
		client: ral.NewClient("getui"),
		config: c,
	}
	return pc
}

func (pc *ChannelPushImpl) Execute(ctx context.Context, deli notifier.Delivery) error {
	token, err := pc.getToken()
	if err != nil {
		return err
	}

	var content Content
	err = json.Unmarshal(deli.Content, &content)
	if err != nil {
		return err
	}

	if ctx.Value("req_id") != nil {
		content.RequestID = ctx.Value("req_id").(string)
	}

	userIds := slice.Map[notifier.Receiver, string](deli.Receivers, func(idx int, recv notifier.Receiver) string {
		return recv.UserId
	})
	content.Audience.Cid = append(content.Audience.Cid, userIds...)

	req := ral.Request{
		Header: map[string]string{
			"content-type": "application/json;charset=utf-8",
			"token":        token,
		},
		PathParams: map[string]string{"app_id": pc.config.AppId},
		Body:       content,
	}

	var resp map[string]any
	err = pc.client.Ral(context.TODO(), "Send", req, &resp, map[string]any{})

	logger.Default().Info("push send done", logger.Any("resp", resp))
	return err
}

func (pc *ChannelPushImpl) Name() string {
	return "push"
}

func (pc *ChannelPushImpl) getToken() (token string, err error) {
	ts, sign := pc.getSign()
	req := ral.Request{
		Header: map[string]string{"content-type": "application/json;charset=utf-8"},
		Body: map[string]interface{}{
			"sign":      sign,
			"timestamp": ts,
			"appkey":    pc.config.AppKey,
		},
		PathParams: map[string]string{"app_id": pc.config.AppId},
	}

	var respSucc Result
	err = pc.client.Ral(context.TODO(), "Auth", req, &respSucc, map[string]any{})
	if err != nil {
		return
	}
	var ok bool
	token, ok = respSucc.Data["token"]

	if !ok {
		err = errors.New("[push] 获取token失败")
	}
	return
}

func (pc *ChannelPushImpl) getSign() (timestamp string, sign string) {
	timestamp = strconv.FormatInt(time.Now().UnixMilli(), 10)
	dataToSign := pc.config.AppKey + timestamp + pc.config.MasterSecret

	// 计算SHA-256哈希
	sha256Hash := sha256.Sum256([]byte(dataToSign))

	// 将哈希结果转换为十六进制字符串
	sign = fmt.Sprintf("%x", sha256Hash)

	return
}
