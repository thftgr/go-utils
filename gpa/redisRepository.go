package gpa

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisEntityId interface {
	Id
	ToString() string // redis key 가 string 임.
}

type RedisEntity[ID RedisEntityId] interface {
	Entity[ID]

	// 필요한 경우 아래 두 인터페이스를 구현할것.
	//json.Marshaler
	//json.Unmarshaler

	// entity struct 구현시 redis 태그 사용
	// UserId int `redis:"userId"`
	// Email  int `redis:"email"`
}

type RedisRepository[E RedisEntity[ID], ID RedisEntityId] interface {
	CrudRepository[E, ID]
}

type RedisRepositoryImpl[E RedisEntity[ID], ID RedisEntityId] struct {
	ctx   context.Context
	Model E               // 제너릭을 통해 자동 설정됨.
	pipe  redis.Pipeliner // 필수로 추가해야함
}

func (r *RedisRepositoryImpl[E, ID]) Save(e E) error {
	return r.pipe.HMSet(r.ctx, e.GetId().ToString(), e).Err()
}

func (r *RedisRepositoryImpl[E, ID]) SaveAll(e ...E) (count int64, err error) {
	for i := range e {
		err2 := r.pipe.HMSet(r.ctx, e[i].GetId().ToString(), e[i]).Err()
		if err2 != nil {
			err = err2
			return
		}
		count++
	}
	_, err = r.pipe.Exec(r.ctx)
	return
}

func (r *RedisRepositoryImpl[E, ID]) FindById(id ID) (e E, err error) {
	err = r.pipe.HMGet(r.ctx, id.ToString()).Scan(&e)
	return
}

func (r *RedisRepositoryImpl[E, ID]) FindAllById(id ...ID) (res []E, err error) {
	for i := range id {
		e, err2 := r.FindById(id[i])
		if err2 != nil {
			err = err2
			return
		}
		res = append(res, e)
	}
	return
}

func (r *RedisRepositoryImpl[E, ID]) Delete(e E) error {
	return r.pipe.Del(r.ctx, e.GetId().ToString()).Err()
}

func (r *RedisRepositoryImpl[E, ID]) DeleteAll(e ...E) (count int64, err error) {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].GetId().ToString()
	}
	return r.pipe.Del(r.ctx, ids...).Result()

}

func (r *RedisRepositoryImpl[E, ID]) DeleteById(id ID) error {
	return r.pipe.Del(r.ctx, id.ToString()).Err()
}

func (r *RedisRepositoryImpl[E, ID]) DeleteAllById(id ...ID) (count int64, err error) {
	ids := make([]string, len(id))
	for i := range id {
		ids[i] = id[i].ToString()
	}
	return r.pipe.Del(r.ctx, ids...).Result()
}
