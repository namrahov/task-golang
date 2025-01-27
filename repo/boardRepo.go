package repo

import (
	"task-golang/model"
)

type IBoardRepo interface {
	SaveBoard(board *model.Board) (*model.Board, error)
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
