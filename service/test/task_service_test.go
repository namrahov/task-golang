package test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	repoMock "task-golang/repo/mocks"
	"task-golang/service"
	utilMocks "task-golang/util/mocks"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"task-golang/model"
)

// mockContext creates a test context with a logger
func mockContext() context.Context {
	ctx := context.Background()
	logger := logrus.NewEntry(logrus.New())
	return context.WithValue(ctx, model.ContextLogger, logger)
}

var (
	mockBoardRepo = new(repoMock.BoardRepoMock)
	mockTaskRepo  = new(repoMock.TaskRepoMock)
	mockFileRepo  = new(repoMock.FileRepoMock)
	mockUserUtil  = new(utilMocks.UserUtilMock)

	taskService = service.TaskService{
		BoardRepo: mockBoardRepo,
		TaskRepo:  mockTaskRepo,
		FileRepo:  mockFileRepo,
		UserUtil:  mockUserUtil,
	}

	testBoardId  int64 = 1
	testUserId   int64 = 123
	testTaskId   int64 = 456
	testTaskName       = "Test Task"
	testPriority       = model.High

	testDto = &model.TaskRequestDto{
		Name:     testTaskName,
		Priority: testPriority,
		Deadline: "2025-12-01T00:00:00Z",
	}

	testBoard = &model.Board{
		Id:   testBoardId,
		Name: "Test Board",
	}

	testUser = &model.User{
		Id:       testUserId,
		UserName: "Test User",
	}

	testTask = &model.Task{
		Id:          testTaskId,
		Name:        testTaskName,
		Priority:    testPriority,
		Status:      model.NotStarted,
		BoardID:     &testBoardId,
		CreatedByID: &testUserId,
		CreatedAt:   time.Now(),
	}
)

func TestTaskService_CreateTask_Success(t *testing.T) {
	mockBoardRepo.On("GetBoardById", testBoardId).Return(testBoard, nil).Once()
	mockUserUtil.On("GetUserFromRequest", mock.Anything).Return(testUser, (*model.ErrorResponse)(nil)).Once()
	mockTaskRepo.On("SaveTask", mock.Anything).Return(testTask, nil).Once()

	err := taskService.CreateTask(mockContext(), testDto, testBoardId)

	assert.Nil(t, err)
	mockBoardRepo.AssertExpectations(t)
	mockUserUtil.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_CreateTask_BoardNotFound(t *testing.T) {
	mockBoardRepo.On("GetBoardById", testBoardId).Return(nil, errors.New("board not found")).Once()

	err := taskService.CreateTask(mockContext(), testDto, testBoardId)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("%s.cant-take-board", model.Exception), err.Error)
	mockBoardRepo.AssertExpectations(t)
}

func TestTaskService_CreateTask_UserFetchError(t *testing.T) {
	mockBoardRepo.On("GetBoardById", testBoardId).Return(testBoard, nil).Once()
	mockUserUtil.On("GetUserFromRequest", mock.Anything).Return(nil, &model.ErrorResponse{
		Error:   "user not found",
		Message: "User retrieval failed",
		Code:    http.StatusNotFound,
	}).Once()

	err := taskService.CreateTask(mockContext(), testDto, testBoardId)

	assert.NotNil(t, err)
	assert.Equal(t, "user not found", err.Error)
	mockBoardRepo.AssertExpectations(t)
	mockUserUtil.AssertExpectations(t)
}

func TestTaskService_CreateTask_SaveTaskError(t *testing.T) {
	mockBoardRepo.On("GetBoardById", testBoardId).Return(testBoard, nil).Once()
	mockUserUtil.On("GetUserFromRequest", mock.Anything).Return(testUser, (*model.ErrorResponse)(nil)).Once()
	mockTaskRepo.On("SaveTask", mock.Anything).Return(nil, errors.New("database error")).Once()

	err := taskService.CreateTask(mockContext(), testDto, testBoardId)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("%s.cant-save-task", model.Exception), err.Error)
	mockBoardRepo.AssertExpectations(t)
	mockUserUtil.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_GetTask_Success(t *testing.T) {
	mockTaskRepo.On("GetTaskById", testTaskId).Return(testTask, nil).Once()
	mockFileRepo.On("FindTaskAttachmentsFileByTaskId", testTaskId).Return(&[]model.TaskAttachmentFile{}, nil).Once()
	mockFileRepo.On("FindTaskTaskVideosByTaskId", testTaskId).Return(&[]model.TaskTaskVideo{}, nil).Once()

	taskResponse, err := taskService.GetTask(mockContext(), testTaskId)

	assert.Nil(t, err)
	assert.NotNil(t, taskResponse)
	assert.Equal(t, testTaskId, taskResponse.Id)
	mockTaskRepo.AssertExpectations(t)
	mockFileRepo.AssertExpectations(t)
}

func TestTaskService_GetTask_NotFound(t *testing.T) {
	mockTaskRepo.On("GetTaskById", testTaskId).Return(nil, errors.New("task not found")).Once()

	taskResponse, err := taskService.GetTask(mockContext(), testTaskId)

	assert.Nil(t, taskResponse)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("%s.cant-get-task", model.Exception), err.Error)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_GetTasks_Success(t *testing.T) {
	mockTaskRepo.On("GetTasks", testTaskName, string(testPriority), testBoardId, 1, 10).
		Return(&model.TaskPageResponseDto{
			Tasks:          []*model.TaskResponseDto{{Id: testTaskId, Name: testTaskName}},
			HasNextPage:    false,
			LastPageNumber: 1,
			TotalCount:     1,
		}, nil).Once()

	response, err := taskService.GetTasks(mockContext(), testTaskName, string(testPriority), testBoardId, 1, 10)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int64(1), response.TotalCount)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_GetTasks_Error(t *testing.T) {
	mockTaskRepo.On("GetTasks", testTaskName, string(testPriority), testBoardId, 1, 10).
		Return(nil, errors.New("database error")).Once()

	response, err := taskService.GetTasks(mockContext(), testTaskName, string(testPriority), testBoardId, 1, 10)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("%s.cant-get-tasks", model.Exception), err.Error)
	mockTaskRepo.AssertExpectations(t)
}
