package redisrepository

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/thftgr/go-utils/gpa"
	"time"
)

// RedisEntityId struct 를 id 로 사용하는경우 fmt.Stringer 인터페이스를 구현.
type RedisEntityId interface {
	gpa.Id
}

// RedisEntity struct 를 ID 로 사용하는경우 fmt.Stringer 인터페이스를 구현.
type RedisEntity[ID RedisEntityId] interface {
	gpa.Entity[ID]

	// 필요한 경우 아래 두 인터페이스를 구현.
	//json.Marshaler
	//json.Unmarshaler

	// entity struct 구현시 redis 태그 사용
	// UserId int `redis:"userId"`
	// Email  int `redis:"email"`
}

type RedisRepository[E RedisEntity[ID], ID RedisEntityId] interface {
	gpa.CrudRepository[E, ID]
	Flush()
	FlushWithContext(ctx context.Context)
}

// RedisRepositoryImpl
type RedisRepositoryImpl[E RedisEntity[ID], ID RedisEntityId] struct {
	Context context.Context
	DB      *redis.Client // 필수로 추가해야함
	Timeout *time.Duration
}

func NewRedisRepositoryImpl[E RedisEntity[ID], ID RedisEntityId](ctx context.Context, pipe *redis.Client) *RedisRepositoryImpl[E, ID] {
	return &RedisRepositoryImpl[E, ID]{Context: ctx, DB: pipe}
}
func NewRedisRepositoryImplWithTimeout[E RedisEntity[ID], ID RedisEntityId](ctx context.Context, pipe *redis.Client, timeout time.Duration) *RedisRepositoryImpl[E, ID] {
	return &RedisRepositoryImpl[E, ID]{Context: ctx, DB: pipe, Timeout: &timeout}
}
func (r *RedisRepositoryImpl[E, ID]) GetTimeoutContext() (context.Context, context.CancelFunc) {
	if r.Timeout == nil {
		return r.Context, func() {}
	}
	return context.WithTimeout(r.Context, *r.Timeout)
}

func (r *RedisRepositoryImpl[E, ID]) Save(e E) error {
	return r.DB.HSet(r.Context, fmt.Sprint(e.GetId()), e).Err()
}

func (r *RedisRepositoryImpl[E, ID]) SaveAll(e ...E) (count int64, err error) {
	ctx, cancel := r.GetTimeoutContext()
	defer cancel()
	pipe := r.DB.Pipeline()
	for i := range e {
		err2 := pipe.HSet(ctx, fmt.Sprint(e[i].GetId()), e[i]).Err() // 여기선 에러가 잘 발생하지 않음.
		if err2 != nil {
			err = err2
			return
		}
	}
	ce, err := pipe.Exec(ctx)
	for i := range ce {
		if ceerr := ce[i].Err(); ceerr != nil {
			err = errors.Join(err, ce[i].Err())
		} else {
			count++
		}
	}
	return
}

func (r *RedisRepositoryImpl[E, ID]) FindById(id ID) (e E, err error) {
	err = r.DB.HGetAll(r.Context, fmt.Sprint(id)).Scan(&e)
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
	return r.DB.Del(r.Context, fmt.Sprint(e.GetId())).Err()
}

func (r *RedisRepositoryImpl[E, ID]) DeleteAll(e ...E) (count int64, err error) {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = fmt.Sprint(e[i].GetId())
	}
	return r.DB.Del(r.Context, ids...).Result()

}

func (r *RedisRepositoryImpl[E, ID]) DeleteById(id ID) error {
	return r.DB.Del(r.Context, fmt.Sprint(id)).Err()
}

func (r *RedisRepositoryImpl[E, ID]) DeleteAllById(id ...ID) (count int64, err error) {
	ids := make([]string, len(id))
	for i := range id {
		ids[i] = fmt.Sprint(id[i])
	}
	return r.DB.Del(r.Context, ids...).Result()
}
