package repo

import (
	"task-golang/model"
)

type IUserRepo interface {
	GetUserByEmail(email string) (*model.User, error)
	SaveUser(user *model.User) (*model.User, error)
}

type UserRepo struct {
}

func (r UserRepo) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := Db.Model(&user).
		Where("email = ?", email).
		Limit(1).
		Select()

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r UserRepo) SaveUser(user *model.User) (*model.User, error) {
	tx, err := Db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Model(user).Insert()
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r UserRepo) AddRolesToUser(userId int64, roles []*model.Role) error {
	userRoles := make([]*model.UserRole, len(roles))
	for i, role := range roles {
		userRoles[i] = &model.UserRole{
			UserId: userId,
			RoleId: role.Id,
		}
	}
	_, err := Db.Model(&userRoles).Insert()
	return err
}
