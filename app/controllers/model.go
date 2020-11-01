package controllers

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// DB構造体へInsert用のメソッドを定義
// JSONファイルから読み込んだバイトスライスを渡し、MongoDBへInsert
func (db *DB) InsertMongoDB(json []byte, table_name string) {
	// 3秒でタイムアウトするコンテキストを作成
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	bsonMap := bson.M{}
	// JSONのバイトスライスをMongoDBのドキュメント型であるbsonへマップ
	err := bson.UnmarshalExtJSON([]byte(json), false, &bsonMap)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Insert先のコレクション名からクライアント作成
	collection := db.Client.Database(Dbname).Collection(table_name) //colname is the parameter for table_name
	fmt.Println("bsonMap:", bsonMap)
	//fmt.Println("ctx:", ctx)

	if table_name == Colname {
		readOne, _ := collection.Find(context.Background(), bson.D{{"url", bsonMap["url"]}})
		if readOne != nil {
			fmt.Println("there already exists:", bsonMap["company"])
			return
		}
	}

	if table_name == ColnameUser {
		readOne, _ := collection.Find(context.Background(), bson.D{{"email", bsonMap["email"]}})
		if readOne != nil {
			//This user is already registered.
			return
		}
	}

	_, err = collection.InsertOne(ctx, bsonMap)
	if err != nil {
		fmt.Println(err)
		return
	}
}
