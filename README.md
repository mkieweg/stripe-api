# Stripe API Demo

This application is a demo for a Stripe subscription payment process. It provides 3 API endpoints used to create a user, subscribe to a product, and process incoming Stripe webhooks.

A Makefile is provided for local verification deployment using `minikube`. A Docker deployment can be created using the `docker compose` command. The `.yml` files in the `kubernetes` directory can be used to deploy the application to an existing k8s cluster.

## API

The API endpoints are defined in the OpenAPI 3 specification in `doc/openapi.yml`. The server is configured to listen on poer 80 by default.

### `/api/v1/register`

Expects a JSON object containing `number`, `exp_month`, `exp_year`, and `cvc` strings. Returns a JSON object containing `priceId`, `paymentMethodId`, `customerId` and `uuid` strings.

### `/api/v1/subscribe`

Expects a JSON object containing `priceId`, `paymentMethodId`, `customerId` and `uuid` strings.

### `/api/v1/webhook`

Expects a Stripe webhook and processes it.

## Run Unit Tests

To run the unit tests certain parameters need to be set using environment variables

```sh
$ STRIPE_API_KEY=<stripe secret key> \
    STRIPE_PRICE_ID=<price id> \
    STRIPE_WEBHOOK_SECRET=<webhook secret> \
    make test
```

For the tests to pass valid values need to be passed to `Test_createSubscription` in `server_test.go`.

## Run Locally (Docker Compose)

Run `docker compose up` to build the Docker image and run a Container. This more lightweight than using `minikube`. To use Docker Compose a `.env` file containing the following values is required:

```env
STRIPE_API_KEY=<stripe secret key>
STRIPE_PRICE_ID=<price id>
STRIPE_WEBHOOK_SECRET=<webhook secret>
```

## Run Locally (Minikube)

The application is built, deployed, and run locally using the `make` command.

### Requirements

The local deployment needs `minikube` and `kubectl` installed and expects a `secrets.yml` file to be present in the `kubernetes` directory. A sample file is shown below. The required Tokens and IDs need to be base64 encoded.

```yml
apiVersion: v1
kind: Secret
metadata:
  name: stripe-credentials
type: Opaque
data:
  api-token: <base64 encoded token>
  webhook-token: <base64 encoded token>
  price-id: <base64 encoded token>
```

### `make build`

The Docker container can either be built using the GitHub Action defined in the `.github` directory or using `make build` with the GitHub account name set as environment variable.

```sh
REPO=<github account> make build
```

### `make init`

The `minikube` cluster is initialised using the `make init` macro with the same environment variable as used in the build step. It echoes the ingress IP address upon completion.

```sh
REPO=<github account> make init
```

### `make create`

This macro creates the Deployment, Service, and Ingress used. It requires no further parameters.

### `make delete`

This macro deletes the Deployment, Service, and Ingress used. It requires no further parameters.

### `make stop`

This macro is only an alias for `minikube stop`.

## Security considerations

The API access is neither authenticated nor encrypted. It should not be used as is. For encryption a HTTPS reverse proxy could be considered.
