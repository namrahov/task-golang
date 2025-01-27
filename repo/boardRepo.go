package repo

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"task-golang/model"
)

type IBoardRepo interface {
	SaveBoard(board *model.Board) (*model.Board, error)
	SaveUserBoard(ctx context.Context, userId int64, boardId int64) error
}

type BoardRepo struct {
}

func (r BoardRepo) SaveBoard(board *model.Board) (*model.Board, error) {
	result := Db.Create(board)
	if result.Error != nil {
		return nil, result.Error
	}

	return board, nil
}

// SaveUserBoard associates a user with a board in the users_boards join table.
func (r *BoardRepo) SaveUserBoard(ctx context.Context, userId int64, boardId int64) error {
	// Start a transaction
	tx := Db.WithContext(ctx)

	// Check if the board exists
	var board model.Board
	if err := tx.First(&board, boardId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("board not found")
		}
		return err
	}

	// Check if the user exists
	var user model.User
	if err := tx.First(&user, userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Associate the user with the board
	if err := tx.Model(&board).Association("Users").Append(&user); err != nil {
		return err
	}

	return nil
}
