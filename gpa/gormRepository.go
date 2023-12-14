package gpa

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormEntityId interface {
	Id
}

type GormEntity[ID GormEntityId] interface {
	Entity[ID]
	TableName() string
}

type GormRepository[E GormEntity[ID], ID GormEntityId] interface {
	CrudRepository[E, ID]
}

type GormRepositoryImpl[E GormEntity[ID], ID GormEntityId] struct {
	db *gorm.DB // 필수로 추가해야함
}

func NewGormRepositoryImpl[E GormEntity[ID], ID GormEntityId](tx *gorm.DB) *GormRepositoryImpl[E, ID] {
	return &GormRepositoryImpl[E, ID]{db: tx}
}

func NewGormRepository[E GormEntity[ID], ID GormEntityId](tx *gorm.DB) GormRepository[E, ID] {
	return &GormRepositoryImpl[E, ID]{db: tx}
}

func (r *GormRepositoryImpl[E, ID]) Save(e E) error {
	return r.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&e).Error
}

func (r *GormRepositoryImpl[E, ID]) SaveAll(e ...E) (int64, error) {
	res := r.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&e)
	return res.RowsAffected, res.Error
}

func (r *GormRepositoryImpl[E, ID]) FindById(id ID) (e E, err error) {
	err = r.db.Find(&e, id).Error
	return
}

func (r *GormRepositoryImpl[E, ID]) FindAllById(id ...ID) (e []E, err error) {
	err = r.db.Find(&e, id).Error
	return
}

func (r *GormRepositoryImpl[E, ID]) Delete(e E) error {
	return r.db.Delete(&e).Error
}

func (r *GormRepositoryImpl[E, ID]) DeleteAll(e ...E) (int64, error) {
	res := r.db.Delete(&e)
	return res.RowsAffected, res.Error
}

func (r *GormRepositoryImpl[E, ID]) DeleteById(id ID) error {
	e, err := r.FindById(id)
	if err != nil {
		return err
	}
	return r.db.Delete(&e).Error
}

func (r *GormRepositoryImpl[E, ID]) DeleteAllById(id ...ID) (int64, error) {
	e, err := r.FindAllById(id...)
	if err != nil {
		return 0, err
	}
	res := r.db.Delete(&e)
	return res.RowsAffected, res.Error
}
