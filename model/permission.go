package model

type Permission struct {
	ID          int64  `gorm:"primaryKey;column:id" json:"id"`
	URL         string `gorm:"column:url;size:255;not null" json:"url"`
	HTTPMethod  string `gorm:"column:http_method;size:10;not null" json:"httpMethod"`
	Description string `gorm:"column:description;type:text" json:"description"`
}

// TableName overrides the default table name
func (Permission) TableName() string {
	return "permissions"
}
