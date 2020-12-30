package controllers

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Job struct {
	URL     []string
	Title   []string
	Company []string
}

type JsonJob struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	Company   string `json:"company"`
	DateAdded string `json:"dateadded"`
}

// Array for Job struct
type Jobs []*Job

// mongo-driverのクライアントを自前で定義した構造体DBへセット
type DB struct {
	Client *mongo.Client
}

const (
	// 接続先のDB情報を入力
	MongoDBHost       = "127.0.0.1" //mongodb.default.svc.cluster.local.
	MongoDBPort       = "27017"     //"27016" is used by dev
	MongoUser         = "Ken"
	MongoPassword     = "k0668466425"
	Dbname            = "test" //"databases"
	Colname           = "Job"
	ColnameUser       = "User"
	POSTGRES_HOST     = "host.docker.internal"
	POSTGRES_USER     = "postgres"
	POSTGRES_PASSWORD = "k0668466425"
)

type User struct {
	//Capital letter means public
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
