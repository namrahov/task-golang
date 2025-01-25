package model

type Role struct {
	Id          int64         `gorm:"primaryKey;column:id" json:"id"`
	Name        string        `gorm:"column:name;size:32;not null;unique" json:"name"`
	Permissions []*Permission `gorm:"many2many:roles_permissions;joinForeignKey:RoleID;joinReferences:PermissionID" json:"permissions"`
}

// TableName overrides the default table name
func (Role) TableName() string {
	return "roles"
}
