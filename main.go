package main

import (
	"Scraping/app/controllers"
	_ "Scraping/app/controllers"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "time"

	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
)

/*
type Job struct {
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

// mongo-driverのクライアントを自前で定義した構造体DBへセット
type DB struct {
	client *mongo.Client
}
*/

//this month use for mongo find filtering
//var thisYear string = time.RFC3339[0:8]

/*
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
		//log.Println("job_struct[0]:", job_struct[0])
		t.Execute(w, job_struct)
	}
}
*/
func errorInResponse(w http.ResponseWriter, status int, error controllers.Error) {
	w.WriteHeader(status) // HTTP status code such as 400, 500
	json.NewEncoder(w).Encode(error)
	return
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	var user controllers.User
	//var error Error

	// Working Directory
	wd, err := os.Getwd()
	t, err := template.ParseFiles(wd + "/app/view/signup.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)

	//we ganna insert into User collection (use later)
	mongoClient, _ := controllers.ConnectMongoDB()

	//get ID by number of users + 1
	collection := mongoClient.Client.Database(controllers.Dbname).Collection(controllers.ColnameUser)
	cur, err := collection.Find(context.Background(), bson.D{})
	numOfUsers := 0
	for cur.Next(context.Background()) {
		numOfUsers += 1
	}
	user.ID = numOfUsers + 1
	//get email from html file
	email := r.FormValue("email")
	if email != "" {
		log.Println("email:", email)
		user.Email = email
	}
	password := r.FormValue("password")
	if password != "" && len(password) > 3 {
		user.Password = password
	}
	json.NewDecoder(r.Body).Decode(&user)
	/*
		if user.Email == "" {
			errorInResponse(w, http.StatusBadRequest, error)
			return
		}
		if user.Password == "" {
			errorInResponse(w, http.StatusBadRequest, error)
			return
		}

	*/

	//hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//user.Password = string(hash)

	userJSON, err := json.Marshal(user)
	log.Println("email:", user.Email, "password:", user.Password)
	if err != nil {
		return
	}
	mongoClient.InsertMongoDB(userJSON, controllers.ColnameUser)
}

func ticker() {
	//t := time.NewTicker(24 * time.Hour) //24 Hour周期の ticker
	//defer t.Stop()

	url := "https://www.linkedin.com/jobs/search/?geoId=101174742&keywords=intern&location=Canada"

	// 1. で定義したMongoDBクライアント作成関数から構造体を取得
	mongoClient, err := controllers.ConnectMongoDB() //mongoClient is a pointer of address to DB.
	fmt.Println("my mongoClient:", mongoClient)
	if err != nil {
		fmt.Println("Error from ConnectMongoDB()!")
		fmt.Println(err)
		os.Exit(1)
	}
	// web crawl　and store into mongo
	mongoClient.GetURL(url)
}

func main() {
	ticker()

	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/", controllers.HomeHandler)
	http.HandleFunc("/signup", signUpHandler)
	//add css below
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css")))) //http.Handle("/css/")
	server.ListenAndServe()
}

/*
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
		mongoClient.(controllers.InsertMongoDB(jsonJobJSON, colname))
	} //end of for loop of each array
}
*/
/*
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
*/
/*
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
	collection := db.client.Database(dbname).Collection(table_name) //colname is the parameter for table_name
	fmt.Println("bsonMap:", bsonMap)
	//fmt.Println("ctx:", ctx)

	if table_name == colname {
		readOne, _ := collection.Find(context.Background(), bson.D{{"url", bsonMap["url"]}})
		if readOne != nil {
			//fmt.Println("there already exists:", bsonMap["company"])
			return
		}
	}

	if table_name == colnameUser {
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
*/

/*
func (db *DB) readMongo(user_iput ...string) []JsonJob {
	//log.Println("readMongo: user input is ", user_iput)
	// get table(=collection)
	collection := db.client.Database(dbname).Collection(colname)

	findOptions := options.Find()
	// Sort by `price` field descending
	findOptions.SetSort(bson.D{{"dateadded", -1}})

	cur, err := collection.Find(context.Background(), bson.D{}, findOptions)
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
*/
