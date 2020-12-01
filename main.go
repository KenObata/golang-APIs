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

/*
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
*/

func main() {
	//ticker()
	log.Println("os.Getenv:", os.Getenv("MONGO_SERVER"))
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

	server := http.Server{} //if you use kubectl
	if os.Getenv("MONGO_SERVER") == "" {
		server = http.Server{Addr: "127.0.0.1:3000"}
	}

	//Addr: "127.0.0.1:8080",
	//}
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
