package influxrepository

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
	*InfluxRepositoryImpl[RepositoryTestEntity]
}

func TestInfluxRepositoryImpl_FullFindTest(t *testing.T) {
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
		InfluxRepositoryImpl: NewInfluxRepositoryImpl[RepositoryTestEntity](
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
	var startD, stopD = now.Sub(now.Add(-(time.Second * 5))), time.Duration(0)
	t.Logf("now: %+v, startT: %+v, stopT: %+v, startD: %+v, stopD: %+v", now, startT, stopT, startD, stopD)

	if e, err := repo.FindAllByTime(startT, stopT); err != nil { // -5s ~ now
		t.Error(err)
	} else if !reflect.DeepEqual(e, testEntities) {
		t.Errorf("FindAllByTime() gotRes = %v, want %v", e, testEntities)
	}

	if e, err := repo.FindAllByDuration(startD, stopD); err != nil { // -5s ~ now
		t.Error(err)
	} else if !reflect.DeepEqual(e, testEntities) {
		t.Errorf("FindAllByDuration() gotRes = %v, want %v", e, testEntities)
	}

	if e, err := repo.FindAllByTagsAndTime(startT, stopT, testTags); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(e, testEntities) {
		t.Errorf("FindAllByTagsAndTime() gotRes = %v, want %v", e, testEntities)
	}

	if e, err := repo.FindAllByTagsAndDuration(startD, stopD, testTags); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(e, testEntities) {
		t.Errorf("FindAllByTagsAndDuration() gotRes = %v, want %v", e, testEntities)
	}

	if err := repo.DeleteByTime(startT, stopT, nil); err != nil {
		t.Error(err)
	}

}

func TestInfluxRepositoryImpl_FullDeleteTest(t *testing.T) {
	INFLUXDB_TOKEN := os.Getenv("INFLUXDB_TOKEN")
	INFLUXDB_BUCKET := os.Getenv("INFLUXDB_BUCKET")
	INFLUXDB_ORG := os.Getenv("INFLUXDB_ORG")
	INFLUXDB_URL := os.Getenv("INFLUXDB_URL")

	//t.Logf("INFLUXDB_TOKEN:  %s", INFLUXDB_TOKEN)
	t.Logf("INFLUXDB_BUCKET: %s", INFLUXDB_BUCKET)
	t.Logf("INFLUXDB_ORG:    %s", INFLUXDB_ORG)
	t.Logf("INFLUXDB_URL:    %s", INFLUXDB_URL)
	now := time.Now().UTC()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	repo := &RepositoryTestEntityRepository{
		InfluxRepositoryImpl: NewInfluxRepositoryImpl[RepositoryTestEntity](
			INFLUXDB_ORG, INFLUXDB_BUCKET, influxdb2.NewClient(INFLUXDB_URL, INFLUXDB_TOKEN), ctx, time.Second*600,
		),
	}
	var testTags = []*protocol.Tag{
		{Key: "name", Value: "test-name"},
	}
	var testEntities = []RepositoryTestEntity{
		{Name: "test-name", LogString: "log string-5s", Status: http.StatusOK, Time: now.Add(-(time.Second * 5))},
		{Name: "test-name", LogString: "log string+0s", Status: http.StatusOK, Time: now},
		{Name: "test-name", LogString: "log string+5s", Status: http.StatusOK, Time: now.Add(time.Second * 5)},
	}

	var startT, stopT = now.Add(-(time.Second * 5)), now.Add(time.Second)
	var startD, stopD = now.Sub(now.Add(-(time.Second * 5))), time.Duration(0)
	_, _, _, _, _ = startT, stopT, startD, stopD, testTags

	for i := range testEntities {
		if err := repo.SaveAndFlush(testEntities[i]); err != nil {
			t.Error(err)
		}
	}
	if err := repo.DeleteAllByTime(startT, stopT); err != nil {
		t.Error(err)
	} else {
		if e, err := repo.FindAllByTime(startT, stopT); err != nil {
			t.Error(err)
		} else if len(e) > 0 {
			t.Errorf("FindAllByTagsAndTimeAfter() gotRes = %v, want = []", e)
		}
	}
	if err := repo.DeleteByTime(now.Add(-(time.Second * 6)), now.Add(time.Second*6), nil); err != nil {
		t.Error(err)
	}

	for i := range testEntities {
		if err := repo.SaveAndFlush(testEntities[i]); err != nil {
			t.Error(err)
		}
	}
	if err := repo.DeleteAllByTime(now, stopT); err != nil {
		t.Error(err)
	} else {
		if e, err := repo.FindAllByTime(startT, stopT); err != nil {
			t.Error(err)
		} else if len(e) != 1 {
			t.Errorf("FindAllByTagsAndTimeAfter() gotRes = %v, want = []", e)
		}
	}
	if err := repo.DeleteByTime(now.Add(-(time.Second * 6)), now.Add(time.Second*6), nil); err != nil {
		t.Error(err)
	}

	for i := range testEntities {
		if err := repo.SaveAndFlush(testEntities[i]); err != nil {
			t.Error(err)
		}
	}
	if err := repo.DeleteAllByDuration(startD, stopD); err != nil {
		t.Error(err)
	} else {
		if e, err := repo.FindAllByDuration(startD, stopD); err != nil {
			t.Error(err)
		} else if len(e) > 0 {
			t.Errorf("FindAllByTagsAndTimeAfter() gotRes = %v, want = []", e)
		}
	}
	if err := repo.DeleteByTime(now.Add(-(time.Second * 6)), now.Add(time.Second*6), nil); err != nil {
		t.Error(err)
	}

}
