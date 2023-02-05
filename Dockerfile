FROM golang:alpine as go
WORKDIR /app
ENV GO111MODULE=on

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o private-message-service ./cmd

FROM alpine

WORKDIR /app

COPY --from=go /app/private-message-service ./private-message-service
COPY run/config.yaml ./config.yaml
CMD ["./private-message-service"]