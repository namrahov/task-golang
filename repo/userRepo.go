package repo

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"task-golang/model"
)

type IUserRepo interface {
	FindUserById(id int64) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	SaveUser(tx *gorm.DB, user *model.User) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	AddRolesToUser(tx *gorm.DB, userId int64, roles []*model.Role) error
	FindActiveUserByEmailOrUsername(emailOrNickname string) (*model.User, error)
	//BeginTransaction() *gorm.DB
}

type UserRepo struct {
}

func (r UserRepo) FindUserById(id int64) (*model.User, error) {
	var user model.User
	err := Db.First(&user, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Return nil user and no error if the record is not found
		return nil, nil
	}

	if err != nil {
		// Return any other error
		return nil, err
	}

	return &user, nil
}

func (r UserRepo) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := Db.Where("email = ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Return nil user and no error if the record is not found
		return nil, nil
	}

	if err != nil {
		// Return any other error
		return nil, err
	}

	return &user, nil
}

func (r UserRepo) SaveUser(tx *gorm.DB, user *model.User) (*model.User, error) {
	err := tx.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r UserRepo) UpdateUser(user *model.User) (*model.User, error) {
	err := Db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}}, // Specifies the column to check for conflict
			DoUpdates: clause.AssignmentColumns([]string{"username", "email", "password", "phone_number", "accept_notification", "is_active", "inactivated_date", "full_name", "description", "updated_at"}),
		},
	).Create(user).Error

	if err != nil {
		fmt.Println("err=", err)
		return nil, err
	}
	return user, nil
}

func (r UserRepo) AddRolesToUser(tx *gorm.DB, userId int64, roles []*model.Role) error {
	userRoles := make([]model.UserRole, len(roles))
	for i, role := range roles {
		userRoles[i] = model.UserRole{
			UserId: userId,
			RoleId: role.Id,
		}
	}

	// Batch insert userRoles
	err := tx.Create(&userRoles).Error
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepo) FindActiveUserByEmailOrUsername(emailOrNickname string) (*model.User, error) {
	var user model.User

	err := Db.Model(&model.User{}).
		Where("(email = ? OR username = ?) AND is_active = ?", emailOrNickname, emailOrNickname, true).
		Preload("Roles", func(db *gorm.DB) *gorm.DB {
			return db.Joins("JOIN users_roles ur ON ur.role_id = roles.id").
				Joins("JOIN users ON users.id = ur.user_id") // Explicitly join the `users` table
		}).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, err // Other errors
	}

	return &user, nil
}
