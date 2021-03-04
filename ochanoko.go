package main

//sample json unmarshal:https://tutorialedge.net/golang/go-json-tutorial/

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	API_KEY = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiIxIiwianRpIjoiNjI0M2QxOWU2NTQwMTM2Y2JlYjllZTQ3ZWYyMDg4MjdiMDM2YmQ1OWM2NzBjNTY0NmFjZDQxOWVjOWE5NGZjMWFkMGU0MzQ4NDYzZDNmNDYiLCJpYXQiOjE2MTQ2NTk2OTQsIm5iZiI6MTYxNDY1OTY5NCwiZXhwIjozMTkyNDk2NDk0LCJzdWIiOiI5NDczMyIsInNjb3BlcyI6W119.SkMgCTVCUiDOLZLL_qkicCV5ODRazaaEZfUIy6Q1Or9Vm6CWB7OrHEaKsIbce2IevTM669Kp1YHdZJtgJEceRWVn7tB_o2k5yN-pCsiVPIMlPs9fTnL7gwYaBDYAffH4Eqbf1LQilicyC2D0bwvWjdPkbpwXOOjbHxfDB6rLEAYk7DpZlWu4cS2Agm30zGZUVSz8x6bne8EFgyfNKUfC-ms71d5xT0QqhVZDek8Dxwg2q26gRFbzON3gYqe6VPlb7oW8mX9pgGtLLEAXNxcXZ7590iOghwOQ1gC5GTlX2yC0uTxnt_h3rI5D3ZJtwNFEqnc8sOtnUbjQMaN_FfKd5c66tAWFAuHg1gtX0f9AbrxmRz1lPAiqd-fZM1wBYvdgmD4Y5BaxvguicdgcF2TecxpgKM_45CXlxggP1Fz0M8-MK2kBXcgJGiEJ8sFZxqrUQpPDBwaLASZaujggtzOo0VtgTsMLHZ3XdSfeWfS-l2hGEZGAaVAUkTkLkmMxK-laLLIzrbPtYFUG4h6l1AY3WR37_qPQnUgZcpx1gNxHPROLQWC-jeoYInI7wTmgBA6dQdMXQ66NMRSBgvwpttniInfv8t7v1PRwrSLF1mdsDm4dfuZHvdaC2Va5jZia_CuehW2rjBrh1s_yG015FbBKnDF0Z1FoPmbeKq-f046WFps"
	bearer  = "Bearer " + API_KEY
)

type StockNodes struct {
	Data []Stock `json:"data"`
}

type Stock struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}

func updateStock(id int, stock int) {
	log.Println("---------updateStock()--------")
	//1.create JSON object
	var newStock Stock
	newStock.Id = id
	newStock.Stock = stock
	//json, err := json.Marshal(newStock)
	//log.Println("jsonNewStock:", string(json))
	//---------------------------
	//2.http PUT json object
	URL := "https://api.ocnk.net/v1/products"
	//req, err := http.NewRequest("PUT", URL, bytes.NewBuffer(json))
	//attempt2
	json, _ := os.Open("stock.json")
	content, err := ioutil.ReadFile("stock.json")
	log.Println("json file:", string(content))
	req, err := http.NewRequest(http.MethodPut, URL, json)
	if err != nil {
		log.Print("err from http.NewRequest:", err)
	}
	// set the request header Content-Type for json
	//req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	// Send req using http Client
	client := &http.Client{}
	response, err2 := client.Do(req)
	//defer response.Body.Close()
	if err2 != nil {
		log.Println("err2 from client.Do(req):", err2)
	} else {
		log.Println("Status code:", response.StatusCode, "means:", response.Status)
	}

	/*buf := make([]byte, 64)
	reader, _ := response.Body.Read(buf)
	log.Println("reader:", reader)
	*/
	//---------------------------
}

func main() {

	//var urls [10]string
	productURL := "https://api.ocnk.net/v1/products"

	/*
		for i := 0; i < 10; i++ {
			temp := []string{"https://api.ocnk.net/v1/products/", strconv.Itoa(i)}
			urls[i] = strings.Join(temp, "")
			//log.Println(temp)
			//url := "https://api.ocnk.net/v1/products/2"
		}*/

	//track previous stock number in key value.

	/*Get stock status of all items.*/
	// Create a new request using http
	//for i := 0; i < 10; i++ {

	req, err := http.NewRequest("GET", productURL, nil)
	if err != nil {
		log.Print("err from http.NewRequest:", err)
	}
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	client := &http.Client{}
	response, err2 := client.Do(req)
	defer response.Body.Close()

	//response, err := http.Get(url)
	if err2 != nil {
		log.Println("error from http.Get:", err2)
	} else {
		data, err3 := ioutil.ReadAll(response.Body)
		if err3 != nil {
			log.Println("error from ioutil.ReadAl:", err3)
		}
		log.Println(string(data))
		//unmarshal to Stock structure
		stocks := StockNodes{}
		err4 := json.Unmarshal([]byte(data), &stocks)
		if err4 != nil {
			log.Print("error from unmarshal:", err4)
		}
		log.Println("---------------------")
		//fmt.Printf("%+v\n", stock.Id)
		//need to debug
		log.Println("stocks Data:", stocks.Data)
	}
	//} //end of for loop

	//updateStock(2, 11)
}
