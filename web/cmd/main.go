package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/double-nibble/algothon/web"
)

// TODO: stats
// TODO: Store session IDs in Redis

var log = logrus.WithFields(logrus.Fields{"service": "web"})

func main() {
	ll := flag.String("log_level", "WARN", "Logging level (ERROR|WARN|INFO)")
	flag.Parse()
	loggingInit(ll)
	web.ListenAndServe()
}

func loggingInit(level *string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
	var ll logrus.Level
	switch *level {
	case "DEBUG":
		ll = logrus.DebugLevel
	case "INFO":
		ll = logrus.InfoLevel
	case "WARN":
		ll = logrus.WarnLevel
	case "ERROR":
		ll = logrus.ErrorLevel
	default:
		flag.Usage()
		log.Fatalf("Invalid log level '%s'", *level)
	}
	logrus.SetLevel(ll)
}
