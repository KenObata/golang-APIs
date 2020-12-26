package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func HomeHandlerAfterLogin(w http.ResponseWriter, r *http.Request) {
	mongoClient, err := ConnectMongoDB()
	if err != nil {
		fmt.Println("Error from ConnectMongoDB()!")
		fmt.Println(err)
		os.Exit(1)
	}
	// Working Directory
	wd, err := os.Getwd()
	t, err := template.ParseFiles(wd + "/app/view/index-login-success.html")
	if err != nil {
		log.Println("Error from template.ParseFiles()!")
		log.Println(err)
	}

	name := r.FormValue("name")

	var filter_condition [3]bool = [3]bool{false, false, false}
	checkSoftware := r.FormValue("filterSoftware")
	if checkSoftware == "true" {
		filter_condition[0] = true
	}
	checkDataScience := r.FormValue("filterDataScience")
	if checkDataScience == "true" {
		filter_condition[1] = true
	}
	checkThisWeek := r.FormValue("filterThisWeek")
	log.Println("checkThisWeek:", checkThisWeek)
	if checkThisWeek == "true" {
		filter_condition[2] = true
	}
	//pass filter condition
	new_job_struct := mongoClient.ReadMongo(name, filter_condition[0], filter_condition[1], filter_condition[2])
	t.Execute(w, new_job_struct)
	time.Sleep(1000)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
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
	var filter_condition [3]bool = [3]bool{false, false, false}

	log.Println("filter_condition", filter_condition)
	name := r.FormValue("name")
	log.Println("Company:", name, "is input.")

	checkSoftware := r.FormValue("filterSoftware")
	if checkSoftware == "true" {
		filter_condition[0] = true
	}
	checkDataScience := r.FormValue("filterDataScience")
	if checkDataScience == "true" {
		filter_condition[1] = true
	}
	checkThisWeek := r.FormValue("filterThisWeek")
	if checkThisWeek == "true" {
		filter_condition[2] = true
	}
	log.Println("checkThisWeek:", filter_condition)
	//pass filter condition
	new_job_struct := mongoClient.ReadMongo(name, filter_condition[0], filter_condition[1], filter_condition[2])
	t.Execute(w, new_job_struct)
	time.Sleep(1000)
}

// 実際にMongoDBへ接続するクライアントを内包したDB addressを返却
func ConnectMongoDB() (*DB, error) {
	log.Println("ConnectMongoDB() is called.")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 認証が必要な場合は、options.Credentialを作成
	credential := options.Credential{
		AuthSource: Dbname,
		Username:   MongoUser,
		Password:   MongoPassword,
	}
	log.Println("credential.Username:", credential.Username)
	// 認証情報・接続情報を元にclientを作成
	var host string
	var port string
	if os.Getenv("MONGO_SERVER") == "" {
		host = MongoDBHost
		port = "27016"
	} else { // running on kubernetes
		host = os.Getenv("MONGO_SERVER")
		port = MongoDBPort
	}
	log.Println("host:", host, "port:", port)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+host+":"+port).SetAuth(credential))
	if err != nil {
		fmt.Println("error from mongo.Connect(ctx,")
		fmt.Println(err)
		return nil, err
	}
	return &DB{client}, nil
}

func GetGlassdoor(doc *goquery.Document) (urls []string, companies []string, titles []string) {
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		class, _ := s.Attr("class")
		url, _ := s.Attr("href")
		url = "https://www.glassdoor.ca" + url
		title := s.Text()

		if strings.Contains(class, "jobLink") && !(strings.Contains(s.Text(), ";")) && s.Text() != "" && len(companies)-len(titles) == 0 {
			company := s.Text()
			companies = append(companies, company)
			//log.Println("Inside of doc.Find(), company name:", company)
		} else if strings.Contains(class, "jobTitle") && len(companies)-len(titles) == 1 {
			title = s.Text()
			titles = append(titles, title)
			urls = append(urls, url)
			//log.Println("Inside of doc.Find(), titles name:", title)
		}
	})
	if len(companies) > len(titles) {
		companies = companies[:len(companies)-1]
	}
	return urls, companies, titles
}

func (mongoClient *DB) GetURL(URL string) error {
	log.Println("GetURL function is called from main.")
	log.Println("URL:", URL)
	//set timeout 5 sec
	_, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Load the URL
	res, e := http.Get(URL)
	if e != nil {
		log.Println("Error from http.Get(URL)")
		return e
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println("Error from goquery.NewDocumentFromReader.")
		//log.Fatal(err)
		return err
	}
	var urls []string
	var companies []string
	var titles []string

	if strings.Contains(URL, "glassdoor") {
		urls, companies, titles = GetGlassdoor(doc)
	} else {
		doc.Find("a").Each(func(_ int, s *goquery.Selection) {
			url, _ := s.Attr("href")
			title := s.Text()

			if len(companies) < len(titles)-1 {
				log.Println("something is wrong with:", titles)
				os.Exit(1)
			}
			if strings.Contains(url, "/company/") {
				//get company name
				company := s.Text()
				//log.Println("Get URL(), company name:", company)
				companies = append(companies, company)
			} else if strings.Contains(url, "https://ca.linkedin.com/jobs/view/") && len(companies) == len(titles) {
				urls = append(urls, url)
				titles = append(titles, title)
			}
		})
	}

	job := &Job{
		Title:   titles,
		URL:     urls,
		Company: companies,
	}
	// Unmarshal結果の格納先である構造体のポインターを取得
	jsonJob := new(JsonJob)
	//create json
	var i int
	currentTime := time.Now()
	log.Println("len(job.Company):", len(job.Company))
	log.Println("len(job.Title):", len(job.Title))
	log.Println("len(job.URL):", len(job.URL))
	for i = 0; i < len(companies); i++ {
		jsonJob.URL = job.URL[i]
		jsonJob.Title = job.Title[i]
		jsonJob.Company = job.Company[i]
		jsonJob.DateAdded = currentTime.Format("2006-01-02")

		// convert Job struct to JSON by marshal
		jsonJobJSON, err := json.Marshal(jsonJob)
		if err != nil {
			log.Println("error from json.Marshal(jsonJob)", jsonJobJSON)
			log.Println(err)
			return err
		}
		// Insert JSON data to MongoDB
		err = mongoClient.InsertMongoDB(jsonJobJSON, Colname)
		if err != nil {
			//log.Fatal(err)
			return err
		}
	} //end of for loop of each array
	log.Println("End of for loop to insert jsonJobJSON.")
	return nil
}

/*
func (db *DB) DoMongoImport() {
	log.Println("DoMongoImport called.")
	wd, _ := os.Getwd()
	docsPath, _ := filepath.Abs(wd + "/app/Job.json")
	byteValues, err1 := ioutil.ReadFile(docsPath)
	if err1 != nil {
		fmt.Println(err1.Error())
		os.Exit(1)
	}
	var mongoExport []MongoExportJobs
	err := json.Unmarshal(byteValues, &mongoExport)
	if err != nil {
		log.Println("error from unmarshal!")
		log.Println(err)
	}
	for _, data := range mongoExport {
		log.Println(data.Company)
		jsonJobJSON, _ := json.Marshal(data)
		db.InsertMongoDB(jsonJobJSON, Colname)
	}
}
*/
