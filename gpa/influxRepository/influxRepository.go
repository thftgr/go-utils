package influxRepository

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	protocol "github.com/influxdata/line-protocol"
	"github.com/thftgr/go-utils/gpa"
	"strings"
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
	SaveAndFlush(E) error

	FindAllByTime(time.Time, time.Time) ([]E, error)
	FindAllByDuration(time.Duration, time.Duration) ([]E, error)
	FindAllByTagsAndTime(time.Time, time.Time, []*protocol.Tag) ([]E, error)
	FindAllByTagsAndDuration(time.Duration, time.Duration, []*protocol.Tag) ([]E, error)

	//DeleteAllByTime(time.Time, time.Time) error
	//DeleteAllByDuration(time.Duration, time.Duration) error
	//DeleteAllByTagsTime(time.Time, time.Time, []*protocol.Tag) error
	//DeleteAllByTagsDuration(time.Duration, time.Duration, []*protocol.Tag) error
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

// FromPoints E 가 InfluxEntityDecoder 를 구현하지 않은경우 reflection 으로 처리됨. 이로인해 성능문제가 발생할수있음.
// If FromPoints E does not implement InfluxEntityDecoder, it will be treated as reflection, which may cause performance problems.
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

// TimeBaseQueryBuilder
//
//	from(bucket: "%s")
//	    |> range(start: %s, stop:%s)
//	    |> filter(fn: (r) => r._measurement == "%s" and ...)
//	    |> keep(columns: [%s])
//	    |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
func (r *InfluxRepositoryImpl[E]) TimeBaseQueryBuilder(start time.Time, stop time.Time, tags []*protocol.Tag, columns ...string) string {
	buf := bytes.Buffer{}
	_, _ = fmt.Fprintf(&buf, `from(bucket: "%s")`, r.Bucket)
	_, _ = fmt.Fprintf(&buf, `  |> range(start: %s, stop: %s)`, start.Format(time.RFC3339), stop.Format(time.RFC3339))

	var opers = []string{fmt.Sprintf(`r._measurement == "%s"`, r.EntityCache.measurement)}
	for _, t := range tags {
		opers = append(opers, fmt.Sprintf(`r.%s == "%s"`, t.Key, t.Value))
	}
	_, _ = fmt.Fprintf(&buf, `  |> filter(fn: (r) => %s)`, strings.Join(opers, " and "))

	if len(columns) > 0 {
		for i := range columns {
			columns[i] = `"` + columns[i] + `"`
		}
		_, _ = fmt.Fprintf(&buf, `  |> keep(columns: [%s])`, strings.Join(columns, ", "))
	}

	_, _ = fmt.Fprint(&buf, `  |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")`)
	return buf.String()
}

// DurationBaseQueryBuilder
//
//	from(bucket: "%s")
//	    |> range(start: -%s, stop: -%s)
//	    |> filter(fn: (r) => r._measurement == "%s" and ...)
//	    |> keep(columns: [%s])
//	    |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
func (r *InfluxRepositoryImpl[E]) DurationBaseQueryBuilder(start time.Duration, stop time.Duration, tags []*protocol.Tag, columns ...string) string {
	buf := bytes.Buffer{}
	_, _ = fmt.Fprintf(&buf, `from(bucket: "%s")`, r.Bucket)
	_, _ = fmt.Fprintf(&buf, `  |> range(start: -%s, stop: -%s)`, start, stop)

	var opers = []string{fmt.Sprintf(`r._measurement == "%s"`, r.EntityCache.measurement)}
	for _, t := range tags {
		opers = append(opers, fmt.Sprintf(`r.%s == "%s"`, t.Key, t.Value))
	}
	_, _ = fmt.Fprintf(&buf, `  |> filter(fn: (r) => %s)`, strings.Join(opers, " and "))

	if len(columns) > 0 {
		for i := range columns {
			columns[i] = `"` + columns[i] + `"`
		}
		_, _ = fmt.Fprintf(&buf, `  |> keep(columns: [%s])`, strings.Join(columns, ", "))
	}

	_, _ = fmt.Fprint(&buf, `  |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")`)
	return buf.String()
}

