package influxRepository

import (
	"context"
	"errors"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
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
	gpa.TimeSeriesEntity
}

type InfluxRepository[E InfluxEntity] interface {
	gpa.TimeSeriesRepository[E]
	FindAllByTime(time.Time, time.Time) ([]E, error)             // |<---------->|
	FindAllByDuration(time.Duration, time.Duration) ([]E, error) // |<---------->|
	FindAllByTimeAfter(time.Time) ([]E, error)                   // |--------------->latest
	FindAllByDurationAfter(time.Duration) ([]E, error)           // |--------------->latest
	DeleteAllByTimeAfter(time.Time) error
}

type InfluxRepositoryImpl[E InfluxEntity] struct {
	Org         string
	Bucket      string
	DB          influxdb2.Client
	WriteAPI    api.WriteAPI
	QueryAPI    api.QueryAPI
	DeleteAPI   api.DeleteAPI // https://docs.influxdata.com/influxdb/v2/reference/syntax/delete-predicate/ 참고
	Context     context.Context
	EntityCache *InfluxEntityTagHelper[E]
	Timeout     time.Duration
}

func NewInfluxRepositoryImpl[E InfluxEntity](org string, bucket string, DB influxdb2.Client, context context.Context, timeout time.Duration) *InfluxRepositoryImpl[E] {
	return &InfluxRepositoryImpl[E]{
		Org:         org,
		Bucket:      bucket,
		DB:          DB,
		WriteAPI:    DB.WriteAPI(org, bucket),
		QueryAPI:    DB.QueryAPI(org),
		DeleteAPI:   DB.DeleteAPI(),
		Context:     context,
		EntityCache: NewInfluxEntityTagHelper[E](),
		Timeout:     timeout,
	}
}

// ToPoint E 가 InfluxEntityEncoder 를 구현하지 않은경우 reflection 으로 처리됨. 이로인해 성능문제가 발생할수있음.
// If ToPoint E does not implement InfluxEntityEncoder, it will be treated as reflection, which may cause performance problems.
func (r *InfluxRepositoryImpl[E]) ToPoint(e E) (p *write.Point) {
	if encoder, ok := (any(e)).(InfluxEntityEncoder); ok {
		p = write.NewPointWithMeasurement(encoder.GetMeasurement())
		for _, t := range encoder.GetTags() {
			p.AddTag(t.Key, t.Value)
		}
		for _, f := range encoder.GetField() {
			p.AddField(f.Key, f.Value)
		}
		p.SetTime(encoder.GetTime())
	} else {
		return r.EntityCache.ToPoint(&e)
	}

	return
}

func (r *InfluxRepositoryImpl[E]) FromPoints(rows *api.QueryTableResult) (res []E, err error) {
	for rows.Next() {
		record := rows.Record()
		var entity E
		if decoder, ok := (any(entity)).(InfluxEntityDecoder); ok {
			decoder.SetTime(record.Time())
			if err = decoder.SetValue(record.Values()); err != nil {
				return nil, err
			}
		} else {
			_, err = r.EntityCache.FromRows(rows)
		}
		res = append(res, entity)
	}
	return
}

func (r *InfluxRepositoryImpl[E]) Save(e E) error {
	//c := r.WriteAPI.Errors()
	r.WriteAPI.WritePoint(r.ToPoint(e))
	r.WriteAPI.Flush()
	return nil
}

func (r *InfluxRepositoryImpl[E]) FindAllByTime(start time.Time, stop time.Time) ([]E, error) {
	panic("")
}
func (r *InfluxRepositoryImpl[E]) FindAllByDuration(start time.Duration, stop time.Duration) ([]E, error) {
	panic("")
}
func (r *InfluxRepositoryImpl[E]) FindAllByTimeAfter(start time.Time) (res []E, err error) {
	ctx, cancel := context.WithTimeout(r.Context, r.Timeout)
	defer cancel()
	rows, err := r.QueryAPI.Query(ctx, fmt.Sprintf(`
from(bucket: "%s")
	|> range(start: %s)
	|> filter(fn: (r) => r._measurement == "%s")
    |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
	`, r.Bucket, start.Format(time.RFC3339), r.EntityCache.measurement),
	)
	if err != nil {
		return
	}
	defer func(rows *api.QueryTableResult) {
		e := rows.Close()
		if e != nil {
			err = errors.Join(err, e)
		}
	}(rows)
	return r.EntityCache.FromRows(rows)
}
func (r *InfluxRepositoryImpl[E]) FindAllByDurationAfter(time.Duration) ([]E, error) {
	panic("")
}

func (r *InfluxRepositoryImpl[E]) DeleteAllByTimeAfter(start time.Time) (err error) {
	ctx, cancel := context.WithTimeout(r.Context, r.Timeout)
	defer cancel()
	return r.DeleteAPI.DeleteWithName(ctx, r.Org, r.Bucket, start, time.Now(), "")
}
