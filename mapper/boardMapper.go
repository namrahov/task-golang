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

func BuildBoards(boards *[]model.Board) *[]model.BoardResponseDto {
	if boards == nil {
		return nil
	}

	// Create a slice to hold the response DTOs
	boardDtos := make([]model.BoardResponseDto, 0, len(*boards))

	// Iterate over the boards and map them to BoardResponseDto
	for _, board := range *boards {
		boardDto := model.BoardResponseDto{
			Id:   board.Id,
			Name: board.Name,
		}
		boardDtos = append(boardDtos, boardDto)
	}

	return &boardDtos
}
