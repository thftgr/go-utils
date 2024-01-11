package redisrepository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func NewRedisClient() (DB *redis.Client) {
	// 환경변수로 redis 접속 url 을 제공할것.
	// REDIS_URL=redis://user:password@localhost:6379/0?protocol=3
	opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}
	DB = redis.NewClient(opts)
	tctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	status := DB.Ping(tctx)
	if status.Err() != nil {
		fmt.Printf("Failed to connect redis. status: %s err:%+v", status.String(), status.Err())
	} else {
		fmt.Println("REDIS Connected")
	}
	return
}

type RedisTestEntity struct {
	Id      int    `redis:"id,omitempty"`
	Name    string `redis:"name,omitempty"`
	LoginId string `redis:"loginId,omitempty"`
}

func (r RedisTestEntity) GetId() string {
	return "user_test::" + strconv.Itoa(r.Id)
}

//=====================================================================================================================
//=====================================================================================================================
//=====================================================================================================================

func TestRedisRepositoryImpl_Save(t *testing.T) {
	client := NewRedisClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	type testCase[E RedisEntity[ID], ID RedisEntityId] struct {
		name string
		repo *RedisRepositoryImpl[E, ID]
		args E
	}
	tests := []testCase[RedisTestEntity, string]{
		{"", NewRedisRepositoryImpl[RedisTestEntity, string](ctx, client), RedisTestEntity{1, "1", "1"}},
		{"", NewRedisRepositoryImpl[RedisTestEntity, string](ctx, client), RedisTestEntity{2, "2", "2"}},
		{"", NewRedisRepositoryImpl[RedisTestEntity, string](ctx, client), RedisTestEntity{2, "3", "3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.repo.Save(tt.args); err != nil {
				t.Errorf("Save() error = %v", err)
			}
			if e, err := tt.repo.FindById(tt.args.GetId()); err != nil {
				t.Errorf("FindById() error = %v", err)
			} else if !reflect.DeepEqual(e, tt.args) {
				t.Errorf("FindById() want = %v, get %v", tt.args, e)
			} else {
				t.Logf("FindById() want = %v, get %v", tt.args, e)
			}
		})
	}
}

