package tippler

import (
	"encoding/json"
	"fmt"
	"github.com/darkhelmet/twitterstream"
	"github.com/joho/godotenv"
	"github.com/kellydunn/golang-geo"
	"log"
	"os"
	"strings"
)

var SINGAPORE_POINT1 = twitterstream.Point{1.227781, 103.602909}
var SINGAPORE_POINT2 = twitterstream.Point{1.461176, 104.015583}

type broadcastMsg struct {
	msg     string
	channel string
}

type tipplerTweet struct {
	Msg  string `json:"msg"`
	User string `json:"user"`
	Loc  string `json:"loc"`
	Time string `json:"time"`
}

func RunStream() {
	_ = godotenv.Load()
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	tw := twitterstream.NewClient(consumerKey, consumerSecret, accessToken, accessTokenSecret)
	conn, err := tw.Locations(SINGAPORE_POINT1, SINGAPORE_POINT2)
	if err == nil {
		for {
			if tweet, err := conn.Next(); err == nil {
				if tweet.RetweetedStatus == nil && tweet.InReplyToUserId == nil {
					if tweet.Coordinates != nil {
						coords := *tweet.Coordinates
						if area := FindClosestArea(geo.NewPoint(float64(coords.Long), float64(coords.Lat))); area != nil {
							msg := fmt.Sprintf("[%s] @%s: %s\n", area.Name, tweet.User.ScreenName, tweet.Text)
							log.Print(msg)

							tt := tipplerTweet{
								Msg:  tweet.Text,
								User: tweet.User.ScreenName,
								Loc:  area.Name,
								Time: fmt.Sprintf("%02d:%02d", tweet.CreatedAt.Hour(), tweet.CreatedAt.Minute()),
							}

							json, err := json.Marshal(tt)
							if err != nil {
								log.Println(err)
								continue
							}
							channel := strings.ToLower(area.Name)
							channel = strings.Replace(channel, " ", "_", -1)

							h.broadcast <- &broadcastMsg{msg: string(json), channel: channel}
						}
					}
				}
			}
		}
	} else {
		log.Fatal(err)
	}
}
