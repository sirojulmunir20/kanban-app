package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"
	"errors"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasks(ctx context.Context, id int) ([]entity.Task, error)
	StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error)
	GetTaskByID(ctx context.Context, id int) (entity.Task, error)
	GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) GetTasks(ctx context.Context, id int) ([]entity.Task, error) {
	data := []entity.Task{}
	err := r.db.WithContext(ctx).Model(&entity.Task{}).Where("user_id = ?", id).Find(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.Task{}, nil
		}
		return nil, err
	}
	return data, nil // TODO: replace this
}

func (r *taskRepository) StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error) {
	result := r.db.WithContext(ctx).Model(&entity.Task{}).Create(&task)
	if result.Error != nil {
		return 0, result.Error
	}
	taskId = task.ID 
	return taskId, nil // TODO: replace this
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	res:= entity.Task{}
	err := r.db.WithContext(ctx).Model(&entity.Task{}).Where("id = ?", id).Take(&res).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Task{}, nil
		}
		return entity.Task{}, err
	}

	return res, nil // TODO: replace this
}

func (r *taskRepository) GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error) {
	res:= []entity.Task{}
	err := r.db.WithContext(ctx).Model(&entity.Task{}).Where("category_id = ?", catId).Find(&res).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.Task{}, nil
		}
		return nil, err
	}
	return res, nil // TODO: replace this
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *entity.Task) error {
	err := r.db.WithContext(ctx).Model(&entity.Task{}).Where("id = ?",task.ID).Updates(&task)

	if err.Error != nil {
		return err.Error
	}
	return nil // TODO: replace this
}

func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {
	err := r.db.WithContext(ctx).Model(&entity.Task{}).Where("id = ?", id).Delete(&entity.Task{})

	if err.Error != nil {
		return err.Error
	}
	return nil // TODO: replace this
}
