# rmTweet
This application uses the <a href="https://github.com/ChimeraCoder/anaconda">Anaconda</a> library to access a Twitter account via API keys and remove all tweets from it's home timeline. This application was originally designed as a free alternative to online applications such as <a href="https://www.tweetdeleter.com/">tweetdeleter</a> and will probably never be improved beyond it's existing functionality. 

## Usage: 

Either provide API information via command line arguments or in the source itself: 
```
consumerKey    = ""
consumerSecret = ""
accessToken    = ""
accessSecret   = ""
```
Options include:

```
$ rmTweet --help
Usage of rmTweet:
  -as string
        Your Twitter API Access Secret
  -at string
        Your Twitter API Access Token
  -ck string
        Your Twitter API Consumer Key
  -cs string
        Your Twitter API Consumer Secret
  -days int
        Only affect tweets created before n day(s) ago (Default 0)
  -likes
        This will unlike past tweets
  -retweets
        This will un-retweet past tweets (must use along with -tweets option)
  -test
        This will prevent any actual changes from occurring
  -tweets
        This will delete past tweets
```

Example for deleting tweets, retweets, and likes that are older than 60 days:

`go run rmTweet.go -tweets -retweets -likes -days 60`