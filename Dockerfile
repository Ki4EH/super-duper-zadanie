FROM golang AS builder

WORKDIR /app

COPY service/go.mod service/go.sum ./

ENV GOPROXY=direct

RUN go mod download

COPY service/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd

FROM alpine:latest

RUN apk add --no-cache build-base

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

CMD ["./main"]