func (r *InfluxRepositoryImpl[E]) PredicateBuilder(tags []*protocol.Tag) string {
	var opers = []string{fmt.Sprintf(`"_measurement"="%s"`, r.EntityCache.measurement)}
	for _, t := range tags {
		opers = append(opers, fmt.Sprintf(`"%s"="%s"`, t.Key, t.Value))
	}
	return strings.Join(opers, " and ")
}

func (r *InfluxRepositoryImpl[E]) QueryByTime(start time.Time, stop time.Time, tags []*protocol.Tag, fields ...string) (res []E, err error) {
	ctx, cancel := context.WithTimeout(r.Context, r.Timeout)
	defer cancel()
	rows, err := r.QueryAPI.Query(ctx, r.TimeBaseQueryBuilder(start, stop, tags, fields...))
	if err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()
	return r.EntityCache.FromRows(rows)
}

func (r *InfluxRepositoryImpl[E]) QueryByDuration(start time.Duration, stop time.Duration, tags []*protocol.Tag, fields ...string) (res []E, err error) {
	ctx, cancel := context.WithTimeout(r.Context, r.Timeout)
	defer cancel()
	rows, err := r.QueryAPI.Query(ctx, r.DurationBaseQueryBuilder(start, stop, tags, fields...))
	if err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()
	return r.EntityCache.FromRows(rows)
}

func (r *InfluxRepositoryImpl[E]) DeleteByTime(start time.Time, stop time.Time, tags []*protocol.Tag) (err error) {
	ctx, cancel := context.WithTimeout(r.Context, r.Timeout)
	defer cancel()
	return r.DeleteAPI.DeleteWithName(ctx, r.Org, r.Bucket, start, stop, r.PredicateBuilder(tags))
}

func (r *InfluxRepositoryImpl[E]) DeleteByDuration(start time.Duration, stop time.Duration, tags []*protocol.Tag) (err error) {
	now := time.Now()
	return r.DeleteByTime(now.Add(-start), now.Add(-stop), tags)
}

// ====================================================================================================================
// ====================================================================================================================
// ====================================================================================================================

func (r *InfluxRepositoryImpl[E]) Save(e E) error {
	r.WriteAPI.WritePoint(r.ToPoint(e))
	return nil // 비동기라 error가 발생하지 않음.
}

func (r *InfluxRepositoryImpl[E]) SaveAndFlush(e E) error {
	errCh := r.WriteAPI.Errors()
	r.WriteAPI.WritePoint(r.ToPoint(e))
	r.WriteAPI.Flush()
	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func (r *InfluxRepositoryImpl[E]) FindAllByTime(start time.Time, stop time.Time) (res []E, err error) {
	return r.QueryByTime(start, stop, nil)
}
func (r *InfluxRepositoryImpl[E]) FindAllByDuration(start time.Duration, stop time.Duration) (res []E, err error) {
	return r.QueryByDuration(start, stop, nil)
}
func (r *InfluxRepositoryImpl[E]) FindAllByTagsAndTime(start time.Time, stop time.Time, tags []*protocol.Tag) (res []E, err error) {
	return r.QueryByTime(start, stop, tags)
}
func (r *InfluxRepositoryImpl[E]) FindAllByTagsAndDuration(start time.Duration, stop time.Duration, tags []*protocol.Tag) (res []E, err error) {
	return r.QueryByDuration(start, stop, tags)
}

func (r *InfluxRepositoryImpl[E]) DeleteAllByTime(start time.Time, stop time.Time) error {
	return r.DeleteByTime(start, stop, nil)
}
func (r *InfluxRepositoryImpl[E]) DeleteAllByDuration(start time.Duration, stop time.Duration) error {
	return r.DeleteByDuration(start, stop, nil)
}
func (r *InfluxRepositoryImpl[E]) DeleteAllByTagsTime(start time.Time, stop time.Time, tags []*protocol.Tag) error {
	return r.DeleteByTime(start, stop, tags)
}
func (r *InfluxRepositoryImpl[E]) DeleteAllByTagsDuration(start time.Duration, stop time.Duration, tags []*protocol.Tag) error {
	return r.DeleteByDuration(start, stop, tags)
}
