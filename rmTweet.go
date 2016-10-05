package main

import (
	"github.com/ChimeraCoder/anaconda"

	"log"
)

var (
	consumerKey    = "" // FILL ME IN BEFORE USING!!!
	consumerSecret = "" // FILL ME IN BEFORE USING!!!
	accessToken    = "" // FILL ME IN BEFORE USING!!!
	accessSecret   = "" // FILL ME IN BEFORE USING!!!
)

func main() {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessSecret)
	timelineNotEmpty := true
	for timelineNotEmpty {
		timeline, err := api.GetUserTimeline(nil)
		if err != nil {
			log.Fatal(err)
		}
		if len(timeline) > 0 {
			for _, tweet := range timeline {
				tweet, err := api.DeleteTweet(tweet.Id, true)
				if err != nil {
					log.Fatal(err)
				}
				log.Println("[DELETED] " + tweet.Text)
			}
		} else {
			timelineNotEmpty = false
		}
	}
	log.Println("DONE!")
}
