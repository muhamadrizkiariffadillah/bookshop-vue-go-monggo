package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client
var defaulDbName = "test_tb"

func init() {

	LoadEnvVariable()

}

func InitDatabase() (*mongo.Client, error) {

	dbUrl := GetEnvProperties("MONGODB_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUrl))

	if err != nil {
		log.Fatal("db_config error :", err.Error())
	}

	fmt.Println("database connection successfully")

	return client, nil
}

func GetDatabaseCollection(dbName *string, collectionName string) *mongo.Collection {

	if *dbName == "" || dbName != nil {
		dbName = &defaulDbName
	}

	if dbClient == nil {
		dbClient, _ = InitDatabase()
	}

	collection := dbClient.Database(*dbName).Collection(collectionName)

	return collection
}

func InitDb() (*mongo.Client, error) {
	db, err := InitDatabase()

	if err != nil {
		log.Fatal("error init_db: ", err.Error())
	}

	return db, nil
}
