package stripeapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mkieweg/stripe-api/config"
	"github.com/mkieweg/stripe-api/payment"
	log "github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
)

var registerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	var req *payment.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	p, err := payment.NewPayment(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	json.NewEncoder(w).Encode(p)
})

var subscribeHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var req *payment.Payment
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Errorf("json.NewDecoder.Decode: %v", err)
		return
	}

	p, err := store.GetPayment(req.UUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Errorf("store.GetPayment: %v", err)
		return
	}
	if err := p.AttachCustomer(req.CustomerID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Errorf("paymentmethod.Attach: %v %s", err, req.MethodID)
		return
	}

	_, err = p.CreateSubscription()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Errorf("sub.New: %v", err)
		return
	}
})

var webhookHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event := stripe.Event{}
	if err := json.Unmarshal(payload, &event); err != nil {
		log.Errorf("⚠️  Webhook error while parsing basic request. %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	endpointSecret := config.WebhookSecret()
	signatureHeader := r.Header.Get("Stripe-Signature")
	event, err = webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
		log.Errorf("⚠️  Webhook signature verification failed. %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "customer.subscription.created":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			log.Errorf("Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Infof("Subscription created %v", subscription.ID)

	case "customer.subscription.updated":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			log.Errorf("Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Infof("Subscription updated %v", subscription.ID)

	default:

		log.Errorf("Unhandled event type: %s\n", event.Type)

	}

	w.WriteHeader(http.StatusOK)
})
