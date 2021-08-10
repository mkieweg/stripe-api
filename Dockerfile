FROM golang:1.16-alpine AS installer
WORKDIR /src/bin
RUN apk add --no-cache git
COPY go.* .
RUN go mod download

FROM installer as builder
COPY . .
RUN GOOS=linux go build -o api ./cmd/api

FROM alpine
EXPOSE 80
COPY --from=builder /src/bin/api .

ENTRYPOINT [ "./api" ]
CMD [""]
