package stripeapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mkieweg/stripe-api/config"
)

func Test_createCustomer(t *testing.T) {
	type args struct {
		payload []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "default",
			args: args{
				payload: []byte(`{"number": "4242424242424242", "exp_month": "8", "exp_year": "2022", "cvc": "314"}`),
			},
			want: http.StatusOK,
		},
		{
			name: "no payload",
			args: args{
				payload: []byte(`{}`),
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "expired card",
			args: args{
				payload: []byte(`{"number": "4242424242424242", "exp_month": "8", "exp_year": "2020", "cvc": "314"}`),
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.InitialiseConfig()
			req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(tt.args.payload))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.Handler(registerHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.want {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want)
			}
		})
	}
}

func Test_createSubscription(t *testing.T) {
	type args struct {
		payload []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "default",
			args: args{
				payload: []byte(`{"paymentMethodId": "pm_1JMZYJGM7BqXcOjMDyE77Pl9", "customerId": "cus_K0ajfXI7Vntewh", "priceId": "price_1JKRj0GM7BqXcOjMdnMpsj2t"}`),
			},
			want: http.StatusOK,
		},
		{
			name: "no payload",
			args: args{
				payload: []byte(`{}`),
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "expired card",
			args: args{
				payload: []byte(`{"number": "4242424242424242", "exp_month": "8", "exp_year": "2020", "cvc": "314"}`),
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.InitialiseConfig()
			req, err := http.NewRequest("POST", "/subscribe", bytes.NewBuffer(tt.args.payload))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.Handler(subscribeHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.want {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want)
			}
		})
	}
}
