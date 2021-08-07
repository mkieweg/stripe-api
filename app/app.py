import json
import logging
import os
from multiprocessing import Process

import stripe
from flask import Flask, request, jsonify

from datamodel.payment import Store

logging.basicConfig(level=logging.DEBUG,
                    format='[%(asctime)s]: {} %(levelname)s %(message)s'.format(os.getpid()),
                    datefmt='%Y-%m-%d %H:%M:%S',
                    handlers=[logging.StreamHandler()])

logger = logging.getLogger()


def create_app():
    app = Flask(__name__)
    payment_store = Store()

    @app.route('/checkout', methods=['POST'])
    def checkout():
        payload = json.loads(request.data)
        create_payment = Process(target=payment_store.add_payment, args=(payload,))
        create_payment.start()
        return jsonify({'status': 'success'})

    @app.route('/webhook', methods=['POST'])
    def webhook():
        webhook_secret = os.getenv("STRIPE_WEBHOOK_SECRET")
        request_data = json.loads(request.data)

        if webhook_secret:
            signature = request.headers.get('stripe-signature')
            try:
                event = stripe.Webhook.construct_event(
                    payload=request.data, sig_header=signature, secret=webhook_secret)
                data = event['data']
            except Exception as e:
                return e
            event_type = event['type']
        else:
            data = request_data['data']
            event_type = request_data['type']

        data_object = data['object']

        if event_type == 'customer.subscription.created':
            try:
                payment_store.change_status(data_object['id'], data_object['status'])
            except AttributeError as e:
                logging.error(e)

        if event_type == 'customer.subscription.updated':
            try:
                payment_store.change_status(data_object['id'], data_object['status'])
            except Exception as e:
                logging.error(e)

        return jsonify({'status': 'success'})
    return app


if __name__ == "__main__":
    app = create_app()
    app.run(host='0.0.0.0', debug=True)
