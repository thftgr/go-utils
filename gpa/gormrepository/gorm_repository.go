package gormrepository

import (
	"database/sql"
	"github.com/thftgr/go-utils/gpa"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormEntityId interface {
	gpa.Id
}

type GormEntity[ID GormEntityId] interface {
	gpa.Entity[ID]
	TableName() string
}

type GormRepository[E GormEntity[ID], ID GormEntityId] interface {
	gpa.CrudRepository[E, ID]
	Begin(opts ...*sql.TxOptions) GormRepository[E, ID]
	Rollback() error
	Commit() error
}

type GormRepositoryImpl[E GormEntity[ID], ID GormEntityId] struct {
	DB *gorm.DB // 필수로 추가해야함
}

func NewGormRepositoryImpl[E GormEntity[ID], ID GormEntityId](tx *gorm.DB) *GormRepositoryImpl[E, ID] {
	return &GormRepositoryImpl[E, ID]{DB: tx}
}

func NewGormRepository[E GormEntity[ID], ID GormEntityId](tx *gorm.DB) GormRepository[E, ID] {
	return &GormRepositoryImpl[E, ID]{DB: tx}
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
func (r *GormRepositoryImpl[E, ID]) Begin(opts ...*sql.TxOptions) GormRepository[E, ID] {
	return NewGormRepositoryImpl[E, ID](r.DB.Begin(opts...))
}
func (r *GormRepositoryImpl[E, ID]) Rollback() error {
	return r.DB.Rollback().Error
}
func (r *GormRepositoryImpl[E, ID]) Commit() error {
	return r.DB.Commit().Error
}
