package controllers

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func errorInResponse(w http.ResponseWriter, status int, error Error) {
	w.WriteHeader(status) // HTTP status code such as 400, 500
	json.NewEncoder(w).Encode(error)
	return
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	// Working Directory
	wd, err := os.Getwd()
	t, err := template.ParseFiles(wd + "/app/view/signup.html")
	//path := os.Getenv(myPath) //create ENV in Dockerfile when we deploy to kubernetes.
	//t, err := template.ParseFiles("/app/view/signup.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)

	//we ganna insert into User collection (use later)
	var error Error
	mongoClient, _ := ConnectMongoDB()

	//get ID by number of users + 1
	collection := mongoClient.Client.Database(Dbname).Collection(ColnameUser)
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
	} else {
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}
	password := r.FormValue("password")
	if password != "" && len(password) > 3 {
		user.Password = password
	} else {
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}
	json.NewDecoder(r.Body).Decode(&user)

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
	mongoClient.InsertMongoDB(userJSON, ColnameUser)
}
