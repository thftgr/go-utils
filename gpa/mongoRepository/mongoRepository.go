package mongoRepository

import (
	"context"
	"github.com/thftgr/go-utils/gpa"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoEntityId interface {
	gpa.Id
}

type MongoEntity[ID MongoEntityId] interface {
	gpa.Entity[ID]
	Collection() string
}

type MongoRepository[E MongoEntity[ID], ID MongoEntityId] interface {
	gpa.CrudRepository[E, ID]
}

type MongoRepositoryImpl[E MongoEntity[ID], ID MongoEntityId] struct {
	collection *mongo.Collection // 필수로 추가해야함
	ctx        context.Context
}

func NewMongoRepositoryImpl[E MongoEntity[ID], ID MongoEntityId](collection *mongo.Collection, ctx context.Context) *MongoRepositoryImpl[E, ID] {
	return &MongoRepositoryImpl[E, ID]{collection: collection, ctx: ctx}
}

func NewMongoRepository[E MongoEntity[ID], ID MongoEntityId](collection *mongo.Collection, ctx context.Context) MongoRepository[E, ID] {
	return &MongoRepositoryImpl[E, ID]{collection: collection, ctx: ctx}
}

func (m *MongoRepositoryImpl[E, ID]) Save(e E) (err error) {
	_, err = m.collection.UpdateOne(m.ctx, bson.M{"_id": e.GetId()}, bson.M{"$set": e}, options.Update().SetUpsert(true))
	return
}

func (m *MongoRepositoryImpl[E, ID]) SaveAll(e ...E) (count int64, err error) {
	for i := range e {
		res, err2 := m.collection.UpdateOne(m.ctx, bson.M{"_id": e[i].GetId()}, bson.M{"$set": e[i]}, options.Update().SetUpsert(true))
		if err2 != nil {
			err = err2
			break
		}
		count += res.ModifiedCount + res.UpsertedCount
	}
	return
}

func (m *MongoRepositoryImpl[E, ID]) FindById(id ID) (e E, err error) {
	err = m.collection.FindOne(m.ctx, bson.M{"_id": id}).Decode(&e)
	return
}

func (m *MongoRepositoryImpl[E, ID]) FindAllById(id ...ID) (e []E, err error) {
	cur, err := m.collection.Find(m.ctx, bson.M{"_id": bson.M{"$in": id}})
	if err != nil {
		return
	}
	defer cur.Close(m.ctx)
	err = cur.All(m.ctx, &e)
	return
}

func (m *MongoRepositoryImpl[E, ID]) Delete(e E) (err error) {
	_, err = m.collection.DeleteOne(m.ctx, bson.M{"_id": e.GetId()})
	return
}

func (m *MongoRepositoryImpl[E, ID]) DeleteAll(e ...E) (count int64, err error) {
	ids := make([]ID, len(e))
	for i := range e {
		ids[i] = e[i].GetId()
	}
	res, err := m.collection.DeleteMany(m.ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return
	}
	count = res.DeletedCount
	return
}

func (m *MongoRepositoryImpl[E, ID]) DeleteById(id ID) (err error) {
	_, err = m.collection.DeleteOne(m.ctx, bson.M{"_id": id})
	return
}

func (m *MongoRepositoryImpl[E, ID]) DeleteAllById(id ...ID) (count int64, err error) {
	res, err := m.collection.DeleteMany(m.ctx, bson.M{"_id": bson.M{"$in": id}})
	if err != nil {
		return
	}
	count = res.DeletedCount
	return
}
