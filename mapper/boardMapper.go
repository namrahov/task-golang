package mapper

import (
	"task-golang/model"
)

func BuildBoard(name string, createdBy string) *model.Board {
	board := &model.Board{
		Name:      name,
		CreatedBy: createdBy,
	}

	return board
}
