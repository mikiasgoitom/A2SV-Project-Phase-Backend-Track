package controller

import (
	"clean-architecture/delivary/handlerdto"
	"clean-architecture/domain"
	"clean-architecture/usecase/contract"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	TaskUseCase contract.ITaskUseCase
}

func NewTaskHandler(taskUC contract.ITaskUseCase) *TaskHandler {
	return &TaskHandler{
		TaskUseCase: taskUC,
	}
}

func (h *TaskHandler) CreateNewTask(ctx *gin.Context) {
	var taskDto handlerdto.TaskDto
	if err := ctx.ShouldBindJSON(&taskDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !domain.IsTaskStatusValid(taskDto.Status) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid task status"})
		return
	}
	task, err := h.TaskUseCase.CreateNewTask(ctx.Request.Context(), taskDto.Title, taskDto.Description, domain.TaskStatus(taskDto.Status), taskDto.DueDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"task": task})
}
func (h *TaskHandler) UpdateTask(ctx *gin.Context) {
	var taskDto handlerdto.TaskDto
	if err := ctx.ShouldBindJSON(&taskDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if taskDto.TaskID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "task ID is required"})
		return
	}
	if !domain.IsTaskStatusValid(taskDto.Status) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid task status"})
		return
	}
	updatedTask, err := h.TaskUseCase.UpdateTask(ctx.Request.Context(), taskDto.TaskID, taskDto.Title, taskDto.Description, domain.TaskStatus(taskDto.Status), taskDto.DueDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"task": updatedTask})
	ctx.JSON(http.StatusNoContent, nil)
}

func (h *TaskHandler) CompleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := h.TaskUseCase.CompleteTask(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func (h *TaskHandler) OverDueTasks(ctx *gin.Context) {
	tasks, err := h.TaskUseCase.OverDueTasks(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (h *TaskHandler) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := h.TaskUseCase.DeleteTask(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
