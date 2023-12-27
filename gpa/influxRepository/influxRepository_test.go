package influxRepository

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"math"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

type RepositoryTestEntity struct {
	Measurement `json:"-"  influxdb:"measurement:test"`
	Name        string     `influxdb:"tag:name"`
	LogString   *string    `influxdb:"field:log"`
	Status      *int64     `influxdb:"field:status"`
	Time        *time.Time `influxdb:"time"`
}

func (r RepositoryTestEntity) GetTime() time.Time {
	return *r.Time
}

type RepositoryTestEntityRepository struct {
	InfluxRepository[RepositoryTestEntity]
}

func TestInfluxRepositoryImpl_ToPoint(t *testing.T) {
	INFLUXDB_TOKEN := os.Getenv("INFLUXDB_TOKEN")
	INFLUXDB_BUCKET := os.Getenv("INFLUXDB_BUCKET")
	INFLUXDB_ORG := os.Getenv("INFLUXDB_ORG")
	INFLUXDB_URL := os.Getenv("INFLUXDB_URL")
	t.Logf("INFLUXDB_TOKEN: %s, INFLUXDB_BUCKET: %s, INFLUXDB_ORG: %s, INFLUXDB_URL: %s", INFLUXDB_TOKEN, INFLUXDB_BUCKET, INFLUXDB_ORG, INFLUXDB_URL)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	repo := &RepositoryTestEntityRepository{
		InfluxRepository: NewInfluxRepositoryImpl[RepositoryTestEntity](
			INFLUXDB_ORG, INFLUXDB_BUCKET, influxdb2.NewClient(INFLUXDB_URL, INFLUXDB_TOKEN), ctx, time.Second*5,
		),
	}
	_name := "test"
	_logstring := "this is log"
	_status := int64(http.StatusOK)
	_time := time.Now()
	entity := &RepositoryTestEntity{
		Name:      _name,
		LogString: &_logstring,
		Status:    &_status,
		Time:      &_time,
	}
	testTime := time.Now()
	if err := repo.Save(*entity); err != nil {
		t.Error(err)
	} else {
		t.Log("saved.")
	}

	if res, err := repo.FindAllByTimeAfter(testTime); err != nil {
		t.Error(err)
	} else if res == nil || len(res) < 1 {
		t.Errorf("result should return 1 or more rows but return nil or empty slice. err: %+v", err)
		t.Errorf("res: %+v", res)
	} else {
		t.Logf("res: %+v", res)
	}

	//if err := repo.DeleteAllByTimeAfter(testTime); err != nil {
	//	t.Error(err)
	//} else {
	//	t.Log("deleted.")
	//}

}

func Test_reflectSet(t *testing.T) {
	i8v := int64(math.MaxInt64)
	type S struct {
		I64 int8
	}
	s := S{}
	reflect.ValueOf(&s).Elem().FieldByIndex([]int{0}).SetInt(reflect.ValueOf(i8v).Int())
	t.Log(s.I64)

}
func Test_reflectSet2(t *testing.T) {
	i8v := int64(math.MaxInt64)
	type S struct {
		I64 int8
	}
	s := S{}
	reflect.ValueOf(&s).Elem().FieldByIndex([]int{0}).Set(reflect.ValueOf(i8v))
	t.Log(s.I64)

}
