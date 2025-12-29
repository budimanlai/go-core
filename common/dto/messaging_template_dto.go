package dto

type MessageTemplateDTO struct {
	Channel      string
	TemplateName string
	Subject      string
	ContentHtml  string
	ContentText  string
}