func TestRedisRepositoryImpl_SaveAll(t *testing.T) {
	client := NewRedisClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	type testCase[E RedisEntity[ID], ID RedisEntityId] struct {
		name string
		repo *RedisRepositoryImpl[E, ID]
		args []E
		want []E
	}
	tests := []testCase[RedisTestEntity, string]{
		{
			name: "",
			repo: NewRedisRepositoryImpl[RedisTestEntity, string](ctx, client),
			args: []RedisTestEntity{{1, "1", "1"}, {2, "2", "2"}, {2, "3", "3"}},
			want: []RedisTestEntity{{1, "1", "1"}, {2, "3", "3"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if count, err := tt.repo.SaveAll(tt.args...); err != nil {
				t.Errorf("SaveAll() error = %v", err)
			} else if int(count) != len(tt.args) {
				t.Errorf("save count not match. args:%d. count:%d", len(tt.args), count)
			}
			var ids []string
			for i := range tt.want {
				ids = append(ids, tt.want[i].GetId())
			}
			if e, err := tt.repo.FindAllById(ids...); err != nil {
				t.Errorf("FindAllById() error = %v", err)
			} else if !reflect.DeepEqual(e, tt.want) {
				t.Errorf("FindAllById() want = %v, get %v", tt.want, e)
			} else {
				t.Logf("FindAllById() want = %v, get %v", tt.want, e)
			}
		})
	}
}

func TestRedisRepositoryImpl_Delete(t *testing.T) {
	client := NewRedisClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	type testCase[E RedisEntity[ID], ID RedisEntityId] struct {
		name   string
		repo   *RedisRepositoryImpl[E, ID]
		data   []E
		delete E
		want   []E
	}
	tests := []testCase[RedisTestEntity, string]{
		{
			name:   "",
			repo:   NewRedisRepositoryImpl[RedisTestEntity, string](ctx, client),
			data:   []RedisTestEntity{{1, "1", "1"}, {2, "2", "2"}},
			delete: RedisTestEntity{1, "1", "1"},
			want:   []RedisTestEntity{{2, "2", "2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if count, err := tt.repo.SaveAll(tt.data...); err != nil {
				t.Errorf("SaveAll() error = %v", err)
			} else if int(count) != len(tt.data) {
				t.Errorf("save count not match. args:%d. count:%d", len(tt.data), count)
			}

			if err := tt.repo.Delete(tt.delete); err != nil {
				t.Errorf("Delete() error = %v", err)
			}

			var ids []string
			for i := range tt.want {
				ids = append(ids, tt.want[i].GetId())
			}
			if e, err := tt.repo.FindAllById(ids...); err != nil {
				t.Errorf("SaveAll() error = %v", err)
			} else if !reflect.DeepEqual(e, tt.want) {
				t.Errorf("SaveAll() want = %v, get %v", tt.want, e)
			} else {
				t.Logf("SaveAll() want = %v, get %v", tt.want, e)
			}
		})
	}
}

func TestRedisRepositoryImpl_DeleteAll(t *testing.T) {
	client := NewRedisClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	type testCase[E RedisEntity[ID], ID RedisEntityId] struct {
		name   string
		repo   *RedisRepositoryImpl[E, ID]
		data   []E
		delete []E
		want   []E
	}
	tests := []testCase[RedisTestEntity, string]{
		{
			name:   "",
			repo:   NewRedisRepositoryImpl[RedisTestEntity, string](ctx, client),
			data:   []RedisTestEntity{{1, "1", "1"}, {2, "2", "2"}, {3, "3", "3"}},
			delete: []RedisTestEntity{{1, "1", "1"}, {2, "2", "2"}},
			want:   []RedisTestEntity{{3, "3", "3"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if count, err := tt.repo.SaveAll(tt.data...); err != nil {
				t.Errorf("SaveAll() error = %v", err)
			} else if int(count) != len(tt.data) {
				t.Errorf("save count not match. args:%d. count:%d", len(tt.data), count)
			}

			if count, err := tt.repo.DeleteAll(tt.delete...); err != nil {
				t.Errorf("Delete() error = %v", err)
			} else if int(count) != len(tt.delete) {
				t.Errorf("deleteAll count not match. args:%d. count:%d", len(tt.delete), count)
			}

			var ids []string
			for i := range tt.want {
				ids = append(ids, tt.want[i].GetId())
			}
			if e, err := tt.repo.FindAllById(ids...); err != nil {
				t.Errorf("SaveAll() error = %v", err)
			} else if !reflect.DeepEqual(e, tt.want) {
				t.Errorf("SaveAll() want = %v, get %v", tt.want, e)
			} else {
				t.Logf("SaveAll() want = %v, get %v", tt.want, e)
			}
		})
	}
}

func TestRedisRepositoryImpl_DeleteById(t *testing.T) {
	client := NewRedisClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	type testCase[E RedisEntity[ID], ID RedisEntityId] struct {
		name   string
		repo   *RedisRepositoryImpl[E, ID]
		data   []E
		delete E
		want   []E
	}
	tests := []testCase[RedisTestEntity, string]{
		{
			name:   "",
			repo:   NewRedisRepositoryImpl[RedisTestEntity, string](ctx, client),
			data:   []RedisTestEntity{{1, "1", "1"}, {2, "2", "2"}},
			delete: RedisTestEntity{1, "1", "1"},
			want:   []RedisTestEntity{{2, "2", "2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if count, err := tt.repo.SaveAll(tt.data...); err != nil {
				t.Errorf("SaveAll() error = %v", err)
			} else if int(count) != len(tt.data) {
				t.Errorf("save count not match. args:%d. count:%d", len(tt.data), count)
			}

			if err := tt.repo.DeleteById(tt.delete.GetId()); err != nil {
				t.Errorf("Delete() error = %v", err)
			}

			var ids []string
			for i := range tt.want {
				ids = append(ids, tt.want[i].GetId())
			}
			if e, err := tt.repo.FindAllById(ids...); err != nil {
				t.Errorf("SaveAll() error = %v", err)
			} else if !reflect.DeepEqual(e, tt.want) {
				t.Errorf("SaveAll() want = %v, get %v", tt.want, e)
			} else {
				t.Logf("SaveAll() want = %v, get %v", tt.want, e)
			}
		})
	}
}

func TestRedisRepositoryImpl_DeleteAllById(t *testing.T) {
	client := NewRedisClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	type testCase[E RedisEntity[ID], ID RedisEntityId] struct {
		name   string
		repo   *RedisRepositoryImpl[E, ID]
		data   []E
		delete []E
		want   []E
	}
	tests := []testCase[RedisTestEntity, string]{
		{
			name:   "",
			repo:   NewRedisRepositoryImpl[RedisTestEntity, string](ctx, client),
			data:   []RedisTestEntity{{1, "1", "1"}, {2, "2", "2"}, {3, "3", "3"}},
			delete: []RedisTestEntity{{1, "1", "1"}, {2, "2", "2"}},
			want:   []RedisTestEntity{{3, "3", "3"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if count, err := tt.repo.SaveAll(tt.data...); err != nil {
				t.Errorf("SaveAll() error = %v", err)
			} else if int(count) != len(tt.data) {
				t.Errorf("save count not match. args:%d. count:%d", len(tt.data), count)
			}
			var dids []string
			for i := range tt.delete {
				dids = append(dids, tt.delete[i].GetId())
			}
			if count, err := tt.repo.DeleteAllById(dids...); err != nil {
				t.Errorf("Delete() error = %v", err)
			} else if int(count) != len(dids) {
				t.Errorf("deleteAll count not match. args:%d. count:%d", len(tt.delete), count)
			}

			var ids []string
			for i := range tt.want {
				ids = append(ids, tt.want[i].GetId())
			}
			if e, err := tt.repo.FindAllById(ids...); err != nil {
				t.Errorf("SaveAll() error = %v", err)
			} else if !reflect.DeepEqual(e, tt.want) {
				t.Errorf("SaveAll() want = %v, get %v", tt.want, e)
			} else {
				t.Logf("SaveAll() want = %v, get %v", tt.want, e)
			}
		})
	}
}
