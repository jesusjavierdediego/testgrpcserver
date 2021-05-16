FROM golang:1.16.4-alpine3.13 as build-env

MAINTAINER mahendrabagul <bagulm123@gmail.com>

ENV GO111MODULE=on

WORKDIR /app

RUN apk add --update --no-cache ca-certificates git

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/golang-grpc-server server/main.go

FROM scratch
COPY --from=build-env /go/bin/golang-grpc-server /app/golang-grpc-server
EXPOSE 50051
ENTRYPOINT ["/app/golang-grpc-server"]
