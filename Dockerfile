FROM golang:1.15-alpine

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./app .

CMD ["./app"]