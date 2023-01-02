package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"
	"errors"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error)
	StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error)
	StoreManyCategory(ctx context.Context, categories []entity.Category) error
	GetCategoryByID(ctx context.Context, id int) (entity.Category, error)
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error) {
	res:= []entity.Category{}
	err:= r.db.WithContext(ctx).Model(&entity.Category{}).Where("user_id = ?",id).Find(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.Category{}, nil
		}
		return nil,err
	}

	return res, nil // TODO: replace this
}

func (r *categoryRepository) StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error) {
	result := r.db.WithContext(ctx).Model(&entity.Category{}).Create(&category)
	if result.Error != nil{
		return 0,result.Error
	}
	categoryId = category.ID
	return categoryId, nil // TODO: replace this
}

func (r *categoryRepository) StoreManyCategory(ctx context.Context, categories []entity.Category) error {
	result:= r.db.WithContext(ctx).Model(&entity.Category{}).Create(&categories)
	if result.Error != nil{
		return result.Error
	}
	return nil// TODO: replace this
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (entity.Category, error) {
	res:= entity.Category{}
	err:= r.db.WithContext(ctx).Model(&entity.Category{}).Where("id = ?",id).Find(&res).Error
	if err != nil{
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Category{}, nil
		}
		return entity.Category{},err
	}
	return res, nil // TODO: replace this
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *entity.Category) error {
	// data := entity.Category{}
	err := r.db.WithContext(ctx).Model(&entity.Category{}).Where("id = ?",category.ID).Updates(category).Error
	if err != nil {
		return err
	}
	return nil // TODO: replace this
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	res:= r.db.WithContext(ctx).Where("id = ?",id).Delete(&entity.Category{})
	if  res.Error != nil{
		return res.Error
	}
	return nil // TODO: replace this
}
