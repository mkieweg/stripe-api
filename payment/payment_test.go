package payment

import (
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/mkieweg/stripe-api/config"

	"github.com/stripe/stripe-go/v72"
)

func TestMain(m *testing.M) {
	setupMock()
	exit := m.Run()
	os.Exit(exit)
}

func setupMock() {
	stripe.Key = "sk_test_51JKROHGM7BqXcOjM2TK8MFshg9psVzMWiDe12mzvMLwTlFB7MaYeaTsbYIPzE7dNIVJV0B81P4frPgGOn70hsdl600l4YJhoc3"
	/*
		httpClient := &http.Client{Timeout: time.Second * 1}
		backend := stripe.BackendImplementation{
			Type:              stripe.APIBackend,
			URL:               "http://localhost:12111",
			HTTPClient:        httpClient,
			LeveledLogger:     stripe.DefaultLeveledLogger,
			MaxNetworkRetries: stripe.DefaultMaxNetworkRetries,
		}

		stripe.SetBackend(stripe.APIBackend, &backend)
	*/
}

func TestNewPayment(t *testing.T) {
	type args struct {
		req *PaymentRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				req: &PaymentRequest{
					Number:   "4242424242424242",
					ExpMonth: "8",
					ExpYear:  "2022",
					Cvc:      "123",
				},
			},
		},
		{
			name: "no number",
			args: args{
				req: &PaymentRequest{
					ExpMonth: "8",
					ExpYear:  "2022",
					Cvc:      "123",
				},
			},
			wantErr: true,
		},
		{
			name: "expired",
			args: args{
				req: &PaymentRequest{
					Number:   "4242424242424242",
					ExpMonth: "8",
					ExpYear:  "2020",
					Cvc:      "123",
				},
			},
			wantErr: true,
		},
		{
			name: "no cvc",
			args: args{
				req: &PaymentRequest{
					Number:   "4242424242424242",
					ExpMonth: "8",
					ExpYear:  "2022",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPayment(tt.args.req)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewPayment() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if reflect.DeepEqual(got.CustomerID, "") ||
				reflect.DeepEqual(got.MethodID, "") ||
				reflect.DeepEqual(got.CustomerID, "") {
				t.Errorf("NewPayment() = %v", got)
			}
		})
	}
}

func Test_createPaymentMethod(t *testing.T) {
	type args struct {
		req *PaymentRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createPaymentMethod(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("createPaymentMethod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createPaymentMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createCustomer(t *testing.T) {
	type args struct {
		id *string
	}
	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createCustomer(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("createCustomer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createCustomer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPayment_AttachCustomer(t *testing.T) {
	type fields struct {
		PriceID        string
		MethodID       string
		CustomerID     string
		SubscriptionID string
		UUID           uuid.UUID
	}
	type args struct {
		customerID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Payment{
				PriceID:        tt.fields.PriceID,
				MethodID:       tt.fields.MethodID,
				CustomerID:     tt.fields.CustomerID,
				SubscriptionID: tt.fields.SubscriptionID,
				UUID:           tt.fields.UUID,
			}
			if err := p.AttachCustomer(tt.args.customerID); (err != nil) != tt.wantErr {
				t.Errorf("Payment.AttachCustomer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPayment_CreateSubscription(t *testing.T) {
	type fields struct {
		PriceID        string
		MethodID       string
		CustomerID     string
		SubscriptionID string
		UUID           uuid.UUID
	}
	type args struct {
		params *stripe.SubscriptionParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Payment{
				PriceID:        tt.fields.PriceID,
				MethodID:       tt.fields.MethodID,
				CustomerID:     tt.fields.CustomerID,
				SubscriptionID: tt.fields.SubscriptionID,
				UUID:           tt.fields.UUID,
			}
			got, err := p.CreateSubscription()
			if (err != nil) != tt.wantErr {
				t.Errorf("Payment.CreateSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Payment.CreateSubscription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Functionality(t *testing.T) {
	type args struct {
		req *PaymentRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				req: &PaymentRequest{
					Number:   "4242424242424242",
					ExpMonth: "8",
					ExpYear:  "2022",
					Cvc:      "123",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.InitialiseConfig()
			p, err := NewPayment(tt.args.req)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewPayment() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			s, err := p.CreateSubscription()
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Payment.CreateSubscription() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if s == nil {
				t.Errorf("Payment.CreateSubscription() = %v", s)
			}
		})
	}
}
