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

package mysql

import (
	"context"
	"time"

	"github.com/ecodeclub/notify-go/pkg/notifier"
	"xorm.io/xorm"
)

type notifyGoDAO struct {
	engine *xorm.Engine
}

type INotifyGoDAO interface {
	InsertRecord(ctx context.Context, templateId int64, target notifier.Receiver, msgContent string) error
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

func (n *notifyGoDAO) InsertRecord(ctx context.Context, templateId int64, target notifier.Receiver,
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
		DeliveryId: delivery.Id,
		Status:     1, // 创建状态
		MsgContent: msgContent,
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
