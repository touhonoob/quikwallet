FROM golang:1.16

# cache go modules
RUN mkdir -p /go/src/github.com/touhonoob/quikwallet/
WORKDIR /go/src/github.com/touhonoob/quikwallet/
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# build
COPY . .
RUN go mod vendor -v
RUN go build -v -o /go/bin/quikwallet-migrate -mod=vendor cmd/migrate.go
RUN go build -v -o /go/bin/quikwallet-seeder -mod=vendor cmd/seeder.go
RUN go build -v -o /go/bin/quikwallet-queue-worker -mod=vendor cmd/queue_worker.go
RUN go build -v -o /go/bin/quikwallet -mod=vendor cmd/main.go