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
