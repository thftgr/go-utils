package gpa

import "time"

type Id interface {
	comparable
}

type Entity[ID Id] interface {
	GetId() ID
}

type TimeSeriesEntity interface {
	Time() time.Time
}
