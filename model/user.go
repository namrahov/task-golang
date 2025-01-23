package model

type User struct {
	tableName struct{} `sql:"users" pg:",discard_unknown_columns"`

	Id                 int64   `sql:"id"  json:"id"`
	UserName           string  `sql:"username" json:"userName"`
	Email              string  `sql:"email" json:"email"`
	Password           []byte  `sql:"password" json:"-"`
	PhoneNumber        string  `sql:"phone_number" json:"phoneNumber"`
	AcceptNotification bool    `sql:"accept_notification" json:"acceptNotification,omitempty"`
	IsActive           bool    `sql:"is_active" json:"isActive"`
	InactivatedDate    string  `sql:"inactivated_date" json:"inactivatedDate,omitempty"`
	FullName           string  `sql:"full_name" json:"fullName"`
	Description        string  `sql:"description" json:"description"`
	Roles              []*Role `pg:"many2many:users_roles,joinFK:user_id" json:"roles"`
	CreatedAt          string  `sql:"created_at" json:"-"`
	UpdatedAt          string  `sql:"updated_at" json:"-"`
}

type UserRole struct {
	tableName struct{} `sql:"users_roles" pg:",discard_unknown_columns"`

	UserId int64 `sql:"user_id"`
	RoleId int64 `sql:"role_id"`
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
