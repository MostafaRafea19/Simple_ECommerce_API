FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY ../../go.mod ../../go.sum ./

RUN go mod download

COPY ../../ .

RUN CGO_ENABLED=0 GOOS=linux go build -o customers-app main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/customers-app .

EXPOSE 8080

CMD ["./customers-app"]
