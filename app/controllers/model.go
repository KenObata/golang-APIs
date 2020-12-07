package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB構造体へInsert用のメソッドを定義
// JSONファイルから読み込んだバイトスライスを渡し、MongoDBへInsert
func (db *DB) InsertMongoDB(json []byte, table_name string) {
	log.Println("InsertMongoDB is called.")
	// 60秒でタイムアウトするコンテキストを作成
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	bsonMap := bson.M{}
	// JSONのバイトスライスをMongoDBのドキュメント型であるbsonへマップ
	err := bson.UnmarshalExtJSON([]byte(json), false, &bsonMap)
	if err != nil {
		log.Println("error from InsertMongo(), bson.UnmarshalExtJSON")
		log.Println(err)
		return
	}
	// Insert先のコレクション名からクライアント作成
	collection := db.Client.Database(Dbname).Collection(table_name)
	//log.Println("Mongo DB name:", db.Client.Database(Dbname).Name())
	if table_name == Colname { //table_name==Job
		var episodesFiltered JsonJob
		filter := bson.D{{"company", bsonMap["company"]}, {"title", bsonMap["title"]}} //bsonMap["url"]
		err := collection.FindOne(context.Background(), filter).Decode(&episodesFiltered)
		if err != nil {
			//log.Println("Error from collection.Find.")
			//log.Println(err)
		}

		//log.Println("episodesFiltered:", episodesFiltered)
		if len(episodesFiltered.Company) > 0 {
			log.Println("there already exists:", bsonMap["company"])
			return
		} else {
			log.Println(bsonMap["company"], "will be inserted.")
		}
	}

	if table_name == ColnameUser {
		var result User
		readOne, _ := collection.Find(context.Background(), bson.D{{"email", bsonMap["email"]}})
		readOne.Decode(&result)
		if len(result.Email) > 0 {
			log.Println("This user is already registered.")
			return
		}
	}

	_, err = collection.InsertOne(ctx, bsonMap)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (db *DB) ReadMongo(user_iput ...string) []JsonJob {
	log.Println("ReadMongo: user input is ", user_iput)
	// get table(=collection)
	collection := db.Client.Database(Dbname).Collection(Colname)

	findOptions := options.Find()
	// Sort by `date` field descending
	findOptions.SetSort(bson.D{{"dateadded", -1}})

	cur, err := collection.Find(context.Background(), bson.D{{"dateadded", bson.D{{"$gt", "2020-11-01"}}}}, findOptions)
	if err != nil {
		return nil
	}

	if len(user_iput) > 0 {
		cur, err = collection.Find(context.Background(), bson.M{"company": user_iput[0]}, findOptions)
		if err != nil {
			log.Println("err from user input:", err)
			return nil
		}
	}

	var jobs []JsonJob
	var doc JsonJob
	for cur.Next(context.Background()) {
		//var doc JsonJob
		err := cur.Decode(&doc)
		if err != nil {
			fmt.Println("error at cur.Decode(&doc)")
			return nil
		}
		//append to jobs
		jobs = append(jobs, doc)
		//log.Println("searched company:", doc.Company)
	}
	return jobs
}
