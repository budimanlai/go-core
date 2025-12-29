package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/common/domain/entity"
	"github.com/budimanlai/go-core/common/domain/repository"
	"github.com/budimanlai/go-core/common/domain/usecase"
	"github.com/budimanlai/go-core/common/dto"
	"gorm.io/gorm"
)

type MessagingTemplateUsecaseImpl struct {
	base.BaseUsecase[entity.MessagingTemplate]
}

func NewMessagingTemplateUsecaseImpl(db *gorm.DB, repo repository.MessagingTemplateRepository) usecase.MessagingTemplateUsecase {
	return &MessagingTemplateUsecaseImpl{
		BaseUsecase: base.NewBaseUsecase(repo, db),
	}
}

// RenderTemplate renders a messaging template with the provided data
func (u *MessagingTemplateUsecaseImpl) RenderTemplate(ctx context.Context, channel, templateName string, data map[string]interface{}) (*dto.MessageTemplateDTO, error) {
	tpl, err := u.FindOne(ctx, func(d *gorm.DB) *gorm.DB {
		return d.Where("template_name = ? and channel = ?", templateName, channel)
	})
	if err != nil {
		return nil, err
	}
	if tpl == nil {
		return nil, errors.New("template not found")
	}

	// append additional data if needed
	data["datetime"] = time.Now().Format("2006-01-02 15:04:05")

	var out dto.MessageTemplateDTO = dto.MessageTemplateDTO{
		Channel:      channel,
		TemplateName: templateName,
		Subject:      u.renderContent(tpl.Subject, data),
		ContentHtml:  u.renderContent(tpl.ContentHTML, data),
		ContentText:  u.renderContent(tpl.ContentText, data),
	}

	return &out, nil
}

func (u *MessagingTemplateUsecaseImpl) renderContent(content string, data map[string]interface{}) string {
	// simple placeholder replacement
	for key, value := range data {
		placeholder := "{{" + key + "}}"
		strValue := ""
		switch v := value.(type) {
		case string:
			strValue = v
		case int:
			strValue = fmt.Sprintf("%d", v)
		case float64:
			strValue = fmt.Sprintf("%f", v)
		default:
			strValue = fmt.Sprintf("%v", v)
		}
		content = strings.ReplaceAll(content, placeholder, strValue)
	}
	return content
}
