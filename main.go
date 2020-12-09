package main

import (
	"Scraping/app/controllers"
	_ "Scraping/app/controllers"
	"fmt"
	"log"
	"net/http"
	"os"
	_ "time"

	_ "go.mongodb.org/mongo-driver/bson"
)

func main() {
	log.Println("os.Getenv:", os.Getenv("MONGO_SERVER"))
	// 1. で定義したMongoDBクライアント作成関数から構造体を取得
	mongoClient, err := controllers.ConnectMongoDB() //mongoClient is a pointer of address to DB.
	fmt.Println("my mongoClient:", mongoClient)
	if err != nil {
		fmt.Println("Error from ConnectMongoDB()!")
		fmt.Println(err)
		os.Exit(1)
	}
	var url [2]string
	url[0] = "https://www.linkedin.com/jobs/search/?geoId=101174742&keywords=intern&location=Canada"
	url[1] = "https://www.glassdoor.ca/Job/canada-software-engineer-internship-jobs-SRCH_IL.0,6_IN3_KO7,35.htm"

	server := http.Server{} //if you use kubectl
	if os.Getenv("MONGO_SERVER") == "" {
		server.Addr = ":8080"
	}

	http.HandleFunc("/", controllers.HomeHandler)
	http.HandleFunc("/signup", controllers.SignUpHandler)
	http.HandleFunc("/login", controllers.LoginHandler)
	http.HandleFunc("/internal", controllers.InternalHandler)
	http.HandleFunc("/userpost", controllers.PostHandler) //debug
	//add css below
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css")))) //http.Handle("/css/")
	server.ListenAndServe()
	mongoClient.DoMongoImport()
	for i := range url {
		mongoClient.GetURL(url[i])
	}

	/*
		t := time.NewTicker(2 * time.Hour)
		for {
			select {
			case <-t.C:
				// every t hour, run get URL function.
				// web crawl　and store into mongo
				mongoClient.GetURL(url)
			}
		}
		t.Stop()*/

}
