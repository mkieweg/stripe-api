import os

import stripe


class Payment:
    def __init__(self, number="", exp_month=0, exp_year=0, cvc=0):
        self.price_id = os.getenv("STRIPE_PRICE_ID")
        resp = create_payment_method(number, exp_month, exp_year, cvc)
        self.method_id = resp
        resp = create_customer(self.method_id)
        self.customer_id = resp
        update_payment_method(self.method_id, self.customer_id)
        self.subscription_id = create_subscription(self.customer_id, self.price_id)


def create_payment_method(number="", exp_month=0, exp_year=0, cvc=0):
    resp = stripe.PaymentMethod.create(
        api_key=os.getenv("STRIPE_API_KEY"),
        type="card",
        card={
            "number": number,
            "exp_month": exp_month,
            "exp_year": exp_year,
            "cvc": cvc
        },
    )
    return resp['id']


def create_customer(method_id):
    resp = stripe.Customer.create(
        api_key=os.getenv("STRIPE_API_KEY"),
        payment_method=method_id,
        invoice_settings={
            "default_payment_method": method_id
        }
    )
    return resp['id']


def update_payment_method(method_id, customer_id):
    stripe.PaymentMethod.attach(
        method_id,
        api_key=os.getenv("STRIPE_API_KEY"),
        customer=customer_id
    )


def create_subscription(customer_id, price_id):
    resp = stripe.Subscription.create(
        api_key=os.getenv("STRIPE_API_KEY"),
        customer=customer_id,
        items=[
            {"price": price_id},
        ],
    )
    return resp['id']


class Store:
    def __init__(self):
        super().__init__()
        self.payments = {}

    def add_payment(self, data):
        number = data['number']
        exp_month = data['exp_month']
        exp_year = data['exp_year']
        cvc = data['cvc']
        payment = Payment(number, exp_month, exp_year, cvc)
        self.payments[payment.subscription_id] = payment

    def change_status(self, subscription_id, status):
        self.payments[subscription_id].subscription_state = status
