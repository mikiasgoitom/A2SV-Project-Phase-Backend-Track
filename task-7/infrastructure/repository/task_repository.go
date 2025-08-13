package repository

import (
	"clean-architecture/domain"
	"clean-architecture/usecase/contract"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ---------- Task DTO -------------------

type TaskModel struct {
	ID          string    `bson:"_id, omitempty"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	Status      string    `bson:"status"`
	DueDate     time.Time `bson:"due_date"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

func (model *TaskModel) ToDomain() (*domain.Task, error) {
	status, err := domain.NewTaskStatusFromString(model.Status)
	if err != nil {
		return nil, fmt.Errorf("invalid status from db for task %s: %w", model.ID, err)
	}
	return &domain.Task{
		TaskID:      model.ID,
		Title:       model.Title,
		Description: model.Description,
		Status:      status,
		DueDate:     model.DueDate,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}, nil
}

func DomainToModel(t *domain.Task) *TaskModel {
	return &TaskModel{
		ID:          t.TaskID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// ---------- end of task DTO ------=-----
type TaskRepository struct {
	Collection *mongo.Collection
	Logger     contract.ILogger
}

// implementing the interface at compile time
var _ contract.ITaskRepository = (*TaskRepository)(nil)

func NewTaskRepository(colln *mongo.Collection, lg contract.ILogger) *TaskRepository {
	if colln == nil {
		panic("em task collection")
	}
	return &TaskRepository{
		Collection: colln,
		Logger:     lg,
	}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task *domain.Task) error {
	taskMdl := DomainToModel(task)
	_, err := r.Collection.InsertOne(ctx, taskMdl)
	if err != nil {
		r.Logger.Error("failed to task store into db")
		return err
	}
	r.Logger.Info("successfully task saved to db")
	return nil
}

func (r *TaskRepository) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	filter := bson.M{"task_id": id}
	var taskMdl *TaskModel
	err := r.Collection.FindOne(ctx, filter).Decode(&taskMdl)
	if taskMdl != nil {
		r.Logger.Error(fmt.Sprintf("failed to fetch task from db with id: %v", id))
		return nil, err
	}
	task, err := taskMdl.ToDomain()
	if err != nil {
		r.Logger.Error(fmt.Sprintf("failed parsing dto to entity: %v", id))
		return nil, err
	}
	return task, nil
}
func (r *TaskRepository) GetAllOverDueTasks(ctx context.Context, currentTime time.Time) ([]*domain.Task, error) {
	cursor, err := r.Collection.Find(ctx, bson.D{})
	if err != nil {
		r.Logger.Error(fmt.Sprintf("failed to get list of tasks: %v", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var overDuetasks []*domain.Task

	for cursor.Next(ctx) {
		var taskMdl TaskModel
		if err = cursor.Decode(&taskMdl); err != nil {
			r.Logger.Error(fmt.Sprintf("failed to get the next task: %s", err))
			continue
		}
		task, err := taskMdl.ToDomain()
		if err != nil {
			r.Logger.Error(fmt.Sprintf("failed parsing dto to entity: %s", err))
			continue
		}

		if task.IsOverDue() {
			overDuetasks = append(overDuetasks, task)
		}
	}
	if err = cursor.Err(); err != nil {
		r.Logger.Error(fmt.Sprintf("curor error: %v", err))
		return nil, err
	}

	r.Logger.Info("successfully fetched all overdue tasks")
	return overDuetasks, nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, task *domain.Task) error {
	var taskMdl = DomainToModel(task)
	filter := bson.M{"task_id": task.TaskID}
	update := bson.D{
		{Key: "$set", Value: taskMdl},
	}
	updateCount, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.Logger.Debug(fmt.Sprintf("failed to store task update: %v", err))
		return err
	}
	if updateCount.MatchedCount == 0 {
		r.Logger.Debug(fmt.Sprintf("failed to find task with task id: %v", task.TaskID))
		return errors.New("task not found")
	}

	r.Logger.Info("task data successfully updated")
	return nil
}

func (r *TaskRepository) DeleteTask(ctx context.Context, id string) error {
	filter := bson.M{"task_id": id}
	count, err := r.Collection.DeleteOne(ctx, filter)
	if err != nil {
		r.Logger.Error(fmt.Sprintf("failed to delete task with ID: %s", id))
		return err
	}

	if count.DeletedCount == 0 {
		r.Logger.Debug(fmt.Sprintf("failed: task not found with ID: %s", id))
		return nil
	}
	r.Logger.Debug(fmt.Sprintf("(%v) user with task (%s) found", count, id))
	return nil
}
