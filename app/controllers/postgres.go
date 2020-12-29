package controllers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/lib/pq"
)

func Init_db() {
	log.Println("Init_db called.")
	time.Sleep(10 * time.Second)
}

func GetURL2(URL string) error {
	log.Println("GetURL function is called from main.")
	log.Println("URL:", URL)
	//set timeout
	_, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Load the URL
	res, e := http.Get(URL)
	if e != nil {
		log.Println("Error from http.Get(URL)")
		return e
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println("Error from goquery.NewDocumentFromReader.")
		return err
	}
	var urls []string
	var companies []string
	var titles []string

	if strings.Contains(URL, "glassdoor") {
		urls, companies, titles = GetGlassdoor(doc)
	} else {
		doc.Find("a").Each(func(_ int, s *goquery.Selection) {
			url, _ := s.Attr("href")
			title := s.Text()

			if len(companies) < len(titles)-1 {
				log.Println("something is wrong with:", titles)
				os.Exit(1)
			}
			if strings.Contains(url, "/company/") {
				//get company name
				company := s.Text()
				//log.Println("Get URL(), company name:", company)
				companies = append(companies, company)
			} else if strings.Contains(url, "https://ca.linkedin.com/jobs/view/") && len(companies) == len(titles) {
				urls = append(urls, url)
				titles = append(titles, title)
			}
		})
	} //end of scraping

	jsonJob := new(JsonJob)
	for i := 0; i < len(companies); i++ {
		jsonJob.URL = urls[i]
		jsonJob.Title = titles[i]
		jsonJob.Company = companies[i]
		jsonJob.DateAdded = time.Now().Format("2006-01-02")
		// Insert JSON data to MongoDB
		err = InsertJob(jsonJob)
		if err != nil {
			return err
		}
	} //end of passing jsonJob to postgres
	log.Println("End of for loop to insert jsonJobJSON.")
	return nil
}

func connectPostgres() *sql.DB {
	var db *sql.DB
	var err error
	if os.Getenv("MONGO_SERVER") == "" {
		db, err = sql.Open("postgres", "host=localhost user=postgres password=k0668466425 dbname=test_db sslmode=disable") //dbname=test_db
	} else {
		db, err = sql.Open("postgres", "host=db user=postgres password=k0668466425 dbname=test_db sslmode=disable") //dbname=test_db
	}
	if err != nil {
		log.Println("error from sql.Open: ", err)
	} else {
		log.Println("postgres successfully connected.")
		return db
	}
	defer db.Close()
	return nil
}
func InsertJob(jsonJob *JsonJob) error {
	db := connectPostgres()
	//check if record already exists
	rows, err := db.Query("SELECT * FROM job where company=$1 and title=$2;", jsonJob.Company, jsonJob.Title)
	if err != nil {
		log.Println("error from  select: ", err.Error())
	}
	if rows.Next() {
		log.Println("company ", jsonJob.Company, "already exists.")
	} else {
		//insert new record.
		_, err = db.Query("INSERT INTO job(company, title, url, dateadded) SELECT $1,$2,$3,$4;", jsonJob.Company, jsonJob.Title, jsonJob.URL, jsonJob.DateAdded)
		if err != nil {
			log.Println("error from  insert into JOB: ", err.Error())
		}
	}
	return nil
}

func DeleteJobDuplicate() error {
	db := connectPostgres()

	//check if duplicate exists
	rows, err := db.Query("SELECT company, title, url, count(*) FROM job GROUP BY company, title, url HAVING COUNT(*) > 1;")
	if err != nil {
		log.Println("error from  SELECT from job: ", err.Error())
	}
	for rows.Next() {
		var company string
		var title string
		var url string
		var count int
		err = rows.Scan(&company, &title, &url, &count)
		count -= 1
		db.Query("DELETE job where company=$1 and title=$2 and url=$3 LIMIT $4;", company, title, url, count)
	}
	return nil
}

func ReadPostgres(user_iput string, checkSoftware bool, checkDataScience bool, checkThisWeek bool) []JsonJob {
	log.Println("ReadMongo: user filter is ", user_iput)
	db := connectPostgres()

	//first, extract last 1 month records.
	currentTime := time.Now()
	lastMonth := time.Date(currentTime.Year(), currentTime.Month()-1, currentTime.Day(), 0, 0, 0, 0, time.Local).Format("2006-01-02")

	//filter by company name if any.
	var rows *sql.Rows
	var selectErr error
	if len(user_iput) > 0 {
		rows, selectErr = db.Query("SELECT company, title, url, dateadded FROM job WHERE dateadded > $1 and company like '%' || $2 || '%';", lastMonth, user_iput)
	} else {
		rows, selectErr = db.Query("SELECT company, title, url, dateadded FROM job WHERE dateadded > $1;", lastMonth)
	}

	if selectErr != nil {
		log.Println("err from collection.Find()")
		return nil
	}

	var jobs []JsonJob
	var doc JsonJob
	var toBeAdded bool = true
	thisWeek := currentTime.AddDate(0, 0, -7).Format("2006-01-02")
	for rows.Next() {
		var company string
		var title string
		var url string
		var dateadded string
		err := rows.Scan(&company, &title, &url, &dateadded)
		if err != nil {
			log.Println("error from rows.Scan() ", err)
		}
		doc.Company = company
		doc.Title = title
		doc.URL = url
		doc.DateAdded = dateadded

		//condition: software developer only
		if checkSoftware && !(strings.Contains(doc.Title, "Software") || strings.Contains(doc.Title, "Developer") || strings.Contains(doc.Title, "Engineer")) {
			toBeAdded = false
		}
		//condition: this week only
		if checkThisWeek && strings.Compare(doc.DateAdded, thisWeek) < 0 {
			toBeAdded = false
		}
		//condition: data science only
		if checkDataScience && !(strings.Contains(doc.Title, "Data") || strings.Contains(doc.Title, "Analytics")) {
			toBeAdded = false
		}
		if toBeAdded {
			//append to jobs
			jobs = append(jobs, doc)
		}
		toBeAdded = true //reset the flag.
	}
	return jobs
}
