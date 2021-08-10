package stripeapi

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mkieweg/stripe-api/config"
	log "github.com/sirupsen/logrus"
)

var stopChan chan os.Signal

func init() {
	stopChan = make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
}

// Run bootstraps the orchestrator and waits for the shutdown signal
func Run() {
	go func() {
		config.InitialiseConfig()
		log.Info("starting to serve")
		if err := serve(); err != nil {
			log.Fatal(err)
		}
	}()

	log.WithFields(log.Fields{}).Info("initialisation finished")
	<-stopChan
}
