package main

import (
	"Scraping/app/controllers"
	_ "Scraping/app/controllers"
	"log"
	"net/http"
	"os"
	"time"
	_ "time"

	_ "go.mongodb.org/mongo-driver/bson"
)

func main() {
	log.Println("os.Getenv:", os.Getenv("MONGO_SERVER"))

	server := http.Server{}
	if os.Getenv("MONGO_SERVER") == "" {
		server.Addr = ":8080"
	}

	http.HandleFunc("/", controllers.HomeHandler)
	http.HandleFunc("/index-login-success", controllers.HomeHandlerAfterLogin)
	http.HandleFunc("/signup", controllers.SignUpHandler)
	http.HandleFunc("/login", controllers.LoginHandler)
	http.HandleFunc("/internal", controllers.InternalHandler)
	http.HandleFunc("/userpost", controllers.PostHandler)
	http.HandleFunc("/about", controllers.AboutHandler)
	//add css below
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css")))) //http.Handle("/css/")

	go ticker()
	//controllers.Init() //initialize Redis
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func ticker() {
	log.Println("ticker() is called.")
	var url [2]string
	url[0] = "https://www.linkedin.com/jobs/search/?geoId=101174742&keywords=intern&location=Canada"
	url[1] = "https://www.glassdoor.ca/Job/canada-software-engineer-internship-jobs-SRCH_IL.0,6_IN3_KO7,35.htm"

	// 1. で定義したMongoDBクライアント作成関数から構造体を取得
	mongoClient, err := controllers.ConnectMongoDB() //mongoClient is a pointer of address to DB.
	if err != nil {
		log.Println("Error from ConnectMongoDB()!", err)
		os.Exit(1)
	}

	if os.Getenv("MONGO_SERVER") != "" {
		mongoClient.DoMongoImport()
	}

	for i := range url {
		mongoClient.GetURL(url[i])
	}
	//clean up DB
	err = mongoClient.DeleteDuplicate()
	if err != nil {
		log.Println("error from DeleteDuplicate()")
		log.Println(err)
	}
	t := time.NewTicker(10 * time.Hour)
	for {
		select {
		case <-t.C:
			log.Println("ticker is working.")

			if os.Getenv("MONGO_SERVER") != "" {
				// this is in kubernetes cluster
				//mongoClient.DoMongoImport()
			} else { //this is in localhost
				//web crawl　and store into mongo
				for i := range url {
					mongoClient.GetURL(url[i])
				}
			}
		} //end of one t transaction.
	} //end of for loop
	t.Stop()
}
