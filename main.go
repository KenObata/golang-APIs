package main

import (
	_ "Scraping/app/controllers"
	"html/template"

	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	_ "time"

	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Job struct {
	//ID string
	URL     []string
	Title   []string
	Company []string
}
type JsonJob struct {
	URL       string `json:"url"`
	Title     string `json:"title"`
	Company   string `json:"company"`
	DateAdded string `json:"dateadded"`
}

// Array for Job struct
type Jobs []*Job

//this month use for mongo find filtering
//var thisYear string = time.RFC3339[0:8]

const (
	// 接続先のDB情報を入力
	mongoDBHost   = "127.0.0.1"
	mongoDBPort   = "27017"
	mongoUser     = "Ken"
	mongoPassword = "k0668466425"
	dbname        = "test" //"databases"
	colname       = "Job"
)

// mongo-driverのクライアントを自前で定義した構造体DBへセット
type DB struct {
	client *mongo.Client
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w,r,"index.html")
	mongoClient, err := ConnectMongoDB()
	if err != nil {
		fmt.Println("Error from ConnectMongoDB()!")
		fmt.Println(err)
		os.Exit(1)
	}
	job_struct := mongoClient.readMongo()
	log.Println("job_struct[0]:", job_struct[0])
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, job_struct)
}

func main() {
	url := "https://www.linkedin.com/jobs/search/?geoId=101174742&keywords=intern&location=Canada"

	// 1. で定義したMongoDBクライアント作成関数から構造体を取得
	mongoClient, err := ConnectMongoDB() //mongoClient is a pointer of address to DB.
	fmt.Println("my mongoClient:", mongoClient)
	if err != nil {
		fmt.Println("Error from ConnectMongoDB()!")
		fmt.Println(err)
		os.Exit(1)
	}

	// web crawl　and store into mongo
	mongoClient.getURL(url)
	//mongoClient.removeDuplicatesMongo()

	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/", homeHandler)
	//add css below
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	server.ListenAndServe()
}

/*
func (db *DB) removeDuplicatesMongo(){
	// get table(=collection)
	collection := db.client.Database(dbname).Collection(colname)
	//remove duplicated URL

}
*/

func (db *DB) readMongo() []JsonJob {
	log.Println("readMongo is called!")
	// get table(=collection)
	collection := db.client.Database(dbname).Collection(colname)
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		// TODO: Do something about the error
		fmt.Println("error at collection.Find()")
		return nil
	}

	var jobs []JsonJob
	var doc JsonJob //こっちに移動した
	for cur.Next(context.Background()) {
		//var doc JsonJob
		err := cur.Decode(&doc)
		if err != nil {
			fmt.Println("error at cur.Decode(&doc)")
			return nil
		}
		//append to jobs
		jobs = append(jobs, doc)
	}
	return jobs
}

func (mongoClient *DB) getURL(URL string) {
	doc, err := goquery.NewDocument(URL)
	if err != nil {
		return
	}

	var urls []string
	var companies []string
	var titles []string
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		title := s.Text()

		if strings.Contains(url, "https://ca.linkedin.com/jobs/view/") {
			urls = append(urls, url)
			titles = append(titles, title)

		} else if strings.Contains(url, "/company/") {
			//get company name
			company := s.Text()
			companies = append(companies, company)
		}
	})
	job := &Job{
		Title:   titles,
		URL:     urls,
		Company: companies,
	}
	//fmt.Println("all titles:",job.Title)
	//fmt.Println("all comapnies:",job.Company)
	//fmt.Println("all urls:",job.URL)

	// Unmarshal結果の格納先である構造体のポインターを取得
	jsonJob := new(JsonJob)
	//create json
	var i int
	currentTime := time.Now()
	for i = 0; i < len(job.Company); i++ {
		jsonJob.URL = job.URL[i]
		jsonJob.Title = job.Title[i]
		jsonJob.Company = job.Company[i]
		jsonJob.DateAdded = currentTime.Format("2006-01-02")

		// 構造体をJSON文字列に変換
		jsonJobJSON, err := json.Marshal(jsonJob)
		if err != nil {
			fmt.Println("error from json.Marshal(jsonJob)")
			fmt.Println(err)
			return
		}

		// Insert JSON data to MongoDB
		mongoClient.InsertMongoDB(jsonJobJSON)
	} //end of for loop of each array
}

// 実際にMongoDBへ接続するクライアントを内包したDB addressを返却
func ConnectMongoDB() (*DB, error) {
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
	return &DB{client}, nil
}

// DB構造体へInsert用のメソッドを定義
// JSONファイルから読み込んだバイトスライスを渡し、MongoDBへInsert
func (db *DB) InsertMongoDB(json []byte) {
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
	collection := db.client.Database(dbname).Collection(colname)
	fmt.Println("bsonMap:", bsonMap)
	//fmt.Println("ctx:", ctx)
	readOne, err := collection.Find(context.Background(), bson.D{{"url", bsonMap["url"]}})
	if readOne != nil {
		fmt.Println("there already exists:", bsonMap["company"])
		return
	}
	//else
	_, err = collection.InsertOne(ctx, bsonMap)
	if err != nil {
		fmt.Println(err)
		return
	}
}
