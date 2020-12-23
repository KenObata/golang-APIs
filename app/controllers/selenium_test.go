package controllers

import (
	"log"
	"testing"
	"time"

	"github.com/sclevine/agouti"
)

//some useful info:https://qiita.com/tenten0213/items/1f897ff8a64bd8b5270c

/*
func init() { //TestMain(m *testing.M)
	log.Println("init() called.")

	//crete DB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + "127.0.0.1" + ":" + MongoDBPort))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	db := client.Database(Dbname)

	//create user in test DB
	err_createUser := db.RunCommand(ctx, bson.D{{"createUser", MongoUser},
		{"pwd", MongoPassword},
		{"roles", []bson.M{{"role": "readWrite", "db": "test"}}}})
	if err_createUser != nil {
		log.Println(err_createUser)
	}
	err_createCollection_Job := db.CreateCollection(ctx, Colname)
	if err_createCollection_Job != nil {
		log.Println(err_createCollection_Job)
	}
	err_createCollection_User := db.CreateCollection(ctx, ColnameUser)
	if err_createCollection_User != nil {
		log.Println(err_createCollection_User)
	}
}
*/

//var agoutiDriver *agouti.WebDriver

func TestSignUpHandler(t *testing.T) {
	driver := agouti.ChromeDriver(agouti.Browser("chrome"))
	err := driver.Start()
	if err != nil {
		log.Fatalf("Failed to start driver:%v", err)
	}
	defer driver.Stop()

	page, err1 := driver.NewPage()
	if err1 != nil {
		log.Fatalf("Failed to open page:%v", err)
	}

	//go to index.html
	err2 := page.Navigate("http://localhost:8080")
	if err2 != nil {
		log.Fatalf("Failed to navigate:%v", err2)
	}
	time.Sleep(1000)
	page.Find("signup").Click()
}

/*
func TestInternalHandler(t *testing.T) {

}
*/
