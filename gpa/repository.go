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

// TimeSeriesRepository 시계열 데이터는 기본적으로 범위데이터임.
// 일반적인 사용 케이스에서는 삽입, 조회를 주로 사용함.
// 삭제는 잘 사용하지 않음으로 정의하지 않았음.
type TimeSeriesRepository[E TimeSeriesEntity] interface {
	Save(E) error // async 의경우 error없음.
}
