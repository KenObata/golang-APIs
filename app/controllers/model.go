package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) DeleteDuplicate() error {
	collection := db.Client.Database(Dbname).Collection(Colname)
	//first delete test data if exists
	filter := bson.D{{"company", primitive.Regex{Pattern: "Test", Options: ""}}}
	result, err := collection.DeleteMany(context.Background(), filter)
	if result.DeletedCount > 0 {
		log.Println(result.DeletedCount, " test records were deleted.")
	}
	if err != nil {
		return err
	}

	//set filter
	currentTime := time.Now()
	lastMonth := time.Date(currentTime.Year(), currentTime.Month()-1, currentTime.Day(), 0, 0, 0, 0, time.Local).Format("2006-01-02")
	findOptions := bson.D{{"dateadded", bson.D{{"$gt", lastMonth}}}}

	//find all, later find duplicate
	readAll, _ := collection.Find(context.Background(), findOptions)
	var results []JsonJob
	readAll.All(context.Background(), &results)
	for i, result := range results {
		//delete if duplicate
		filterDuplicate := bson.D{{"company", result.Company}, {"title", result.Title}}
		fetchAll, _ := collection.Find(context.Background(), filterDuplicate)
		var fetchedResults []JsonJob
		fetchAll.All(context.Background(), &fetchedResults)

		if len(fetchedResults) > 1 {
			log.Println(i, result.Company, "deleted.")
			//add date to filterDuplicate to delete exact record.
			filterDuplicate = bson.D{{"company", result.Company}, {"title", result.Title}, {"dateadded", result.DateAdded}}
			_, err := collection.DeleteOne(context.Background(), filterDuplicate)
			if err != nil {
				log.Println("err from Delete duplicate: ", err)
			} else {
				log.Println(result.Company, " is deleted for duplicate.")
			}
		}
	}

	return nil
}

func (db *DB) InsertMongoDB(json []byte, table_name string) error {
	log.Println("InsertMongoDB is called.")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	bsonMap := bson.M{}
	// Convert json to bson, which is mongodb document.
	err := bson.UnmarshalExtJSON([]byte(json), false, &bsonMap)
	if err != nil {
		log.Println("error from InsertMongo(), bson.UnmarshalExtJSON")
		log.Println(err)
		return err
	}

	collection := db.Client.Database(Dbname).Collection(table_name)
	log.Println("Mongo DB name:", db.Client.Database(Dbname).Name())
	if table_name == Colname { //table_name==Job
		var episodesFiltered JsonJob
		filter := bson.D{{"company", bsonMap["company"]}, {"title", bsonMap["title"]}}
		collection.FindOne(context.Background(), filter).Decode(&episodesFiltered) // no error handler because no doc found is also an error.

		if len(episodesFiltered.Company) > 0 {
			//return fmt.Errorf("there already exists: %s", bsonMap["company"])
			log.Println("there already exists: ", bsonMap["company"])
		} else {
			log.Println(bsonMap["company"], "will be inserted.")
		}
	}

	if table_name == ColnameUser {
		var result User
		var results []User
		readOne, _ := collection.Find(context.Background(), bson.D{{"email", bsonMap["email"]}})
		readOne.Decode(&result)
		readOne.All(context.Background(), &results)
		log.Println("results:", len(results))
		if len(results) > 0 {
			return fmt.Errorf("%s", "This user is already registered.")
		}
	}

	_, err = collection.InsertOne(ctx, bsonMap)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (db *DB) ReadMongo(user_iput ...string) []JsonJob {
	log.Println("ReadMongo: user input is ", user_iput)
	// get table(=collection)
	collection := db.Client.Database(Dbname).Collection(Colname)

	findOptions := options.Find()
	// Sort by `date` field descending
	findOptions.SetSort(bson.D{{"dateadded", -1}})
	currentTime := time.Now()
	lastMonth := time.Date(currentTime.Year(), currentTime.Month()-1, currentTime.Day(), 0, 0, 0, 0, time.Local).Format("2006-01-02")

	cur, err := collection.Find(context.Background(), bson.D{{"dateadded", bson.D{{"$gt", lastMonth}}}}, findOptions)
	if err != nil {
		log.Println("err from collection.Find()")
		return nil
	}

	if len(user_iput) > 0 {
		cur, err = collection.Find(context.Background(), bson.M{"company": user_iput[0], "dateadded": bson.D{{"$gt", lastMonth}}}, findOptions)
		if err != nil {
			log.Println("err from user input:", err)
			return nil
		}
	}

	var jobs []JsonJob
	var doc JsonJob
	for cur.Next(context.Background()) {
		//var doc JsonJob
		err := cur.Decode(&doc)
		if err != nil {
			log.Println("error at cur.Decode(&doc)")
			return nil
		}
		//append to jobs
		jobs = append(jobs, doc)
	}
	return jobs
}
