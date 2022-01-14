package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

var tweets = fetchTweets()
var mu = sync.Mutex{}

type Tweet struct {
	Username string `json:"username"`
	Tweet    string `json:"tweet"`
}

func main() {
	http.HandleFunc("/feed", twitterFeed)
	http.ListenAndServe(":3000", nil)
}

func twitterFeed(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req Tweet
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	persistTweet(req)

	js, _ := json.Marshal(tweets)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func fetchTweets() []Tweet {
	return []Tweet{
		{"PeterMcKinnon", "Happy New Year 🎆"},
		{"Programmer", "I ❤️ GoLang"},
		{"DancingPanda", "I Love 💃🏼 Dancing"},
		{"GingerBread", "I hate 🥛"},
	}
}

func persistTweet(tweet Tweet) {
	mu.Lock()
	defer mu.Unlock()
	tweets = append(tweets, tweet)
}

