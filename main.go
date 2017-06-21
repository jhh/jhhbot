package main

import (
	"fmt"
	"net/url"
	"os"
	"reflect"

	"github.com/ChimeraCoder/anaconda"
	"github.com/Sirupsen/logrus"
	"github.com/coreos/go-systemd/journal"
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
}

type logger struct {
	log *logrus.Logger
}

func (l *logger) Fatal(args ...interface{}) {
	if journal.Enabled() {
		if err := journal.Print(journal.PriEmerg, "%s", msg(args...)); err != nil {
			l.log.Errorf("error printing to systemd journal: %v", err)
		}
		os.Exit(1)
	}
	l.log.Fatal(args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	if journal.Enabled() {
		if err := journal.Print(journal.PriEmerg, format, args...); err != nil {
			l.log.Errorf("error printing to systemd journal: %v", err)
		}
		os.Exit(1)
	}
	l.log.Fatalf(format, args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.log.Panic(args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.log.Panicf(format, args...)
}

func (l *logger) Critical(args ...interface{}) {
	if journal.Enabled() {
		if err := journal.Print(journal.PriCrit, "%s", msg(args...)); err != nil {
			l.log.Errorf("error printing to systemd journal: %v", err)
		}
		return
	}
	l.log.Error(args...)
}

func (l *logger) Criticalf(format string, args ...interface{}) {
	if journal.Enabled() {
		if err := journal.Print(journal.PriCrit, format, args...); err != nil {
			l.log.Errorf("error printing to systemd journal: %v", err)
		}
		return
	}
	l.log.Errorf(format, args...)
}

func (l *logger) Error(args ...interface{}) {
	if journal.Enabled() {
		if err := journal.Print(journal.PriErr, "%s", msg(args...)); err != nil {
			l.log.Errorf("error printing to systemd journal: %v", err)
		}
		return
	}
	l.log.Error(args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	if journal.Enabled() {
		if err := journal.Print(journal.PriErr, format, args...); err != nil {
			l.log.Errorf("error printing to systemd journal: %v", err)
		}
		return
	}
	l.log.Errorf(format, args...)
}

func (l *logger) Warning(args ...interface{}) {
	if journal.Enabled() {
		if err := journal.Print(journal.PriWarning, "%s", msg(args...)); err != nil {
			l.log.Errorf("error printing to systemd journal: %v", err)
		}
		return
	}
	l.log.Warn(args...)
}

func (l *logger) Warningf(format string, args ...interface{}) {
	if journal.Enabled() {
		if err := journal.Print(journal.PriWarning, format, args...); err != nil {
			l.log.Errorf("error printing to systemd journal: %v", err)
		}
		return
	}
	l.log.Warnf(format, args...)
}

func (l *logger) Notice(args ...interface{}) {
	if journal.Enabled() {
		if err := journal.Print(journal.PriNotice, "%s", msg(args...)); err != nil {
			l.log.Errorf("error printing to systemd journal: %v", err)
		}
		return
	}
	l.log.Info(args...)
}

func (l *logger) Noticef(format string, args ...interface{}) {
	if journal.Enabled() {
		if err := journal.Print(journal.PriNotice, format, args...); err != nil {
			l.log.Errorf("error printing to systemd journal: %v", err)
		}
		return
	}
	l.log.Infof(format, args...)
}

func (l *logger) Info(args ...interface{}) {
	if journal.Enabled() {
		if err := journal.Print(journal.PriInfo, "%s", msg(args...)); err != nil {
			l.log.Errorf("error printing to systemd journal: %v", err)
		}
		return
	}
	l.log.Info(args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	if journal.Enabled() {
		if err := journal.Print(journal.PriInfo, format, args...); err != nil {
			l.log.Errorf("error printing to systemd journal: %v", err)
		}
		return
	}
	l.log.Infof(format, args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func msg(a ...interface{}) string {
	var msg string
	prevString := false
	for argNum, arg := range a {
		isString := arg != nil && reflect.TypeOf(arg).Kind() == reflect.String
		// Add a space between two non-string arguments.
		if argNum > 0 && !isString && !prevString {
			msg += " "
		}
		msg += fmt.Sprint(arg)
		prevString = isString
	}
	return msg
}