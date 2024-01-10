FROM golang:1.21

WORKDIR /go/src/genericforum

COPY go.mod ./

RUN go mod download

COPY . ./

RUN go build -o main

EXPOSE 8000

CMD ["./main"]