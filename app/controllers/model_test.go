package controllers

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

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
	result := mongoClient.ReadMongo()
	if len(result) != len(jobs) {
		t.Error("Expected:", len(jobs), "but actually:", len(result))
	}
}
