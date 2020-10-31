package controllers

/*
import (
	"Scraping/data"
	_ "Scraping/data"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//noinspection ALL
func homeHandler(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w,r,"index.html")
	mongoClient, err := ConnectMongoDB()
	if err != nil {
		fmt.Println("Error from ConnectMongoDB()!")
		fmt.Println(err)
		os.Exit(1)
	}
	// Working Directory
	wd, err := os.Getwd()
	t, err := template.ParseFiles(wd + "/app/view/index.html")
	if err != nil {
		log.Println("Error from template.ParseFiles()!")
		log.Println(err)
	}

	name := r.FormValue("name")
	log.Println("Company:", name, "is input.")
	if name != "" {
		new_job_struct := mongoClient.readMongo(name)
		t.Execute(w, new_job_struct)
		time.Sleep(1000) //debug
	} else {
		job_struct := mongoClient.readMongo()
		t.Execute(w, job_struct)
	}
}

const (
	// 接続先のDB情報を入力
	mongoDBHost   = "127.0.0.1"
	mongoDBPort   = "27017"
	mongoUser     = "Ken"
	mongoPassword = "k0668466425"
	dbname        = "test" //"databases"
	colname       = "Job"
)

// 実際にMongoDBへ接続するクライアントを内包したDB addressを返却
//noinspection ALL
func ConnectMongoDB() (*data.DB, error) {
	ctx := context.Background()
	// 認証が必要な場合は、options.Credentialを作成
	credential := options.Credential{
		AuthSource: dbname,
		Username:   mongoUser,
		Password:   mongoPassword,
	}
	// 認証情報・接続情報を元にclientを作成
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+mongoDBHost+":"+mongoDBPort).SetAuth(credential))
	if err != nil {
		fmt.Println("error from mongo.Connect(ctx,")
		fmt.Println(err)
		return nil, err
	}
	//return &DB{client}, nil
	return &data.DB{client}, nil
}

func (db *data.DB) readMongo(user_iput ...string) []data.JsonJob {
	log.Println("readMongo: user input is ", user_iput)
	// get table(=collection)
	collection := db.client.Database(dbname).Collection(colname)

	findOptions := options.Find()
	// Sort by `price` field descending
	findOptions.SetSort(bson.D{{"dateadded", -1}})

	cur, err := collection.Find(context.Background(), bson.D{}, findOptions)
	if err != nil {
		return nil
	}

	//log.Println("user_iput[0]", user_iput[0])
	if len(user_iput) > 0 {
		cur, err = collection.Find(context.Background(), bson.M{"company": user_iput[0]}, findOptions)
		if err != nil {
			log.Println("err from user input:", err)
			return nil
		}
	}

	var jobs []data.JsonJob
	var doc data.JsonJob //こっちに移動した
	for cur.Next(context.Background()) {
		//var doc JsonJob
		err := cur.Decode(&doc)
		if err != nil {
			fmt.Println("error at cur.Decode(&doc)")
			return nil
		}
		//append to jobs
		jobs = append(jobs, doc)
		log.Println("searched company:", doc.Company)
	}
	return jobs
}


*/
