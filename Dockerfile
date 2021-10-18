FROM golang:1.17

WORKDIR /code

COPY . .

RUN go mod vendor

RUN go build cmd/main.go

CMD ["./main"]
