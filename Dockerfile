FROM golang:1.21

WORKDIR /root

COPY go.mod ./

RUN go mod download

COPY . ./

RUN go build -o main

EXPOSE 8080

CMD ["./main"]