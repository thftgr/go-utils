package influxRepository

import "github.com/thftgr/go-utils/gpa"

type InfluxEntityId interface {
	gpa.Id
}

type InfluxEntity[ID InfluxEntityId] interface {
	gpa.Entity[ID]
}

type InfluxRepository[E InfluxEntity[ID], ID InfluxEntityId] interface {
}

type InfluxRepositoryImpl[E InfluxEntity[ID], ID InfluxEntityId] struct {
}
