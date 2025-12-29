package repository

import (
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/common/domain/entity"
	"github.com/budimanlai/go-core/common/models"
)

type MessagingTemplateRepository interface {
	base.BaseRepository[entity.MessagingTemplate, models.MessagingTemplate]
}
