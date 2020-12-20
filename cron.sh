cd /Users/kenobata/go/src/Scraping
echo $PATH >> /Users/kenobata/Desktop/log-cron.txt
echo "/Users/kenobata/go/src/Scraping/cron.sh is called" >> /Users/kenobata/Desktop/log-cron.txt
pwd >> /Users/kenobata/Desktop/log-cron.txt
mongoexport --uri="mongodb://127.0.0.1:27017" --collection=Job --out=/Users/kenobata/go/src/Scraping/app/Job.json --pretty --jsonArray
#chmod 777 /Users/kenobata/go/src/Scraping/app/Job.json
kubectl cp /Users/kenobata/go/src/Scraping/app/Job.json builder-67c8b5fbc5-7sgv5:/app/Job.json
cd /Users/kenobata/go/src/Scraping
#go run /Users/kenobata/go/src/Scraping/.