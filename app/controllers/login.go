package controllers

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

func InternalHandler(w http.ResponseWriter, r *http.Request) {
	//get user input first
	email := r.FormValue("email")
	password := r.FormValue("password")
	//debug
	log.Println("input email in InternalHandler:", email)
	log.Println("input PW in InternalHandler:", password)

	db := ConnectPostgres()
	//check if user certainly exists
	rows, err := db.Query("SELECT id, email, password FROM user_list where email=$1 AND password=$2;", email, password)
	if err != nil {
		log.Println("error from SELECT * from user_list", err.Error())
	}
	if !(rows.Next()) {
		log.Println("Login failed.")
		wd, err := os.Getwd()
		t, err := template.ParseFiles(wd + "/app/view/login-error.html")
		if err != nil {
			log.Println(err)
		}
		t.Execute(w, nil)
	} else {
		var id string
		var email string
		var password string
		err := rows.Scan(&id, &email, &password)
		if err != nil {
			log.Println("error from rows.Scan() ", err)
		}

		//create a session by Redis
		SessionID, err = SetKey(context.Background(), id) //userFiltered[0] = User struct
		if err != nil {
			log.Println(err)
		} else {
			log.Println("SessionID: ", SessionID)
			//pass uuid to cookie -> do it in th HomePageAfterHandler
		}

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

	//get user iput first
	var job JsonJob
	companyName := r.FormValue("companyName")
	jobTitle := r.FormValue("jobTitle")
	jobURL := r.FormValue("jobURL")

	//get dateadded column
	currentTime := time.Now()

	if companyName != "" && jobTitle != "" && jobURL != "" {
		job.Company = companyName
		job.Title = jobTitle
		job.URL = jobURL
		job.DateAdded = currentTime.Format("2006-01-02")

		err = InsertJob(job)

		if err != nil {
			log.Println("Error from InsertJob:", err)
			new_t, _ := template.ParseFiles(wd + "/app/view/userpost-error.html")
			new_t.Execute(w, err.Error())
		} else {
			new_t, _ := template.ParseFiles(wd + "/app/view/userpost-success.html")
			new_t.Execute(w, nil)
		}

	} //end of insert

}
