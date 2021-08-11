package stripeapi

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mkieweg/stripe-api/config"
	log "github.com/sirupsen/logrus"
)

var stopChan chan os.Signal
var ready bool

func init() {
	stopChan = make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
}

// Run bootstraps the orchestrator and waits for the shutdown signal
func Run() {
	config.InitialiseConfig()
	log.Info("starting to serve")
	startHttpServer()
	ready = true
	log.WithFields(log.Fields{}).Info("initialisation finished")
	<-stopChan
	if err := stopHttpServer(); err != nil {
		log.Error(err)
	}
}
