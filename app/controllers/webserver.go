package controllers

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func HomeHandlerAfterLogin(w http.ResponseWriter, r *http.Request) {
	// Working Directory
	wd, err := os.Getwd()
	t, err := template.ParseFiles(wd + "/app/view/index-login-success.html")
	if err != nil {
		log.Println("Error from template.ParseFiles()!", err)
	}

	name := r.FormValue("name")

	var filter_condition [4]bool = [4]bool{false, false, false, false}
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
	done := r.FormValue("filterDone")
	if done == "true" {
		filter_condition[3] = true
	}
	//pass SessionID to cookie
	if SessionID != "" {
		expiration := time.Now().Add(7 * 24 * time.Hour) //1 week
		cookie := &http.Cookie{
			Name:    "SessionID",
			Value:   SessionID,
			Path:    "/",
			Expires: expiration,
		}
		http.SetCookie(w, cookie)
	}

	//get user ID if sessionID exists
	cookie, err := r.Cookie("SessionID")
	var sessionID string = ""
	var userId = "-1"
	if cookie != nil {
		sessionID = cookie.Value
		//get userId
		userId, err = GetKey(context.Background(), sessionID)
		if err != nil {
			log.Println("error from GetKey(): ", err)
		}
		log.Println("userId: ", userId)
	}

	var new_job_struct []JsonJob
	//get all checkboxes and insert to user_job table if checked.
	r.ParseForm()
	jobIdSelected := r.Form["appliedCheckbox"]
	userIdint, _ := strconv.Atoi(userId)
	if userId != "-1" {
		for _, jobId := range jobIdSelected {
			//insert postgres job_user table

			jobIdint, _ := strconv.Atoi(jobId)
			log.Println("jobId: ", jobIdint, " is checked.")
			err = InsertUserJob(userIdint, jobIdint)
			if err != nil {
				log.Println("InsertUserJob failed: ", err)
			}
		}
	}
	new_job_struct = ReadUserCustomPostgres(name, filter_condition[0], filter_condition[1], filter_condition[2], filter_condition[3], userIdint)
	//new_job_struct = ReadPostgres(name, filter_condition[0], filter_condition[1], filter_condition[2])
	log.Println("len(new_job_struct): ", len(new_job_struct))
	//pass filter condition
	t.Execute(w, new_job_struct)
	time.Sleep(1000)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	//get SessionID if exists in cookie
	cookie, err := r.Cookie("SessionID")
	var sessionID string = ""
	if cookie != nil {
		sessionID = cookie.Value
		//httpRedirect if SessionID exists.
		target := "http://" + r.Host + "/index-login-success"
		log.Println("http redirect to ", target)
		http.Redirect(w, r, target, http.StatusFound)
	}
	log.Println("local sessionID:", sessionID)

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
	new_job_struct := ReadPostgres(name, filter_condition[0], filter_condition[1], filter_condition[2])
	t.Execute(w, new_job_struct)
	time.Sleep(1000)
}

func ConnectMongoDB() (*DB, error) {
	log.Println("ConnectMongoDB() is called.")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	credential := options.Credential{
		AuthSource: Dbname,
		Username:   MongoUser,
		Password:   MongoPassword,
	}
	log.Println("credential.Username:", credential.Username)

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
