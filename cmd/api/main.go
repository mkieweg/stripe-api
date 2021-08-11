package main

import (
	"os"

	stripeapi "github.com/mkieweg/stripe-api"
	log "github.com/sirupsen/logrus"
)

func main() {
	level := os.Getenv("STRIPE_API_LOGLEVEL")
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	case "info":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.ErrorLevel)
		log.SetFormatter(&log.JSONFormatter{})
	}
	log.Info(log.GetLevel())
	stripeapi.Run()
}
