package test

import (
	"errors"
	"fmt"
	"net/http"
	"task-golang/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"task-golang/model"
	repoMock "task-golang/repo/mocks"
)

// Mock dependencies
var (
	mockUserRepo = new(repoMock.UserRepoMock)

	boardService = service.BoardService{
		BoardRepo: mockBoardRepo,
		UserRepo:  mockUserRepo,
		UserUtil:  mockUserUtil,
	}

	testBoardDto = &model.BoardRequestDto{
		Name: "Test Board",
	}
)

func TestBoardService_CreateBoard_Success(t *testing.T) {
	mockUserUtil.On("GetUserFromRequest", mock.Anything).Return(testUser, (*model.ErrorResponse)(nil)).Once()
	mockBoardRepo.On("SaveBoard", mock.Anything).Return(testBoard, nil).Once()

	err := boardService.CreateBoard(mockContext(), testBoardDto)

	assert.Nil(t, err)
	mockUserUtil.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardService_CreateBoard_UserNotFound(t *testing.T) {
	mockUserUtil.On("GetUserFromRequest", mock.Anything).Return(nil, &model.ErrorResponse{
		Error:   "user.not.found",
		Message: "User not found",
		Code:    http.StatusNotFound,
	}).Once()

	err := boardService.CreateBoard(mockContext(), testBoardDto)

	assert.NotNil(t, err)
	assert.Equal(t, "user.not.found", err.Error)
	mockUserUtil.AssertExpectations(t)
}

func TestBoardService_CreateBoard_SaveBoardError(t *testing.T) {
	mockUserUtil.On("GetUserFromRequest", mock.Anything).Return(testUser, (*model.ErrorResponse)(nil)).Once()
	mockBoardRepo.On("SaveBoard", mock.Anything).Return(nil, errors.New("database error")).Once()

	err := boardService.CreateBoard(mockContext(), testBoardDto)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("%s.cant-take-board", model.Exception), err.Error)
	mockUserUtil.AssertExpectations(t)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardService_GiveAccessToBoard_Success(t *testing.T) {
	mockBoardRepo.On("SaveUserBoard", mock.Anything, testUser.Id, testBoardId).Return(nil).Once()

	err := boardService.GiveAccessToBoard(mockContext(), testUser.Id, testBoardId)

	assert.Nil(t, err)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardService_GiveAccessToBoard_SaveError(t *testing.T) {
	mockBoardRepo.On("SaveUserBoard", mock.Anything, testUser.Id, testBoardId).Return(errors.New("save error")).Once()

	err := boardService.GiveAccessToBoard(mockContext(), testUser.Id, testBoardId)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("%s.cant-save-user-board", model.Exception), err.Error)
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardService_GetUserBoards_Success(t *testing.T) {
	mockBoardRepo.On("GetUserBoards", testUser.Id).Return(&[]model.Board{*testBoard}, nil).Once()

	boards, err := boardService.GetUserBoards(mockContext(), testUser.Id)

	assert.Nil(t, err)
	assert.NotNil(t, boards)
	assert.Equal(t, 1, len(*boards))
	mockBoardRepo.AssertExpectations(t)
}

func TestBoardService_GetUserBoards_Error(t *testing.T) {
	mockBoardRepo.On("GetUserBoards", testUser.Id).Return(nil, errors.New("fetch error")).Once()

	boards, err := boardService.GetUserBoards(mockContext(), testUser.Id)

	assert.Nil(t, boards)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("%s.cant-get-users-boards", model.Exception), err.Error)
	mockBoardRepo.AssertExpectations(t)
}
