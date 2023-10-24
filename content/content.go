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

package content

import (
	"context"

	"github.com/ecodeclub/notify-go/pkg/notifier"
	"github.com/ecodeclub/notify-go/store/mysql"
)

/*
content服务 根据模版组装发送内容
*/

type Service struct {
	tDAO mysql.ITemplateDAO
}

type IContentService interface {
	GetContent(ctx context.Context, receivers []notifier.Receiver, templateId uint64,
		variable map[string]interface{}) (notifier.Content, error)
}

func NewContentService(td mysql.ITemplateDAO) IContentService {
	return &Service{
		tDAO: td,
	}
}

func (s *Service) GetContent(ctx context.Context, receivers []notifier.Receiver, templateId uint64,
	variable map[string]interface{}) (notifier.Content, error) {
	var cont notifier.Content

	tpl, err := s.tDAO.GetTContent(templateId, "")
	if err != nil {
		return cont, err
	}

	// 通过target获取该target的特定内容

	// 通过模版渲染出发送内容
	cont, err = s.renderContent(ctx, tpl, variable)

	return cont, err
}

func (s *Service) renderContent(ctx context.Context, tpl string,
	variable map[string]interface{}) (notifier.Content, error) {
	var cont notifier.Content
	return cont, nil
}
