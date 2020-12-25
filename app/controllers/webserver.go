package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	if name != "" {
		new_job_struct := mongoClient.ReadMongo(name)
		t.Execute(w, new_job_struct)
		time.Sleep(1000)
	} else {
		job_struct := mongoClient.ReadMongo()
		t.Execute(w, job_struct)
	}
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

	name := r.FormValue("name")
	log.Println("Company:", name, "is input.")
	if name != "" {
		new_job_struct := mongoClient.ReadMongo(name)
		t.Execute(w, new_job_struct)
		time.Sleep(1000) //debug
	} else {
		job_struct := mongoClient.ReadMongo()
		t.Execute(w, job_struct)
	}
}

// 実際にMongoDBへ接続するクライアントを内包したDB addressを返却
func ConnectMongoDB() (*DB, error) {
	log.Println("ConnectMongoDB() is called.")
	ctx := context.Background()
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
	log.Println("host:", host)
	//debug
	//clientOpts := options.Client().ApplyURI("mongodb://" + host + ":" + MongoDBPort + "/?connect=direct")
	//client, err := mongo.Connect(ctx, clientOpts)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+host+":"+port).SetAuth(credential))
	if err != nil {
		fmt.Println("error from mongo.Connect(ctx,")
		fmt.Println(err)
		return nil, err
	}
	return &DB{client}, nil
}
func GetVIATEC(doc *goquery.Document) (urls []string, companies []string, titles []string) {
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		//log.Println(url)
		url = "https://www.viatec.ca/" + url
		title := s.Text()
		if strings.Contains(url, "/jobs/") && len(companies)-len(titles) == 0 {
			company := s.Next().Text()
			companies = append(companies, company)
			log.Println("Inside of doc.Find(), company name:", company)
		} else if strings.Contains(url, "jobs") && len(companies)-len(titles) == 1 {
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

func (mongoClient *DB) GetURL(URL string) {
	log.Println("GetURL function is called from main.")
	log.Println("URL:", URL)

	// Load the URL
	res, e := http.Get(URL)
	if e != nil {
		log.Println("Error from http.Get(URL)")
		log.Fatal(e)
		return
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println("Error from goquery.NewDocumentFromReader.")
		log.Fatal(err)
		return
	}
	var urls []string
	var companies []string
	var titles []string

	if strings.Contains(URL, "glassdoor") {
		urls, companies, titles = GetGlassdoor(doc)
	} else if strings.Contains(URL, "viatec") {
		urls, companies, titles = GetVIATEC(doc)
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
				//log.Println("Get URL(), titles name:", title)
				//log.Println("Get URL(), url name:", url)
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

		// 構造体をJSON文字列に変換
		jsonJobJSON, err := json.Marshal(jsonJob)
		if err != nil {
			log.Println("error from json.Marshal(jsonJob)", jsonJobJSON)
			log.Println(err)
			return
		}
		// Insert JSON data to MongoDB
		err = mongoClient.InsertMongoDB(jsonJobJSON, Colname)
		if err != nil {
			log.Fatal(err)
		}
	} //end of for loop of each array
	log.Println("End of for loop to insert jsonJobJSON.")
}

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
