package usecase

import (
	"clean-architecture/domain"
	"clean-architecture/usecase/contract"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TaskUseCase struct {
	TaskRepository contract.ITaskRepository
	Logger         contract.ILogger
}

func NewTaskUseCase(repo contract.ITaskRepository, lg contract.ILogger) *TaskUseCase {
	return &TaskUseCase{
		TaskRepository: repo,
		Logger:         lg,
	}
}

func (u *TaskUseCase) CreateNewTask(ctx context.Context, title, description string, status domain.TaskStatus, dueDate time.Time) (*domain.Task, error) {
	// check if the required input feilds are valid
	if title == "" || description == "" || status == "" {
		u.Logger.Error("empty input fields")
		return &domain.Task{}, errors.New("invalid inputs")
	}

	if dueDate.IsZero() {
		u.Logger.Debug("CreateNewTask: due date is zero")
		return nil, errors.New("invalid duedate")
	}
	if dueDate.Before(time.Now().UTC()) {
		u.Logger.Info(fmt.Sprintf("CreateNewTask: invalid duedate. it is in the past. duedate: %v", dueDate.Format(time.RFC3339)))
		return nil, errors.New("invalid duedate")
	}

	// store it in domain.task{}
	newTask := domain.Task{
		TaskID:      uuid.NewString(),
		Title:       title,
		Description: description,
		Status:      status,
		DueDate:     dueDate,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	// send it to the db
	// check if an error occured while saving it to the db
	if err := u.TaskRepository.CreateTask(ctx, &newTask); err != nil {
		u.Logger.Error("failed to save new task in db")
		return &domain.Task{}, errors.New("internal server error")
	}

	// return the created task
	return &newTask, nil
}

func (u *TaskUseCase) UpdateTask(ctx context.Context, id, title, description string, status domain.TaskStatus, dueDate time.Time) (*domain.Task, error) {
	// check if task exist
	task, err := u.TaskRepository.GetTaskByID(ctx, id)
	if err != nil {
		u.Logger.Error(fmt.Sprintf("task not found: %s", err))
		return &domain.Task{}, errors.New("task not found")
	}
	// update task with the filled input
	if title != "" {
		task.Title = title
	}
	if description != "" {
		task.Description = description
	}
	status, err = domain.NewTaskStatusFromString(string(status))
	if err != nil {
		return &domain.Task{}, errors.New("invalid status")
	}
	task.Status = status

	if dueDate.Before(time.Now().UTC()) {
		u.Logger.Info(fmt.Sprintf("CreateNewTask: invalid duedate. it is in the past. duedate: %v", dueDate.Format(time.RFC3339)))
		return nil, errors.New("invalid duedate")
	} else {
		task.DueDate = dueDate
	}

	task.UpdatedAt = time.Now().UTC()
	// save updated task to db
	// check if it was saved successfully
	if err = u.TaskRepository.UpdateTask(ctx, task); err != nil {
		return &domain.Task{}, errors.New("internal server error")
	}
	return task, nil
}

func (u *TaskUseCase) DeleteTask(ctx context.Context, id string) error {
	err := u.TaskRepository.DeleteTask(ctx, id)
	if err != nil {
		u.Logger.Error(fmt.Sprintf("failed to delete task: %s", err))
		return errors.New("failed to delete task")
	}
	return nil
}

func (u *TaskUseCase) CompleteTask(ctx context.Context, id string) error {
	task, err := u.TaskRepository.GetTaskByID(ctx, id)
	if err != nil {
		u.Logger.Error(fmt.Sprintf("failed to fetch task: %s", err))
		return errors.New("failed to fetch task")
	}

	task.MarkCompleted()
	task.UpdatedAt = time.Now().UTC()

	if err = u.TaskRepository.UpdateTask(ctx, task); err != nil {
		return errors.New("internal server error")
	}
	return nil
}

func (u *TaskUseCase) OverDueTasks(ctx context.Context) ([]*domain.Task, error) {
	currentTime := time.Now().UTC()
	tasks, err := u.TaskRepository.GetAllOverDueTasks(ctx, currentTime)
	if err != nil {
		u.Logger.Error("failed to fetch overdue tasks")
		return []*domain.Task{}, errors.New("failed to fetch overdue tasks")
	}
	return tasks, nil
}
