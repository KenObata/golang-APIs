package controllers

/*Some info: https://onemuri.space/note/vo6tcv8fq/*/
import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
)

var rClient *redis.Client

func Init() {
	var host string
	if os.Getenv("MONGO_SERVER") == "" {
		host = MongoDBHost
	} else {
		host = os.Getenv("MONGO_SERVER")
	}
	rClient = redis.NewClient(&redis.Options{
		Addr:     host + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	log.Println("Redis Client created.", host+":6379")
}

//userMap={"id": 1, "email":"xxx@gmail.com", "password":"2030"}
func SetKey(userMap bson.M) error {
	//get number of sessions
	session_num := rClient.DBSize().Val() //rClient has ctx already
	var key_str string
	key_str = fmt.Sprint(userMap["id"])
	err := rClient.Set(key_str, session_num+1, 10*time.Second)
	if err != nil {
		log.Println(err)
		return err.Err()
	} else {
		log.Println("SetKey success.")
	}
	return nil
}
