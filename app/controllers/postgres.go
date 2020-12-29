package controllers

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func Init_db() {
	log.Println("Init_db called.")
	//even in docker-compose or kubernetes, first create db manually.

	//connect
	/*var host string
	if os.Getenv("MONGO_SERVER") == "" {
		host = "127.0.0.1"
	} else {
		host = "db" //POSTGRES_HOST
	}
	log.Println("postgres host:", host)
	*/
	//connStr := "postgres://postgres:k0668466425@" + host + ":5432/postgres?sslmode=disable"
	time.Sleep(time.Minute)
	db, err := sql.Open("postgres", "host=db user=postgres password=k0668466425 dbname=test_db sslmode=disable") //dbname=test_db
	//db, err := sql.Open("postgres", connStr) //dbname=test_db
	if err != nil {
		log.Println("error from sql.Open: ", err)
	} else {
		log.Println("postgres successfully connected.")
	}

	//select
	rows, err := db.Query("SELECT * FROM user_list;")
	if err.Error() != "" {
		log.Println("error from  select: ", err.Error())
	}
	for rows.Next() {
		var id int
		var email string
		var password string
		err = rows.Scan(&id, &email, &password)
		log.Println("selected user:", id, email, password)
	}
	defer db.Close()

}
