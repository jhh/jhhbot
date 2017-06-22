package main

import (
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/Sirupsen/logrus"
)

var (
	consumerKey       = getenv("TWITTER_CONSUMER_KEY")
	consumerKeySecret = getenv("TWITTER_CONSUMER_SECRET")
	accessToken       = getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = getenv("TWITTER_ACCESS_TOKEN_SECRET")
	log               = &logger{logrus.New()}
)

func getenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable missing: %s", key)
	}
	return v
}

func main() {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerKeySecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	api.SetLogger(log)

	stream := api.PublicStreamFilter(url.Values{
		"follow": []string{"622020160"},
	})
	defer stream.Stop()

	for v := range stream.C {
		t, ok := v.(anaconda.Tweet)
		if !ok {
			log.Warningf("received unexpected value of type %T", v)
			continue
		}

		if t.RetweetedStatus != nil {
			continue
		}

		_, err := api.Retweet(t.Id, false)
		if err != nil {
			log.Errorf("could not retweet %d: %v", t.Id, err)
			continue
		}
		log.Noticef("retweeted %d", t.Id)
	}
	log.Fatal("stream was not opened")
}
