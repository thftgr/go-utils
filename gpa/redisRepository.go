package gpa

type RedisEntityId interface {
	Id
}

type RedisEntity[ID RedisEntityId] interface {
	Entity[ID]
	KeyPrefix() string
}

type RedisRepository[E RedisEntity[ID], ID RedisEntityId] interface {
	//CrudRepository[E, ID]
}

//type RedisRepositoryImpl[E RedisEntity[ID], ID RedisEntityId] struct {
//	Model E             // 제너릭을 통해 자동 설정됨.
//	DB    *redis.Client // 필수로 추가해야함
//}
