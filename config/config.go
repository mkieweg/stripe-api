package config

import (
	"os"
	"path/filepath"

	"github.com/stripe/stripe-go/v72"
)

var priceID string
var apiKey string
var tlsPath string
var webhookSecret string

func InitialiseConfig() {
	tlsPath = os.Getenv("TLS_PATH")
	priceID = os.Getenv("STRIPE_PRICE_ID")
	apiKey = os.Getenv("STRIPE_API_KEY")
	webhookSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")
	stripe.Key = ApiKey()
}

func PriceID() string {
	return priceID
}

func ApiKey() string {
	return apiKey
}

func WebhookSecret() string {
	return webhookSecret
}
func TlsCert() string {
	return filepath.Join(tlsPath, "cert.pem")
}

func TlsKey() string {
	return filepath.Join(tlsPath, "key.pem")
}
