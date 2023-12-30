package meilisearchRepository

import (
	"errors"
	"fmt"
	"github.com/meilisearch/meilisearch-go"
	"github.com/mitchellh/mapstructure"
	"github.com/thftgr/go-utils/gpa"
	"time"
)

type MeiliSearchEntityId interface {
	gpa.Id
}

// MeiliSearchEntity
//
//	entity 는 반드시 `json:"id"` 태그를 포함해야합니다.
//	이외에도 필요한경우 json 테그를 포함해야합니다.
type MeiliSearchEntity[ID MeiliSearchEntityId] interface {
	gpa.Entity[ID]
}

type MeiliSearchRepository[E MeiliSearchEntity[ID], ID MeiliSearchEntityId] interface {
	gpa.CrudRepository[E, ID]
	Find(text string) ([]E, error)
	FindWithLimit(text string, limit int64) ([]E, error)
	FindWithPage(text string, page, size int) ([]E, error)
}

// MeiliSearchRepositoryImpl sync 방식으로 동작하도록 구현함
type MeiliSearchRepositoryImpl[E MeiliSearchEntity[ID], ID MeiliSearchEntityId] struct {
	Index *meilisearch.Index
	Limit int64
}

func NewMeiliSearchRepositoryImpl[E MeiliSearchEntity[ID], ID MeiliSearchEntityId](index *meilisearch.Index, limit int64) *MeiliSearchRepositoryImpl[E, ID] {
	return &MeiliSearchRepositoryImpl[E, ID]{Index: index, Limit: limit}
}

func (r *MeiliSearchRepositoryImpl[E, ID]) Save(e E) (err error) {
	info, err := r.Index.AddDocuments(e)
	if err != nil {
		return err
	}
	return r.AwaitFinish(info.TaskUID)
}

func (r *MeiliSearchRepositoryImpl[E, ID]) SaveAll(e ...E) (count int64, err error) {
	info, err := r.Index.AddDocuments(e)
	if err != nil {
		return
	}
	err = r.AwaitFinish(info.TaskUID)
	if err == nil {
		count = int64(len(e))
	}
	return
}

func (r *MeiliSearchRepositoryImpl[E, ID]) FindById(id ID) (entity E, err error) {
	err = r.Index.GetDocument(fmt.Sprint(id), nil, &entity)
	return
}

func (r *MeiliSearchRepositoryImpl[E, ID]) FindAllById(id ...ID) (entities []E, err error) {
	for i := range id {
		entity, e := r.FindById(id[i])
		if e != nil {
			return nil, e
		}
		entities = append(entities, entity)
	}
	return
}

func (r *MeiliSearchRepositoryImpl[E, ID]) Delete(e E) error {
	info, err := r.Index.DeleteDocument(fmt.Sprint(e.GetId()))
	if err != nil {
		return err
	}
	return r.AwaitFinish(info.TaskUID)
}

func (r *MeiliSearchRepositoryImpl[E, ID]) DeleteAll(e ...E) (count int64, err error) {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = fmt.Sprint(e[i].GetId())
	}
	info, err := r.Index.DeleteDocuments(ids)
	if err != nil {
		return 0, err
	}
	return int64(len(e)), r.AwaitFinish(info.TaskUID)
}

func (r *MeiliSearchRepositoryImpl[E, ID]) DeleteById(id ID) error {
	info, err := r.Index.DeleteDocument(fmt.Sprint(id))
	if err != nil {
		return err
	}
	return r.AwaitFinish(info.TaskUID)
}

func (r *MeiliSearchRepositoryImpl[E, ID]) DeleteAllById(id ...ID) (int64, error) {
	ids := make([]string, len(id))
	for i := range id {
		ids[i] = fmt.Sprint(id)
	}
	info, err := r.Index.DeleteDocuments(ids)
	if err != nil {
		return 0, err
	}
	return int64(len(id)), r.AwaitFinish(info.TaskUID)
}
func (r *MeiliSearchRepositoryImpl[E, ID]) Find(text string) (res []E, err error) {
	searchResponse, err := r.Index.Search(text, &meilisearch.SearchRequest{
		Limit: r.Limit, // 기본 제한값
	})
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(searchResponse.Hits, &res)
	return
}

func (r *MeiliSearchRepositoryImpl[E, ID]) FindWithLimit(text string, limit int64) (res []E, err error) {
	searchResponse, err := r.Index.Search(text, &meilisearch.SearchRequest{
		Limit: limit, // 기본 제한값
	})
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(searchResponse.Hits, &res)
	return
}
func (r *MeiliSearchRepositoryImpl[E, ID]) FindWithPage(text string, page, size int) (res []E, err error) {
	searchResponse, err := r.Index.Search(text, &meilisearch.SearchRequest{
		Offset: int64((page - 1) * size),
		Limit:  int64(size),
	})
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(searchResponse.Hits, &res)
	return
}

func (r *MeiliSearchRepositoryImpl[E, ID]) AwaitFinish(taskUID int64) (err error) {
	for {
		tr, errr := r.Index.GetTask(taskUID)
		if errr != nil {
			return errr
		}
		switch tr.Status {
		case meilisearch.TaskStatusEnqueued, meilisearch.TaskStatusProcessing:
			time.Sleep(time.Second)
		case meilisearch.TaskStatusSucceeded:
			return nil
		case meilisearch.TaskStatusFailed, meilisearch.TaskStatusCanceled:
			return errors.New(string(tr.Status))
		}
	}
}
