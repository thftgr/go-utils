package gpa

type Id interface {
}

type Entity[ID Id] interface {
	GetId() ID
}
