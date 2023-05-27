FROM golang:alpine

RUN apk update && apk add --no-cache git

RUN mkdir /app
WORKDIR /app

COPY . .

RUN go mod tidy

RUN go mod vendor

RUN go build -o /build ./cmd/main.go

EXPOSE 3000


CMD ["/build"]