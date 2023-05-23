FROM golang:alpine

WORKDIR /app

COPY go.sum go.mod .

RUN go mod download

COPY . .

RUN go build -o booking-api ./api/main.go

EXPOSE 8080

CMD [ "./booking-api" ]
