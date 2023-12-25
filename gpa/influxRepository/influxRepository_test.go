package influxRepository

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/thftgr/go-utils/env"
	"net/http"
	"testing"
	"time"
)

type RepositoryTestEntity struct {
	Measurement `influxdb:"measurement:test"`
	Name        string    `influxdb:"tag:name"`
	LogString   string    `influxdb:"field:log"`
	Status      int       `influxdb:"field:status"`
	Time        time.Time `influxdb:"time"`
}

func (r RepositoryTestEntity) GetTime() time.Time {
	return r.Time
}

type RepositoryTestEntityRepository struct {
	InfluxRepository[RepositoryTestEntity]
}

func TestInfluxRepositoryImpl_ToPoint(t *testing.T) {
	ienv, _ := env.InfluxDB{}.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	repo := &RepositoryTestEntityRepository{
		InfluxRepository: NewInfluxRepositoryImpl[RepositoryTestEntity](
			ienv.ORG, ienv.BUCKET, influxdb2.NewClient(ienv.URL, ienv.TOKEN), ctx, time.Second*5,
		),
	}
	entity := &RepositoryTestEntity{
		Name:      "test",
		LogString: "this is log",
		Status:    http.StatusOK,
		Time:      time.Now(),
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
		t.Errorf("result should return 1 or more rows but return nil or empty slice.")
		t.Errorf("res: %+v", res)
	} else {
		t.Logf("res: %+v", res)
	}

}
