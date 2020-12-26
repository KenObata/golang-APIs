package main

import (
	"Scraping/app/controllers"
	_ "Scraping/app/controllers"
	"log"
	"net/http"
	"os"
	"time"
	_ "time"

	"github.com/sony/gobreaker"
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
		log.Fatal("error from ListenAndServe()", err)
	}
}

var cb *gobreaker.CircuitBreaker

func init() {
	var settings gobreaker.Settings
	settings.Name = "HTTP GET"
	settings.ReadyToTrip = func(counts gobreaker.Counts) bool {
		//circuit breaker will trip when 60% of requests failed an dat least 10 requests were made.
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 5 && failureRatio >= 0.6
	}
	settings.Timeout = time.Millisecond
	settings.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
		if to == gobreaker.StateOpen {
			log.Println("State Open!")
		}
		if from == gobreaker.StateOpen && to == gobreaker.StateHalfOpen {
			log.Println("Going from Open to Half-Open!")
		}
		if from == gobreaker.StateHalfOpen && to == gobreaker.StateClosed {
			log.Println("Going from Half-Open to Close!")
		}
	} //end of settings.OnStateChange
	cb = gobreaker.NewCircuitBreaker(settings)
}

func ticker() {
	log.Println("ticker() is called.")
	var url [2]string
	url[0] = "https://www.linkedin.com/jobs/search/?geoId=101174742&keywords=intern&location=Canada"
	url[1] = "https://www.glassdoor.ca/Job/canada-software-engineer-internship-jobs-SRCH_IL.0,6_IN3_KO7,35.htm"

	mongoClient, err := controllers.ConnectMongoDB() //mongoClient is a pointer of address to DB.
	if err != nil {
		log.Println("Error from ConnectMongoDB()!", err)
		os.Exit(1)
	}

	//wrap by Circuit Breaker
	_, err = cb.Execute(func() (interface{}, error) {
		for i := range url {
			err := mongoClient.GetURL(url[i])
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	}) //end of Circuit Breaker wrapper
	if err != nil {
		log.Println("error from mongoClient.GetURL(): ", err)
	}

	//clean up DB once
	log.Println("Run DeleteDuplicate().")
	deleteError := mongoClient.DeleteDuplicate()
	if deleteError != nil {
		log.Println("erro from DeleteDuplicate(): ", deleteError)
	}

	log.Println("start ticker()")
	t := time.NewTicker(10 * time.Hour)
	for {
		select {
		case <-t.C:
			log.Println("ticker is working.")
			if os.Getenv("MONGO_SERVER") == "" { //this is in localhost
				//web crawl　and store into mongo
				for i := range url {
					mongoClient.GetURL(url[i])
				}
			} else { //on kubernetes cluster
				//clean up DB
				deleteError := mongoClient.DeleteDuplicate()
				if deleteError != nil {
					log.Println("error from DeleteDuplicate()")
					log.Println(deleteError)
				}
			}
		} //end of one t transaction, select.
	} //end of for loop
	t.Stop()
}
