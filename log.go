package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/Sirupsen/logrus"
	"github.com/coreos/go-systemd/journal"
)

// Logger logs messages to systemd journal if available, otherwise to STDOUT.
type Logger struct {
	log *logrus.Logger
}

// Fatal logs messages at log level Fatal and exits.
func (l *Logger) Fatal(args ...interface{}) {
	if l.doJournal(journal.PriEmerg, "%s", args...) {
		os.Exit(1)
	}
	l.log.Fatal(args...)
}

// Fatalf logs messages at log level Fatal and exits.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	if l.doJournal(journal.PriEmerg, format, args...) {
		os.Exit(1)
	}
	l.log.Fatalf(format, args...)
}

// Panic logs messages at log level Panic and panics.
func (l *Logger) Panic(args ...interface{}) {
	l.log.Panic(args...)
}

// Panicf logs messages at log level Panic and panics.
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.log.Panicf(format, args...)
}

// Critical logs messages at log level Critical.
func (l *Logger) Critical(args ...interface{}) {
	if l.doJournal(journal.PriCrit, "%s", args...) {
		return
	}
	l.log.Error(args...)
}

// Criticalf logs messages at log level Critical.
func (l *Logger) Criticalf(format string, args ...interface{}) {
	if l.doJournal(journal.PriCrit, format, args...) {
		return
	}
	l.log.Errorf(format, args...)
}

// Error logs messages at log level Error.
func (l *Logger) Error(args ...interface{}) {
	if l.doJournal(journal.PriErr, "%s", args...) {
		return
	}
	l.log.Error(args...)
}

// Errorf logs messages at log level Error.
func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.doJournal(journal.PriErr, format, args...) {
		return
	}
	l.log.Errorf(format, args...)
}

// Warning logs messages at log level Warning.
func (l *Logger) Warning(args ...interface{}) {
	if l.doJournal(journal.PriWarning, "%s", args...) {
		return
	}
	l.log.Warn(args...)
}

// Warningf logs messages at log level Warning.
func (l *Logger) Warningf(format string, args ...interface{}) {
	if l.doJournal(journal.PriWarning, format, args...) {
		return
	}
	l.log.Warnf(format, args...)
}

// Notice logs messages at log level Notice.
func (l *Logger) Notice(args ...interface{}) {
	if l.doJournal(journal.PriNotice, "%s", args...) {
		return
	}
	l.log.Info(args...)
}

// Noticef logs messages at log level Notice.
func (l *Logger) Noticef(format string, args ...interface{}) {
	if l.doJournal(journal.PriNotice, format, args...) {
		return
	}
	l.log.Infof(format, args...)
}

// Info logs messages at log level Info.
func (l *Logger) Info(args ...interface{}) {
	if l.doJournal(journal.PriInfo, "%s", args...) {
		return
	}
	l.log.Info(args...)
}

// Infof logs messages at log level Info.
func (l *Logger) Infof(format string, args ...interface{}) {
	if l.doJournal(journal.PriInfo, format, args...) {
		return
	}
	l.log.Infof(format, args...)
}

// Debug logs messages at log level Debug. It will not log to systemd journal.
func (l *Logger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

// Debugf logs messages at log level Debug. It will not log to systemd journal.
func (l *Logger) Debugf(format string, args ...interface{}) {
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

func (l *Logger) doJournal(pri journal.Priority, format string, args ...interface{}) bool {
	if !journal.Enabled() {
		return false
	}
	if err := journal.Print(pri, format, args...); err != nil {
		l.log.Errorf("error printing to systemd journal: %v", err)
	}
	return true
}
