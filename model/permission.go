package model

type Permission struct {
	ID          int64  `gorm:"primaryKey;column:id" json:"id"`
	URL         string `gorm:"column:url" json:"url"`
	HTTPMethod  string `gorm:"column:http_method" json:"httpMethod"`
	Description string `gorm:"column:description" json:"description"`
}

// TableName overrides the default table name
func (Permission) TableName() string {
	return "permissions"
}
