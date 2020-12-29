package controllers

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	// Working Directory
	wd, err1 := os.Getwd()
	if err1 != nil {
		log.Println("Error from SignUpHandler.", err1)
	}
	t, err := template.ParseFiles(wd + "/app/view/signup.html")
	if err != nil {
		log.Println("Error from SignUpHandler. (err)", err)
	}
	t.Execute(w, nil)

	//get email from html file
	email := r.FormValue("email")
	if email != "" {
		log.Println("email:", email)
		user.Email = email
	} else {
		log.Println("signup page called without input, return.")
		return
	}
	password := r.FormValue("password")
	if password != "" && len(password) > 3 {
		user.Password = password
	} else {
		log.Println("signup page called without input, return.")
		return
	}
	log.Println("email:", user.Email, "password:", user.Password)

	//now it is guaranteed that email and pw are not nil.
	err = InsertUser(user)
	if err != nil {
		//http.Error(w, err.Error(), 500)
		new_t, _ := template.ParseFiles(wd + "/app/view/signup-error.html")
		new_t.Execute(w, err)
	} else {
		new_t, _ := template.ParseFiles(wd + "/app/view/signup-success.html")
		new_t.Execute(w, nil)
	}
}
