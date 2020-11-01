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
