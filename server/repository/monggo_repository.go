package repository

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/muhamadrizkiariffadillah/bookshop-vue-go-monggo/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepositoryInterface interface {
	Create(data any, ctx mongo.SessionContext) (any, error)
	FindOne(id string, ctx mongo.SessionContext) (fiber.Map, error)
	Update(id string, data any, ctx mongo.SessionContext) (*mongo.UpdateResult, error)
	Delete(id string, ctx mongo.SessionContext) (*mongo.DeleteResult, error)
	FindAll(ctx mongo.SessionContext) ([]fiber.Map, error)
	Aggregate(pipeline mongo.Pipeline, ctx mongo.SessionContext) ([]fiber.Map, error)
}

type MongoRepository struct {
	collection *mongo.Collection
}

func setUpSessionContext(sessionContext mongo.SessionContext) mongo.SessionContext {
	if sessionContext == nil {
		return mongo.NewSessionContext(context.Background(), nil)
	}
	return sessionContext
}

func GetMongoRepository(dbName, collectionName string) *MongoRepository {
	collection := config.GetDatabaseCollection(&dbName, collectionName)

	return &MongoRepository{collection: collection}
}

func (mr *MongoRepository) Create(data any, ctx mongo.SessionContext) (any, error) {
	sessionCtx := setUpSessionContext(ctx)

	result, err := mr.collection.InsertOne(sessionCtx, data)
	if err != nil {
		return nil, errors.New("mongo repository: error inserting document")
	}

	return result, nil
}

func (mr *MongoRepository) FindOne(id string, sessionCtx mongo.SessionContext) (fiber.Map, error) {
	if sessionCtx == nil {
		return nil, errors.New("mongo repository: session context is nil")
	}

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("mongo repository: invalid object ID")
	}

	var document fiber.Map

	err = mr.collection.FindOne(sessionCtx, bson.M{"_id": objId}).Decode(&document)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("mongo repository: document not found")
		}
		return nil, errors.New("mongo repository: error finding document")
	}

	return document, nil
}

func (mr *MongoRepository) Update(id string, data any, ctx mongo.SessionContext) (*mongo.UpdateResult, error) {

	sessionCtx := setUpSessionContext(ctx)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("mongo repository: invalid object ID")
	}

	updateData := bson.M{"$set": data}

	result, err := mr.collection.UpdateOne(sessionCtx, bson.M{"_id": objectId}, updateData)
	if err != nil {
		return nil, errors.New("mongo repository: error updating document")
	}

	return result, nil
}

func (mr *MongoRepository) Delete(id string, ctx mongo.SessionContext) (*mongo.DeleteResult, error) {
	sessionCtx := setUpSessionContext(ctx)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("mongo repository: invalid object ID")
	}

	result, err := mr.collection.DeleteOne(sessionCtx, bson.M{"_id": objectId})
	if err != nil {
		return nil, errors.New("mongo repository: error deleting document")
	}

	if result.DeletedCount == 0 {
		return nil, errors.New("mongo repository: document not found or already deleted")
	}

	return result, nil
}

func (mr *MongoRepository) FindAll(ctx mongo.SessionContext) ([]fiber.Map, error) {
	sessionCtx := setUpSessionContext(ctx)

	cursor, err := mr.collection.Find(sessionCtx, bson.M{})
	if err != nil {
		return nil, errors.New("mongo repository: error fetching documents")
	}
	defer cursor.Close(sessionCtx)

	var results []fiber.Map
	for cursor.Next(sessionCtx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, errors.New("mongo repository: error decoding document")
		}

		results = append(results, fiber.Map(doc))
	}

	if len(results) == 0 {
		return nil, errors.New("mongo repository: no documents found")
	}

	return results, nil
}

func (mr *MongoRepository) Aggregate(pipeline mongo.Pipeline, ctx mongo.SessionContext) ([]fiber.Map, error) {
	sessionCtx := setUpSessionContext(ctx)

	cursor, err := mr.collection.Aggregate(sessionCtx, pipeline)
	if err != nil {
		return nil, errors.New("mongo repository: error executing aggregation pipeline")
	}
	defer cursor.Close(sessionCtx)

	var results []fiber.Map
	for cursor.Next(sessionCtx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, errors.New("mongo repository: error decoding aggregation result")
		}
		results = append(results, fiber.Map(doc))
	}

	if len(results) == 0 {
		return nil, errors.New("mongo repository: no aggregation results found")
	}

	return results, nil
}
