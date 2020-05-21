FROM golang:alpine

RUN apk update && apk add --no-cache git


WORKDIR /app
COPY . /app

RUN go get -d -v

RUN go build -o server server.go

CMD ["./server"]

