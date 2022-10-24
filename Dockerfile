FROM golang:1.19.2 AS builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/server/main.go

FROM alpine:latest
COPY --from=builder /app .
CMD ["./app"]