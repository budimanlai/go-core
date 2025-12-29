package repository

import (
	"github.com/budimanlai/go-core/base"
	"github.com/budimanlai/go-core/common/domain/entity"
	"github.com/budimanlai/go-core/common/domain/repository"
	"github.com/budimanlai/go-core/common/models"
)

type MessagingTemplateRepositoryImpl struct {
	base.BaseRepository[entity.MessagingTemplate, models.MessagingTemplate]
}

func NewMessagingTemplateRepositoryImpl(f *base.Factory) repository.MessagingTemplateRepository {
	return &MessagingTemplateRepositoryImpl{
		BaseRepository: base.NewRepository[entity.MessagingTemplate, models.MessagingTemplate](f),
	}
}
