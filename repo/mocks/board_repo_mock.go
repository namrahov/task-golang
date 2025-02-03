package mocks

import (
	"context"
	"task-golang/model"

	"github.com/stretchr/testify/mock"
)

type BoardRepoMock struct {
	mock.Mock
}

func (r *BoardRepoMock) SaveBoard(board *model.Board) (*model.Board, error) {
	args := r.Called(board)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Board), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *BoardRepoMock) GetBoardById(id int64) (*model.Board, error) {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Board), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *BoardRepoMock) SaveUserBoard(ctx context.Context, userId int64, boardId int64) error {
	args := r.Called(ctx, userId, boardId)
	return args.Error(0)
}

func (r *BoardRepoMock) GetUserBoards(userId int64) (*[]model.Board, error) {
	args := r.Called(userId)
	if args.Get(0) != nil {
		return args.Get(0).(*[]model.Board), args.Error(1)
	}
	return nil, args.Error(1)
}
