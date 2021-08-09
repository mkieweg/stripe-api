class MockStore:
    def __init__(self):
        super().__init__()
        self.payments = {}

    def add_payment(self, data):
        number = data['number']
        exp_month = data['exp_month']
        exp_year = data['exp_year']
        cvc = data['cvc']
        self.payments[number] = {exp_month, exp_year, cvc}

    def change_status(self, subscription_id, status):
        self.payments[subscription_id].subscription_state = status
