package controllers

import (
	"context"
	"log"
	"os"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

var rClient *redis.Client
var ctx = context.Background()

func Init_redis() {
	var host string
	if os.Getenv("REDIS_SERVER") == "" {
		host = MongoDBHost
	} else {
		host = os.Getenv("REDIS_SERVER")
	}
	rClient = redis.NewClient(&redis.Options{
		Addr:     host + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	log.Println("Redis Client created.", host+":6379")
}

//userMap={"id": 1, "email":"xxx@gmail.com", "password":"2030"}
func SetKey(ctx context.Context, userId string) (string, error) {
	//get number of sessions
	session_num := rClient.DBSize(ctx).Val() //rClient has ctx already
	log.Println("session_num:", session_num)

	u, err := uuid.NewRandom()
	if err != nil {
		log.Println("error from UUID generation.", err)
		return "", err
	}
	uu := u.String()

	//log.Println("key_str:", key_str)
	setErr := rClient.Set(ctx, uu, userId, 0) //set no expiration.
	//log.Println("result of set:", err.Err(), "| ", err.Val())
	if setErr.Err() != nil {
		log.Println("error from rClient.Set()", setErr.Err())
		return "", setErr.Err()
	} else {
		log.Println("SetKey success: ", uu)
	}
	return uu, nil
}

func GetKey(ctx context.Context, uuid string) (string, error) {

	res, err := rClient.Get(ctx, uuid).Result()
	if err != nil {
		return "", err
	} else {
		log.Println("GetKey: ", res)
	}
	return res, nil
}
