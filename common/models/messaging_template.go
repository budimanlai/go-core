package models

type MessagingTemplate struct {
	TemplateName string `gorm:"column:template_name;type:varchar(50)"`
	Channel      string `gorm:"column:channel;type:varchar(15)"`
	Subject      string `gorm:"column:subject;type:varchar(256)"`
	ContentHTML  string `gorm:"column:content_html;type:text"`
	ContentText  string `gorm:"column:content_text;type:text"`
}

func (MessagingTemplate) TableName() string {
	return "messaging_template"
}
