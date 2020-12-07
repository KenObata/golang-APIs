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
	ctx := context.Background()
	// 認証が必要な場合は、options.Credentialを作成
	credential := options.Credential{
		AuthSource: Dbname,
		Username:   MongoUser,
		Password:   MongoPassword,
	}
	// 認証情報・接続情報を元にclientを作成
	var host string
	if os.Getenv("MONGO_SERVER") == "" {
		host = MongoDBHost
	} else {
		host = os.Getenv("MONGO_SERVER")
	}
	log.Println("host:", host)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+host+":"+MongoDBPort).SetAuth(credential))
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

	//debugging
	//testResponse, _ := http.Get(URL)
	//res_byte, _ := ioutil.ReadAll(testResponse.Body)
	//log.Println("res_byte:", string(res_byte))

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
	//log.Println("dummy.text:", doc.Text())
	var urls []string
	var companies []string
	var titles []string

	if strings.Contains(URL, "glassdoor") {
		urls, companies, titles = GetGlassdoor(doc)
	} else if strings.Contains(URL, "viatec") {
		urls, companies, titles = GetVIATEC(doc)
	} else {
		doc.Find("a").Each(func(_ int, s *goquery.Selection) {
			//log.Println(s.Text())
			url, _ := s.Attr("href")
			title := s.Text()

			if strings.Contains(url, "https://ca.linkedin.com/jobs/view/") {
				urls = append(urls, url)
				titles = append(titles, title)
				//log.Println("Inside of doc.Find(), titles name:", titles)
			} else if strings.Contains(url, "/company/") {
				//get company name
				company := s.Text()
				//log.Println("Inside of doc.Find(), company name:", company)
				companies = append(companies, company)
				//log.Println("companies:", companies)
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
			log.Println("error from json.Marshal(jsonJob)")
			log.Println(err)
			return
		}
		// Insert JSON data to MongoDB
		mongoClient.InsertMongoDB(jsonJobJSON, Colname)
	} //end of for loop of each array
	log.Println("End of for loop to insert jsonJobJSON.")
}

func (mongoClient *DB) DoMongoImport() {
	log.Println("DoMongoImport called.")
	wd, _ := os.Getwd()
	//t, err := template.ParseFiles(wd + "/app/view/index.html")
	docsPath, _ := filepath.Abs(wd + "/app/events.json")
	log.Println("docsPath:", docsPath)
	byteValues, _ := ioutil.ReadFile(docsPath)
	log.Println("byteValues:", string(byteValues))
	var user User
	err := json.Unmarshal(byteValues, &user)
	if err != nil {
		log.Println("error from unmarshal!")
		log.Println(err)
	}
	log.Println("user", user.Email)
}
