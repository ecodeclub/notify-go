package mysql

import (
	"context"
	"time"

	"github.com/ecodeclub/notify-go/internal/pkg/task"
	"xorm.io/xorm"
)

type notifyGoDAO struct {
	engine *xorm.Engine
}

type INotifyGoDAO interface {
	InsertRecord(ctx context.Context, templateId int64, target task.Receiver, msgContent string) error
}

type ITemplateDAO interface {
	GetTContent(templateId uint64, country string) (string, error)
}

func NewINotifyGoDAO(e *xorm.Engine) INotifyGoDAO {
	return &notifyGoDAO{e}
}

func NewITemplateDAO(e *xorm.Engine) ITemplateDAO {
	return &notifyGoDAO{e}
}

func (n *notifyGoDAO) InsertRecord(ctx context.Context, templateId int64, target task.Receiver,
	msgContent string) error {
	sess := n.engine.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	delivery := Delivery{
		TemplateId:  templateId,
		Status:      1, // 消息创建状态
		SendChannel: 40,
		MsgType:     20,
		Proposer:    "crm",
		Creator:     "chenhaokun",
		Updator:     "chenhaokun",
		IsDelted:    0,
		Created:     time.Now(),
		Updated:     time.Now(),
	}

	// 这里配置好struct的id自增tag，会自动赋值插入的id
	if _, err := n.engine.Insert(&delivery); err != nil {
		return err
	}

	tgt := Target{
		TargetIdType: target.Type(),
		TargetId:     target.Value(),
		DeliveryId:   delivery.Id,
		Status:       1, // 创建状态
		MsgContent:   msgContent,
	}
	if _, err := n.engine.Insert(&tgt); err != nil {
		return err
	}

	return sess.Commit()
}

func (n *notifyGoDAO) GetTContent(templateId uint64, country string) (string, error) {
	tpl := Template{}
	has, err := n.engine.Where("id = ?", templateId).Get(&tpl)
	if err != nil || !has {
		return "", err
	}
	// TODO 根据国家返回模版内容
	return tpl.ChsContent, nil
}
