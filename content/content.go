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
		variable map[string]interface{}) ([]byte, error)
}

func NewContentService(td mysql.ITemplateDAO) IContentService {
	return &Service{
		tDAO: td,
	}
}

func (s *Service) GetContent(ctx context.Context, receivers []notifier.Receiver, templateId uint64,
	variable map[string]interface{}) ([]byte, error) {
	var cont []byte

	tpl, err := s.tDAO.GetTContent(templateId, "")
	if err != nil {
		return cont, err
	}

	// 通过target获取该target的特定内容

	// 通过模版渲染出发送内容
	cont, err = s.renderContent(ctx, tpl, variable)

	return cont, nil
}

func (s *Service) renderContent(ctx context.Context, tpl string,
	variable map[string]interface{}) ([]byte, error) {
	var cont []byte
	return cont, nil
}
