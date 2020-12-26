package controllers

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var testClient *mongo.Client

func init() { //TestMain(m *testing.M)
	log.Println("init() called.")

	//crete DB
	testClient, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://" + "127.0.0.1" + ":" + MongoDBPort))
	/*if err != nil {
		log.Fatal("error creating NewClient: ", err)
	}*/
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	err := testClient.Connect(context.Background())
	if err != nil {
		log.Fatal("error Connect()", err)
	}
	defer testClient.Disconnect(ctx)
	db := testClient.Database(Dbname)

	//create user in test DB
	err_createUser := db.RunCommand(ctx, bson.D{{"createUser", MongoUser},
		{"pwd", MongoPassword},
		{"roles", []bson.M{{"role": "readWrite", "db": "test"}}}})
	if err_createUser.Err() != nil {
		log.Println("error creating user: ", err_createUser)
	}
	err_createCollection_Job := db.CreateCollection(ctx, Colname)
	if err_createCollection_Job != nil {
		log.Println("error creating a collection Job: ", err_createCollection_Job)
	}
	err_createCollection_User := db.CreateCollection(ctx, ColnameUser)
	if err_createCollection_User != nil {
		log.Println("error creating collection User: ", err_createCollection_User)
	}
}

func TestInsertMongoDB(t *testing.T) {
	mongoClient, _ := ConnectMongoDB()
	//create test struct JsonJob
	jsonJob := new(JsonJob)
	jsonJob.URL = "https://leetcode.com/"
	jsonJob.Title = "TestTitle"
	jsonJob.Company = "TestCompany"
	jsonJob.DateAdded = time.Now().Format("2006-01-02")

	// 構造体をJSON文字列に変換
	jsonJobJSON, _ := json.Marshal(jsonJob)
	mongoClient.InsertMongoDB(jsonJobJSON, "Job")

	//check that the row is inserted.
	var result JsonJob
	collection := mongoClient.Client.Database(Dbname).Collection(Colname)
	filter := bson.D{{"company", jsonJob.Company}, {"title", jsonJob.Title}}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		t.Error("Expected insert TestCompany, TestTitle is successful.")
	}

}

func TestReadMongoDB(t *testing.T) {
	log.Println(os.Hostname())
	log.Println(os.Getwd())
	log.Println("os.Getenv:", os.Getenv("MONGO_SERVER"))
	mongoClient, err0 := ConnectMongoDB()
	if err0 != nil {
		log.Println("error from ConnectMongoDB")
		log.Println(err0)
	}
	//create expected result
	collection := mongoClient.Client.Database(Dbname).Collection(Colname)
	expect_result, err1 := collection.Find(context.Background(), bson.D{{"dateadded", bson.D{{"$gt", "2020-11-01"}}}})
	if err1 != nil {
		log.Println("error from collection.Find().")
		log.Println(err1)
	}
	var jobs []JsonJob
	var doc JsonJob
	for expect_result.Next(context.Background()) {
		//var doc JsonJob
		err := expect_result.Decode(&doc)
		if err != nil {
			log.Println("error from expect_result.Decode(&doc).")
			return
		}
		jobs = append(jobs, doc)
	}

	//actual result
	result := mongoClient.ReadMongo("", false, false, false)
	if len(result) != len(jobs) {
		t.Error("Expected:", len(jobs), "but actually:", len(result))
	}
}
