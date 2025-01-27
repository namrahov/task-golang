package repo

import (
	"gorm.io/gorm"
	"task-golang/model"
)

type IBoardRepo interface {
	SaveBoard(tx *gorm.DB, board *model.Board) (*model.Board, error)
}

type BoardRepo struct {
}

func (r BoardRepo) SaveBoard(tx *gorm.DB, board *model.Board) (*model.Board, error) {
	err := tx.Create(board).Error
	if err != nil {
		return nil, err
	}
	return board, nil
}
