package model

type Role struct {
	tableName struct{} `sql:"roles" pg:",discard_unknown_columns"`

	Id    int64   `sql:"id"  json:"id"`
	Name  string  `sql:"name" json:"name"`
	Users []*User `pg:"many2many:users_roles,joinFK:role_id" json:"users"`
}
