FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY config.yaml ./
COPY .env ./

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

EXPOSE 8080

CMD ["./main"]

FROM golang:1.25 AS debug

WORKDIR /app

RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY go.mod go.sum ./
RUN go mod download

COPY --from=builder /app/main .
COPY config.yaml ./config.yaml
COPY .env ./.env

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -gcflags="all=-N -l" -o main ./cmd/app

EXPOSE 8080 40000

CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./main"]
