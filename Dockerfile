FROM golang:1.17-alpine as builder

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go get -d -v ./...

COPY . .

RUN go build -o build

FROM alpine

WORKDIR /app

COPY --from=builder /go/src/app/build .

CMD ["/app/build"]
