package influxRepository

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	protocol "github.com/influxdata/line-protocol"
	"github.com/thftgr/go-utils/utils"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

type RepositoryTestEntity struct {
	Measurement `json:"-"  influxdb:"measurement:test"`
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

func TestInfluxRepositoryImpl_FullTest(t *testing.T) {
	INFLUXDB_TOKEN := os.Getenv("INFLUXDB_TOKEN")
	INFLUXDB_BUCKET := os.Getenv("INFLUXDB_BUCKET")
	INFLUXDB_ORG := os.Getenv("INFLUXDB_ORG")
	INFLUXDB_URL := os.Getenv("INFLUXDB_URL")
	//t.Logf("INFLUXDB_TOKEN:  %s", INFLUXDB_TOKEN)
	t.Logf("INFLUXDB_BUCKET: %s", INFLUXDB_BUCKET)
	t.Logf("INFLUXDB_ORG:    %s", INFLUXDB_ORG)
	t.Logf("INFLUXDB_URL:    %s", INFLUXDB_URL)
	now := time.Now().UTC()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	repo := &RepositoryTestEntityRepository{
		InfluxRepository: NewInfluxRepositoryImpl[RepositoryTestEntity](
			INFLUXDB_ORG, INFLUXDB_BUCKET, influxdb2.NewClient(INFLUXDB_URL, INFLUXDB_TOKEN), ctx, time.Second*5,
		),
	}
	var testTags = []*protocol.Tag{
		{Key: "name", Value: "test-name"},
	}
	var testEntities = []RepositoryTestEntity{
		{Name: "test-name", LogString: "log string1", Status: http.StatusOK, Time: now},
		{Name: "test-name", LogString: "log string2", Status: http.StatusOK, Time: now.Add(-time.Second)},
	}

	for i := range testEntities { // 테스트 데이터 삽입 겸 save 테스트
		if err := repo.SaveAndFlush(testEntities[i]); err != nil {
			t.Error(err)
		}
	}
	utils.ReverseSlice[RepositoryTestEntity](testEntities)
	var startT, stopT = now.Add(-(time.Second * 5)), now.Add(time.Second)
	var startD, stopD = now.Sub(now.Add(-(time.Second * 5))), time.Second

	if e, err := repo.FindAllByTime(startT, stopT); err != nil { // -5s ~ now
		t.Error(err)
	} else if !reflect.DeepEqual(e, testEntities) {
		t.Errorf("FindAllByTime() gotRes = %v, want %v", e, testEntities)
	}

	if e, err := repo.FindAllByDuration(startD, stopD); err != nil { // -5s ~ now
		t.Error(err)
	} else if !reflect.DeepEqual(e, testEntities[:1]) {
		t.Errorf("FindAllByDuration() gotRes = %v, want %v", e, testEntities[:1])
	}

	if e, err := repo.FindAllByTagsAndTime(startT, stopT, testTags); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(e, testEntities) {
		t.Errorf("FindAllByTagsAndTime() gotRes = %v, want %v", e, testEntities)
	}

	if e, err := repo.FindAllByTagsAndDuration(startD, stopD, testTags); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(e, testEntities[:1]) {
		t.Errorf("FindAllByTagsAndDuration() gotRes = %v, want %v", e, testEntities[:1])
	}

	if e, err := repo.FindAllByTimeAfter(startT); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(e, testEntities) {
		t.Errorf("FindAllByTimeAfter() gotRes = %v, want %v", e, testEntities)
	}

	if e, err := repo.FindAllByDurationAfter(startD); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(e, testEntities) {
		t.Errorf("FindAllByDurationAfter() gotRes = %v, want %v", e, testEntities)
	}

	if e, err := repo.FindAllByTagsAndTimeAfter(startT, testTags); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(e, testEntities) {
		t.Errorf("FindAllByTagsAndTimeAfter() gotRes = %v, want %v", e, testEntities)
	}

	if e, err := repo.FindAllByTagsAndDurationAfter(startD, testTags); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(e, testEntities) {
		t.Errorf("FindAllByTagsAndDurationAfter() gotRes = %v, want %v", e, testEntities)
	}

}
