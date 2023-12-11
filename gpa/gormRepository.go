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

type GormRepository[E Entity[ID], ID Id] interface {
	CrudRepository[E, ID]
}

type GormRepositoryImpl[E GormEntity[ID], ID GormEntityId] struct {
	Model E        // 제너릭을 통해 자동 설정됨.
	DB    *gorm.DB // 필수로 추가해야함
}

func (r *GormRepositoryImpl[E, ID]) Save(e E) error {
	return r.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&e).Error
}

func (r *GormRepositoryImpl[E, ID]) SaveAll(e ...E) (int64, error) {
	res := r.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&e)
	return res.RowsAffected, res.Error
}

func (r *GormRepositoryImpl[E, ID]) FindById(id ID) (e E, err error) {
	err = r.DB.Find(&e, id).Error
	return
}

func (r *GormRepositoryImpl[E, ID]) FindAllById(id ...ID) (e []E, err error) {
	err = r.DB.Find(&e, id).Error
	return
}

func (r *GormRepositoryImpl[E, ID]) Delete(e E) error {
	return r.DB.Delete(&e).Error
}

func (r *GormRepositoryImpl[E, ID]) DeleteAll(e ...E) (int64, error) {
	res := r.DB.Delete(&e)
	return res.RowsAffected, res.Error
}

func (r *GormRepositoryImpl[E, ID]) DeleteById(id ID) error {
	e, err := r.FindById(id)
	if err != nil {
		return err
	}
	return r.DB.Delete(&e).Error
}

func (r *GormRepositoryImpl[E, ID]) DeleteAllById(id ...ID) (int64, error) {
	e, err := r.FindAllById(id...)
	if err != nil {
		return 0, err
	}
	res := r.DB.Delete(&e)
	return res.RowsAffected, res.Error
}
