package main

import (
	"Scraping/app/controllers"
	_ "Scraping/app/controllers"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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
	url := "https://www.linkedin.com/jobs/search/?geoId=101174742&keywords=intern&location=Canada"
	mongoClient.GetURL(url)
	t := time.NewTicker(2 * time.Hour)
	for {
		select {
		case <-t.C:
			// every t hour, run get URL function.
			// web crawl　and store into mongo
			mongoClient.GetURL(url)
		}
	}
	t.Stop()

	server := http.Server{} //if you use kubectl
	if os.Getenv("MONGO_SERVER") == "" {
		server.Addr = "127.0.0.1:3000"
	}
	//for docker-compose
	/*
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
	*/
	http.HandleFunc("/", controllers.HomeHandler)
	http.HandleFunc("/signup", controllers.SignUpHandler)
	http.HandleFunc("/login", controllers.LoginHandler)
	http.HandleFunc("/internal", controllers.InternalHandler)
	http.HandleFunc("/userpost", controllers.PostHandler) //debug
	//add css below
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css")))) //http.Handle("/css/")
	server.ListenAndServe()

	//for docker-compose
	/*
		log.Println("** Service Started on Port " + port + " **")
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}*/
}
