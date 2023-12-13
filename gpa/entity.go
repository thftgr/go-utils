package gpa

type Id interface {
	comparable
}

type Entity[ID Id] interface {
	GetId() ID
}
