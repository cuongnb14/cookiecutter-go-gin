package repositories

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IRepository[T any, TK any] interface {
	Create(entity *T) error
	GetByPK(id TK) (*T, error)
	GetAll() ([]*T, error)
	Update(entity *T) error
	Delete(id TK) error
	GetOneBy(query interface{}, args ...interface{}) (*T, error)
	SelectForUpdateByPK(id TK, option string) (*T, error)
	BatchCreate(entities *[]T) error
	SetDB(db *gorm.DB)
}

type repository[T any, TK any] struct {
	db *gorm.DB
}

func (r *repository[T, TK]) SetDB(db *gorm.DB) {
	r.db = db
}

func (r *repository[T, TK]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *repository[T, TK]) GetByPK(id TK) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *repository[T, TK]) SelectForUpdateByPK(id TK, option string) (*T, error) {
	var entity T
	err := r.db.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate, Options: option}).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *repository[T, TK]) GetAll() ([]*T, error) {
	var entities []*T
	err := r.db.Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *repository[T, TK]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

func (r *repository[T, TK]) Delete(id TK) error {
	return r.db.Delete(new(T), id).Error
}

func (r *repository[T, TK]) GetOneBy(query interface{}, args ...interface{}) (*T, error) {
	var model T
	result := r.db.Where(query, args...).Take(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return &model, nil
}

func (r *repository[T, TK]) BatchCreate(entities *[]T) error {
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(entities).Error
}
