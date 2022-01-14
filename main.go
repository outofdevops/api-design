package main

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
	"sync"
)

var validate = validator.New()
var tweets = fetchTweets()
var mu = sync.Mutex{}

type Tweet struct {
	Username string `json:"username" validate:"required"`
	Tweet    string `json:"tweet"`
	Body     string `json:"body"`
	HashTag  string `json:"hash_tag"`
}

func main() {
	validate.RegisterStructValidation(TweetStructValidation, Tweet{})
	http.HandleFunc("/feed", twitterFeed)
	http.ListenAndServe(":3000", nil)
}

func twitterFeed(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var req Tweet
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = validate.Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	persistTweet(req)

	js, _ := json.Marshal(tweets)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func fetchTweets() []Tweet {
	return []Tweet{
		{"PeterMcKinnon", "Happy New Year 🎆", "Happy New Year 🎆", "#2022"},
		{"Programmer", "I ❤️ GoLang", "I ❤️ GoLang", "#coding"},
		{"DancingPanda", "I Love 💃🏼 Dancing", "I Love 💃🏼 Dancing", ""},
		{"GingerBread", "I hate 🥛", "I hate 🥛", ""},
	}
}

func persistTweet(tweet Tweet) {
	if len(tweet.Body) != 0 {
		tweet.Tweet = tweet.Body
	} else {
		tweet.Body = tweet.Tweet
	}
	mu.Lock()
	defer mu.Unlock()
	tweets = append(tweets, tweet)
}

func TweetStructValidation(sl validator.StructLevel) {
	tweet := sl.Current().Interface().(Tweet)

	if len(tweet.Tweet) == 0 && len(tweet.Body) == 0 {
		sl.ReportError(tweet.Tweet, "tweet", "Tweet", "tweet_or_body", "")
		sl.ReportError(tweet.Body, "body", "Body", "tweet_or_body", "")
	}
}
