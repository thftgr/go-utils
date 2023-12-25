package influxRepository

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	protocol "github.com/influxdata/line-protocol"
	"github.com/thftgr/go-utils/gpa"
	"time"
)

// InfluxEntity
// RDB              || influxdb
// ===========================================
// Database         || Database(v1), Bucket(v2)
// Table            || Measurement
// Column           || Tag, Field
// Indexed Column   || Tag (String Only)
// Unindexed Column || Field
// Row              || Point
// ===========================================
//
// The influxdb usually uses "snake_case" if field is PascalCase use "${tag|field}:${snake_case_name}"
// The InfluxEntity implementation structure must have fields that contain the tags below.
//
// * `influxdb:"measurement:${measurement_name}"`
// * `influxdb:"tag:${tag_name}"`
// * `influxdb:"field:${field_name}"`
// * `influxdb:"time"`
type InfluxEntity interface {
	InfluxEntityDecoder
	InfluxEntityEncoder
}

type InfluxRepository[E InfluxEntity] interface {
	gpa.TimeSeriesRepository[E]
	FindAllByTime(time.Time, time.Time, ...string) ([]E, error)             // |<---------->|
	FindAllByDuration(time.Duration, time.Duration, ...string) ([]E, error) // |<---------->|
	FindAllByTimeAfter(time.Time, ...string) ([]E, error)                   // |--------------->latest
	FindAllByDurationAfter(time.Duration, ...string) ([]E, error)           // |--------------->latest
}

type InfluxRepositoryImpl[E InfluxEntity] struct {
	Org         string
	Bucket      string
	DB          influxdb2.Client
	WriteAPI    api.WriteAPI
	QueryAPI    api.QueryAPI
	DeleteAPI   api.DeleteAPI
	Context     context.Context
	EntityCache any
}

func NewInfluxRepositoryImpl[E InfluxEntity](org string, bucket string, DB influxdb2.Client, context context.Context) *InfluxRepositoryImpl[E] {
	return &InfluxRepositoryImpl[E]{
		Org:         org,
		Bucket:      bucket,
		DB:          DB,
		WriteAPI:    DB.WriteAPI(org, bucket),
		QueryAPI:    DB.QueryAPI(org),
		DeleteAPI:   DB.DeleteAPI(),
		Context:     context,
		EntityCache: nil,
	}
}

func (r *InfluxRepositoryImpl[E]) ToPoint(e E) (p *write.Point) {
	p = write.NewPointWithMeasurement(e.GetMeasurement())
	for _, t := range e.GetTags() {
		p.AddTag(t.Key, t.Value)
	}
	for _, f := range e.GetField() {
		p.AddField(f.Key, f.Value)
	}
	p.SetTime(e.GetTime())
	return
}

func (r *InfluxRepositoryImpl[E]) FromPoints(rows *api.QueryTableResult) (res []E, err error) {
	for rows.Next() {
		record := rows.Record()
		var entity E
		entity.SetTime(record.Time())
		if err = entity.SetValue(record.Values()); err != nil {
			return nil, err
		}
		res = append(res, entity)
	}
	return
}

func (r *InfluxRepositoryImpl[E]) Save(e E) error {
	c := r.WriteAPI.Errors()
	r.WriteAPI.WritePoint(r.ToPoint(e))
	return <-c
}

func (r *InfluxRepositoryImpl[E]) FindAllByTime(start time.Time, stop time.Time, tag ...protocol.Tag) (res []E, err error) {
	qs := ``
	rows, err := r.QueryAPI.QueryWithParams(r.Context, qs, []string{})
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.FromPoints(rows)
}

func (r *InfluxRepositoryImpl[E]) FindAllByDuration(start time.Duration, stop time.Duration) ([]E, error) {
	//TODO implement me
	panic("implement me")
}

func (r *InfluxRepositoryImpl[E]) FindAllByTimeAfter(start time.Time) ([]E, error) {
	//TODO implement me
	panic("implement me")
}

func (r *InfluxRepositoryImpl[E]) FindAllByDurationAfter(start time.Duration) ([]E, error) {
	//TODO implement me
	panic("implement me")
}
