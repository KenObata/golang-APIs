package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var (
	query      = flag.String("query", "Google", "pets") //Search term ->pets
	maxResults = flag.Int64("max-results", 25, "Max YouTube results")
)

const developerKey = "AIzaSyDwzOKnr6JtWWE9n2_GMmfpGwA4RR4EuUE"

func main() {
	flag.Parse()

	client := &http.Client{
		//service, err := NewService(developerKey),
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// Make the API call to YouTube.
	//List returns *SearchListCall
	call := service.Search.List("id, snippet").Q(*query).MaxResults(*maxResults)
	response, err := call.Do()
	handleError(err, "")

	// Group video, channel, and playlist results in separate lists.
	videos := make(map[string]string)
	channels := make(map[string]string)
	playlists := make(map[string]string)

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
		case "youtube#channel":
			channels[item.Id.ChannelId] = item.Snippet.Title
		case "youtube#playlist":
			playlists[item.Id.PlaylistId] = item.Snippet.Title
		}
	}

	//printIDs("Videos", videos)
	printIDs("Channels", channels)
	//printIDs("Playlists", playlists)
}

// Print the ID and title of each result in a list as well as a name that
// identifies the list. For example, print the word section name "Videos"
// above a list of video search results, followed by the video ID and title
// of each matching video.
func printIDs(sectionName string, matches map[string]string) {
	fmt.Printf("%v:\n", sectionName)
	for id, title := range matches {
		fmt.Printf("[%v] %v\n", id, title)
	}
	fmt.Printf("\n\n")
}

/*
func main() {
	log.Println("os.Getenv:", os.Getenv("MONGO_SERVER"))

	server := http.Server{}
	if os.Getenv("MONGO_SERVER") == "" {
		server.Addr = ":8080"
	}

		http.HandleFunc("/", controllers.HomeHandler)
		http.HandleFunc("/index-login-success", controllers.HomeHandlerAfterLogin)
		http.HandleFunc("/signup", controllers.SignUpHandler)
		http.HandleFunc("/login", controllers.LoginHandler)
		http.HandleFunc("/internal", controllers.InternalHandler)
		http.HandleFunc("/userpost", controllers.PostHandler)
		http.HandleFunc("/about", controllers.AboutHandler)
		//add css below
		http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css")))) //http.Handle("/css/")

		controllers.Init_db()    //initialize Postgres
		controllers.Init_redis() //initialize Redis

	var url [1]string
	url[0] = "https://www.youtube.com/results?search_query=pet&sp=EgIQAg"

	for i := range url {
		err := controllers.GetURL2(url[i])
		if err != nil {
			log.Println("error from GetURL(): ", err.Error())
		}
	}

	//err := server.ListenAndServe()
	//if err != nil {
	//	log.Fatal("error from ListenAndServe()", err)
	//}
}*/
