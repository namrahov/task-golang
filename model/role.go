package model

type Role struct {
	tableName struct{} `sql:"roles" pg:",discard_unknown_columns"`

	Id          int64         `sql:"id"  json:"id"`
	Name        string        `sql:"name" json:"name"`
	Permissions []*Permission `pg:"many2many:roles_permissions,joinFK:role_id" json:"permissions"`
}
