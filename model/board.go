package model

import "time"

type Board struct {
	Id        int64     `gorm:"primaryKey;column:id" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	CreatedBy string    `gorm:"column:created_by" json:"createdBy"`
	Users     []*User   `gorm:"many2many:users_boards;joinForeignKey:BoardId;joinReferences:UserId" json:"users"`
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
}

// TableName overrides the default table name
func (Board) TableName() string {
	return "boards"
}
