package controllers

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func InternalHandler(w http.ResponseWriter, r *http.Request) {
	//get user input first
	email := r.FormValue("email")
	password := r.FormValue("password")
	//debug
	log.Println("input email in InternalHandler:", email)
	log.Println("input PW in InternalHandler:", password)
	mongoClient, _ := ConnectMongoDB()
	collection := mongoClient.Client.Database(Dbname).Collection(ColnameUser)
	cur, err := collection.Find(context.Background(), bson.M{"email": email, "password": password})
	var episodesFiltered []bson.M
	if err = cur.All(context.Background(), &episodesFiltered); err != nil {
		log.Fatal(err)
	}
	log.Println("cur.Current:", cur.Current)
	if len(episodesFiltered) == 0 {
		//var error Error
		//errorInResponse(w, http.StatusBadRequest, error)
		log.Println("Login failed.")
		wd, err := os.Getwd()
		t, err := template.ParseFiles(wd + "/app/view/login-error.html")
		if err != nil {
			log.Println(err)
		}
		t.Execute(w, nil)
	} else {
		//create a session by Redis
		SetKey(episodesFiltered[0])
		log.Println("episodesFiltered[0]:", episodesFiltered[0])
		log.Println("episodesFiltered[0][\"id\"]:", episodesFiltered[0]["id"])

		//http Redirect
		target := "http://" + r.Host + "/userpost"
		log.Println("http redirect to ", target)
		http.Redirect(w, r, target, http.StatusFound)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	log.Println("LoginHandler called. wd := ", wd)
	t, err := template.ParseFiles(wd + "/app/view/login.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// Working Directory
	wd, err := os.Getwd()
	t, err := template.ParseFiles(wd + "/app/view/userpost.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)

	//get user niput first
	var job JsonJob
	companyName := r.FormValue("companyName")
	jobTitle := r.FormValue("jobTitle")
	jobURL := r.FormValue("jobURL")

	//append to JSON object
	if companyName != "" {
		job.Company = companyName
	}
	if jobTitle != "" {
		job.Title = jobTitle
	}
	if jobURL != "" {
		job.URL = jobURL
	}
	//get dateadded column
	currentTime := time.Now()
	job.DateAdded = currentTime.Format("2006-01-02")

	if companyName != "" && jobTitle != "" && jobURL != "" {
		jsonJobJSON, err := json.Marshal(job)

		mongoClient, _ := ConnectMongoDB()
		err = mongoClient.InsertMongoDB(jsonJobJSON, Colname)
		if err != nil {
			log.Println("User InsertMongoDB:", err)
			new_t, _ := template.ParseFiles(wd + "/app/view/userpost-error.html")
			new_t.Execute(w, err)
		} else {
			new_t, _ := template.ParseFiles(wd + "/app/view/userpost-success.html")
			new_t.Execute(w, nil)
		}

	}

}
