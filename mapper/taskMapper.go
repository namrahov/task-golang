package mapper

import (
	"strconv"
	"task-golang/config"
	"task-golang/model"
	"time"
)

func BuildTask(dto *model.TaskRequestDto, user *model.User, board *model.Board) *model.Task {
	task := &model.Task{
		Name:         dto.Name,
		Priority:     dto.Priority,
		Status:       model.NotStarted,
		CreatedByID:  &user.Id,  // Setting the foreign key
		AssignedByID: &user.Id,  // e.g. assigned by same user
		BoardID:      &board.Id, // referencing the board
		CreatedAt:    time.Now(),
	}
	return task
}

// BuildTaskResponse maps a Task, its Attachment Files, and Task Videos to a TaskResponseDto.
func BuildTaskResponse(task *model.Task, attachments *[]model.TaskAttachmentFile, taskTaskVideos *[]model.TaskTaskVideo) *model.TaskResponseDto {
	// Build the slice of attachment file IDs.
	attachmentFileIds := make([]int64, 0, len(*attachments))
	for _, att := range *attachments {
		// Ensure AttachmentFileID is not nil before dereferencing.
		if att.AttachmentFileID != nil {
			attachmentFileIds = append(attachmentFileIds, *att.AttachmentFileID)
		}
	}

	// Build the slice of task video IDs.
	taskVideoIds := make([]int64, 0, len(*taskTaskVideos))
	for _, tv := range *taskTaskVideos {
		taskVideoIds = append(taskVideoIds, tv.TaskVideoID)
	}

	// Construct and return the TaskResponseDto.
	return &model.TaskResponseDto{
		Id:                task.Id,
		Name:              task.Name,
		Priority:          task.Priority,
		Status:            task.Status,
		CreatedBy:         task.CreatedBy,
		ChangedBy:         task.ChangedBy,
		AssignedBy:        task.AssignedBy,
		AssignedTo:        task.AssignedTo,
		Board:             task.Board,
		Deadline:          task.Deadline,
		AttachmentFileIds: attachmentFileIds,
		TaskVideoId:       taskVideoIds,
		TaskImageUrl:      config.Props.BaseUrl + "/v1/files/get/task-image/" + strconv.FormatInt(task.Id, 10),
	}
}
