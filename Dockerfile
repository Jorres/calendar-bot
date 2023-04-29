FROM golang:1.19-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN apk --no-cache add gcc musl-dev
RUN go build -o bot

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/bot /app/bot
COPY --from=builder /app/token.txt /app/token.txt

WORKDIR /app
CMD ["/app/bot"]