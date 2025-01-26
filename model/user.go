package model

import "time"

type User struct {
	Id                 int64      `gorm:"primaryKey;column:id" json:"id"`
	UserName           string     `gorm:"column:username" json:"userName"`
	Email              string     `gorm:"column:email" json:"email"`
	Password           []byte     `gorm:"column:password" json:"-"`
	PhoneNumber        string     `gorm:"column:phone_number" json:"phoneNumber"`
	AcceptNotification bool       `gorm:"column:accept_notification" json:"acceptNotification,omitempty"`
	IsActive           bool       `gorm:"column:is_active" json:"isActive"`
	InactivatedDate    *time.Time `gorm:"column:inactivated_date" json:"inactivatedDate,omitempty"`
	FullName           string     `gorm:"column:full_name" json:"fullName"`
	Description        string     `gorm:"column:description" json:"description"`
	Roles              []*Role    `gorm:"many2many:users_roles;joinForeignKey:UserId;joinReferences:RoleId" json:"roles"`
	Boards             []*Board   `gorm:"many2many:users_boards;joinForeignKey:UserId;joinReferences:BoardId" json:"boards"`
	CreatedAt          time.Time  `gorm:"column:created_at" json:"-"`
	UpdatedAt          time.Time  `gorm:"column:updated_at" json:"-"`
}

// TableName overrides the default table name
func (User) TableName() string {
	return "users"
}

type UserRole struct {
	UserId int64 `gorm:"column:user_id" json:"userId"`
	RoleId int64 `gorm:"column:role_id" json:"roleId"`
}

// TableName overrides the default table name
func (UserRole) TableName() string {
	return "users_roles"
}

type UserRegistrationDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthRequestDto struct {
	EmailOrNickname string `json:"emailOrNickname"`
	Password        string `json:"password"`
	RememberMe      bool   `json:"rememberMe"`
}

type JwtToken struct {
	Token string `json:"token"`
}
