package gpa

type Repository[E Entity[ID], ID Id] interface {
}

type CrudRepository[E Entity[ID], ID Id] interface {
	Repository[E, ID] // implements

	//C,U

	Save(E) error
	SaveAll(...E) (int64, error)

	//R

	FindById(ID) (E, error)
	FindAllById(...ID) ([]E, error)

	//D

	Delete(E) error
	DeleteAll(...E) (int64, error)
	DeleteById(ID) error
	DeleteAllById(...ID) (int64, error)
}
