FROM golang:1.21-alpine3.18 as builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY key.pem .
COPY cert.pem .

EXPOSE 8000
CMD [ "/app/main" ]