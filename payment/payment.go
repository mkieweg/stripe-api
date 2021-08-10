package payment

import (
	"github.com/google/uuid"
	"github.com/mkieweg/stripe-api/config"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/sub"
)

func NewPayment(req *PaymentRequest) (*Payment, error) {
	p := &Payment{
		PriceID: config.PriceID(),
		UUID:    uuid.New(),
	}

	pmid, err := createPaymentMethod(req)
	if err != nil {
		return nil, err
	}
	p.MethodID = *pmid

	cid, err := createCustomer(pmid)
	if err != nil {
		return nil, err
	}
	p.CustomerID = *cid
	return p, nil
}

type Payment struct {
	PriceID        string    `json:"priceId"`
	MethodID       string    `json:"paymentMethodId"`
	CustomerID     string    `json:"customerId"`
	SubscriptionID string    `json:"subscriptionId"`
	UUID           uuid.UUID `json:"uuid"`
}

func createPaymentMethod(req *PaymentRequest) (*string, error) {
	params := &stripe.PaymentMethodParams{
		Card: &stripe.PaymentMethodCardParams{
			Number:   &req.Number,
			ExpMonth: &req.ExpMonth,
			ExpYear:  &req.ExpYear,
			CVC:      &req.Cvc,
		},
		Type: stripe.String("card"),
	}
	resp, err := paymentmethod.New(params)
	if err != nil {
		return nil, err
	}
	return &resp.ID, nil
}

func createCustomer(id *string) (*string, error) {
	params := &stripe.CustomerParams{
		PaymentMethod: id,
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: id,
		},
	}
	resp, err := customer.New(params)
	if err != nil {
		return nil, err
	}
	return &resp.ID, nil
}

func (p *Payment) AttachCustomer(customerID string) error {
	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customerID),
	}
	_, err := paymentmethod.Attach(p.MethodID, params)
	return err
}

func (p *Payment) CreateSubscription() (*string, error) {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(p.CustomerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(p.PriceID),
			},
		},
	}
	params.AddExpand("latest_invoice.payment_intent")
	resp, err := sub.New(params)
	if err != nil {
		return nil, err
	}
	return &resp.ID, nil
}

type PaymentRequest struct {
	Number   string `json:"number"`
	ExpMonth string `json:"exp_month"`
	ExpYear  string `json:"exp_year"`
	Cvc      string `json:"cvc"`
}
