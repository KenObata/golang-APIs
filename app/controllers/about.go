package controllers

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	// Working Directory
	wd, err := os.Getwd()
	t, err := template.ParseFiles(wd + "/app/view/about.html")
	if err != nil {
		log.Println("Error from template.ParseFiles()!")
		log.Println(err)
	}
	t.Execute(w, nil)
}
