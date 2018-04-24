package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

var (
	consumerKey    = "" // FILL ME IN BEFORE USING!!! (Or provide via -ck argument)
	consumerSecret = "" // FILL ME IN BEFORE USING!!! (Or provide via -cs argument)
	accessToken    = "" // FILL ME IN BEFORE USING!!! (Or provide via -at argument)
	accessSecret   = "" // FILL ME IN BEFORE USING!!! (Or provide via -as argument)
)

var days int
var likes, tweets, retweets, test bool
var now = time.Now().Format("20060102150405")
var logDir = "./log/"

func remove(api *anaconda.TwitterApi, object string) {
	api.EnableThrottling(1*time.Second, 1)

	var err error
	var max int64
	var tweets []anaconda.Tweet
	var deleteTweets []anaconda.Tweet

	v := url.Values{}
	v.Set("count", "200")

	for {
		if object == "tweets" {
			if retweets {
				v.Set("include_rts", "1")
			} else {
				v.Set("include_rts", "0")
			}
			tweets, err = api.GetUserTimeline(v)
		} else if object == "likes" {
			tweets, err = api.GetFavorites(v)
		}
		if err != nil {
			log.Fatal(err)
		}
		if len(tweets) > 0 {
			for _, tweet := range tweets {
				tweetTime, _ := tweet.CreatedAtTime()
				if tweetTime.Before(time.Now().AddDate(0, 0, days)) {
					if !test {
						if object == "tweets" {
							_, err = api.DeleteTweet(tweet.Id, true)
						} else if object == "likes" {
							_, err = api.Unfavorite(tweet.Id)
						}
					}
					if err != nil {
						log.Printf("[ERROR] Could not delete: %v - %v: %v \n", tweet.Id, tweet.CreatedAt, tweet.Text)
					} else {
						deleteTweets = append(deleteTweets, tweet)
						log.Printf("[DELETED] %v: %v \n", tweet.CreatedAt, tweet.Text)
					}
				}
				max = tweet.Id
			}
			v.Set("max_id", strconv.FormatInt(max-1, 10))
		} else {
			break
		}
	}

	if len(deleteTweets) > 0 {
		deleteTweetsJson, _ := json.Marshal(deleteTweets)
		if object == "tweets" {
			err = writeHistory("deleted_tweets", deleteTweetsJson)
		} else if object == "likes" {
			err = writeHistory("unliked_tweets", deleteTweetsJson)
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}

func writeHistory(filename string, data []byte) error {
	fileName := logDir + filename + "_" + now + ".json"
	history, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return errors.New("Error creating or accessing history file: " + fileName)
	}
	defer history.Close()

	_, err = history.Write(data)
	if err != nil {
		return errors.New("Error writing to history file: " + fileName)
	}
	return nil
}

// Handles command line arguments
func handleFlags() {
	flag.StringVar(&consumerKey, "ck", consumerKey, "Your Twitter API Consumer Key")
	flag.StringVar(&consumerSecret, "cs", consumerSecret, "Your Twitter API Consumer Secret")
	flag.StringVar(&accessToken, "at", accessToken, "Your Twitter API Access Token")
	flag.StringVar(&accessSecret, "as", accessSecret, "Your Twitter API Access Secret")
	flag.BoolVar(&tweets, "tweets", false, "This will delete past tweets")
	flag.BoolVar(&retweets, "retweets", false, "This will un-retweet past tweets (must use along with -tweets option)")
	flag.BoolVar(&likes, "likes", false, "This will unlike past tweets")
	flag.BoolVar(&test, "test", false, "This will prevent any actual changes from occurring")
	flag.IntVar(&days, "days", 0, "Only affect tweets created before n day(s) ago (Default 0)")
	flag.Parse()
}

// Configures the logger
func setupLogger(logfile string) *os.File {
	// Ensure log directory exists
	_ = os.Mkdir(logDir, 0755)

	logfile = logDir + logfile

	// Create log file
	file, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("[*] Error creating log file:", logfile)
	}
	log.SetOutput(file)
	return file
}

func main() {
	// Parse command line arguments
	handleFlags()

	days = int(math.Copysign(float64(days), -1))

	// Configure settings for logging
	l := setupLogger("summary_output_" + now + ".txt")
	defer l.Close()

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessSecret)

	if tweets {
		remove(api, "tweets")
	}
	if likes {
		remove(api, "likes")
	}

	fmt.Println("DONE!")
}
