package model

type Permission struct {
	tableName struct{} `sql:"permissions" pg:",discard_unknown_columns"`

	Id          int64  `sql:"id"  json:"id"`
	Url         string `sql:"url" json:"url"`
	HttpMethod  string `sql:"http_method" json:"httpMethod"`
	Description string `sql:"description" json:"description"`
}
