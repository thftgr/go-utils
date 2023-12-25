package influxRepository

import (
	protocol "github.com/influxdata/line-protocol"
	"testing"
	"time"
)

type RepositoryTestEntity struct {
	Measurement `influxdb:"measurement:test"`
	T1          string    `influxdb:"tag:tag_name"`
	F2          string    `influxdb:"field:field_name"`
	Time        time.Time `influxdb:"time"`
}

func (r *RepositoryTestEntity) SetValue(values map[string]interface{}) error {
	for k, v := range values {
		switch k {
		case "t1":
			r.T1, _ = v.(string)
		case "f2":
			r.F2, _ = v.(string)
		}
	}
	return nil
}

func (r *RepositoryTestEntity) SetTime(t time.Time) {
	r.Time = t
}

func (r *RepositoryTestEntity) GetMeasurement() string {
	return "dev"
}

func (r *RepositoryTestEntity) GetTags() []*protocol.Tag {
	return []*protocol.Tag{
		{"t1", r.T1},
	}
}

func (r *RepositoryTestEntity) GetField() []*protocol.Field {
	return []*protocol.Field{
		{"f2", r.F2},
	}
}

func (r *RepositoryTestEntity) GetTime() time.Time {
	return r.Time
}

type RepositoryTestEntityRepository struct {
	InfluxRepository[RepositoryTestEntity]
}

func TestInfluxRepositoryImpl_ToPoint(t *testing.T) {
	//ienv, _ := env.InfluxDB{}.Parse()
	//client := influxdb2.NewClient(ienv.URL, ienv.TOKEN)
	//influxRepository.InfluxRepositoryImpl[Test]{}
}
