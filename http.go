package stripeapi

import (
	"context"
	"net/http"
	"time"

	"github.com/mkieweg/stripe-api/payment"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const apiPrefix = "/api/v1"

var store payment.Store
var httpServer *http.Server

func stopHttpServer() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Info("shutting down http server")
	return httpServer.Shutdown(ctx)
}

func registerHttpHandler() {
	http.Handle(apiPrefix+"/register", registerChain)
	http.Handle(apiPrefix+"/subscribe", subscribeChain)
	http.Handle(apiPrefix+"/webhook", webhookChain)
	http.HandleFunc("/healthz", healthCheck)
	http.HandleFunc("/readyz", readynessCheck)
	http.Handle("/metrics", promhttp.Handler())
	store = &payment.StoreImplementation{}
}

func startHttpServer() {
	prometheus.MustRegister(inFlightGauge, counter, duration, responseSize)
	registerHttpHandler()
	httpServer = &http.Server{Addr: ":80"}
	go func() {
		log.Info(httpServer.ListenAndServe())
	}()
}

func healthCheck(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
}

func readynessCheck(writer http.ResponseWriter, request *http.Request) {
	if !ready {
		writer.WriteHeader(http.StatusServiceUnavailable)
	} else {
		writer.WriteHeader(http.StatusOK)
	}
}
