package model

import "time"

type User struct {
	Id                 int64      `gorm:"primaryKey;column:id" json:"id"`
	UserName           string     `gorm:"column:username;size:255;not null" json:"userName"`
	Email              string     `gorm:"column:email;size:255;not null;unique" json:"email"`
	Password           []byte     `gorm:"column:password;not null" json:"-"`
	PhoneNumber        string     `gorm:"column:phone_number;size:20" json:"phoneNumber"`
	AcceptNotification bool       `gorm:"column:accept_notification;default:false" json:"acceptNotification,omitempty"`
	IsActive           bool       `gorm:"column:is_active;default:false" json:"isActive"`
	InactivatedDate    *time.Time `gorm:"column:inactivated_date" json:"inactivatedDate,omitempty"`
	FullName           string     `gorm:"column:full_name;size:255" json:"fullName"`
	Description        string     `gorm:"column:description;type:text" json:"description"`
	Roles              []*Role    `gorm:"many2many:users_roles;joinForeignKey:UserID;joinReferences:RoleID" json:"roles"`
	CreatedAt          time.Time  `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt          time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"-"`
}

// TableName overrides the default table name
func (User) TableName() string {
	return "users"
}

type UserRole struct {
	UserId int64 `gorm:"column:user_id;not null" json:"userId"`
	RoleId int64 `gorm:"column:role_id;not null" json:"roleId"`
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
