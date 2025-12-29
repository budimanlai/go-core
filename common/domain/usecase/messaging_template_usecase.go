package usecase

import (
	"context"

	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/common/domain/entity"
	"github.com/budimanlai/go-core/common/dto"
)

type MessagingTemplateUsecase interface {
	base.BaseUsecase[entity.MessagingTemplate]

	// RenderTemplate renders a messaging template with the provided data
	RenderTemplate(ctx context.Context, channel, templateName string, data map[string]interface{}) (*dto.MessageTemplateDTO, error)
}
