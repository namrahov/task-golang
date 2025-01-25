package repo

import (
	"fmt"
	"github.com/go-pg/pg"
	"task-golang/model"
)

type IUserRepo interface {
	BeginTransaction() (*pg.Tx, error)
	GetUserByEmail(email string) (*model.User, error)
	SaveUser(tx *pg.Tx, user *model.User) (*model.User, error)
	AddRolesToUser(tx *pg.Tx, userId int64, roles []*model.Role) error
	FindUserById(id int64) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	FindActiveUserByEmailOrUsername(EmailOrNickname string) (*model.User, error)
	//GetPermissions(roles []string) ([]*model.Permission, error)
}

type UserRepo struct {
}

// BeginTransaction starts a database transaction and returns the transaction object.
func (r *UserRepo) BeginTransaction() (*pg.Tx, error) {
	tx, err := Db.Begin()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (r UserRepo) FindUserById(id int64) (*model.User, error) {
	var user model.User
	err := Db.Model(&user).
		Where("id = ?", id).
		Select()

	if err != nil {
		// Check if no rows were found
		if err == pg.ErrNoRows {
			// Return nil user and no error
			return nil, nil
		}
		// Return any other error
		return nil, err
	}

	return &user, nil
}

func (r UserRepo) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := Db.Model(&user).
		Where("email = ?", email).
		Limit(1).
		Select()

	if err != nil {
		// Check if no rows were found
		if err == pg.ErrNoRows {
			// Return nil user and no error
			return nil, nil
		}
		// Return any other error
		return nil, err
	}

	return &user, nil
}

func (r UserRepo) SaveUser(tx *pg.Tx, user *model.User) (*model.User, error) {
	_, err := tx.Model(user).Insert()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r UserRepo) UpdateUser(user *model.User) (*model.User, error) {
	_, err := Db.Model(user).
		OnConflict("(id) DO UPDATE").
		Insert()
	if err != nil {
		fmt.Println("err=", err)
		return nil, err
	}
	return user, nil
}

func (r UserRepo) AddRolesToUser(tx *pg.Tx, userId int64, roles []*model.Role) error {
	userRoles := make([]*model.UserRole, len(roles))
	for i, role := range roles {
		userRoles[i] = &model.UserRole{
			UserId: userId,
			RoleId: role.Id,
		}
	}
	_, err := tx.Model(&userRoles).Insert()
	return err
}

func (r UserRepo) FindActiveUserByEmailOrUsername(emailOrNickname string) (*model.User, error) {
	var user model.User
	err := Db.Model(&user).
		Where("(email = ? OR username = ?) AND is_active = ?", emailOrNickname, emailOrNickname, true).
		Select()

	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err // Other errors
	}
	// Fetch roles explicitly to avoid duplicates
	var roles []*model.Role
	err = Db.Model(&roles).
		Table("roles"). // Explicitly set the table name
		Join("JOIN users_roles ur ON ur.role_id = roles.id").
		Where("ur.user_id = ?", user.Id).
		Select()

	if err != nil {
		fmt.Println("3333333", err)
		return nil, err
	}

	// Assign unique roles to the user
	roleMap := make(map[int64]*model.Role)
	for _, role := range roles {
		roleMap[role.Id] = role
	}

	user.Roles = make([]*model.Role, 0, len(roleMap))
	for _, role := range roleMap {
		user.Roles = append(user.Roles, role)
	}

	return &user, nil
}

//func (r UserRepo) GetPermissions(roles []string) ([]*model.Permission, error) {
//	// Fetch permissions for the user's roles
//	var permissions []*model.Permission
//	err := Db.Model(&model.Permission{}).
//		Join("JOIN roles_permissions rp ON rp.permission_id = permissions.id").
//		Join("JOIN roles r ON r.id = rp.role_id").
//		Where("r.name IN (?)", roles).
//		Select(&permissions)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return permissions, nil
//}
