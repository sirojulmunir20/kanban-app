package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	data := entity.User{}
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Take(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, nil
		}
		return entity.User{}, err
	}
	return data, nil // TODO: replace this
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	data := entity.User{}
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("email = ?", email).Find(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, nil
		}
		return entity.User{}, err
	}

	return data, nil // TODO: replace this
}

func (r *userRepository) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	result := r.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return entity.User{}, result.Error
	}
	return user, nil // TODO: replace this
}

func (r *userRepository) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	err := r.db.WithContext(ctx).Model(&entity.User{}).Updates(user).Error
	if err != nil {
		return entity.User{}, err
	}

	return user, nil // TODO: replace this
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.User{})

	if result.Error != nil {
		return result.Error
	}

	return nil // TODO: replace this
}
