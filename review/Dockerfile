# build stage
FROM golang:latest AS builder

RUN apt-get update && apt-get install -y git gcc libc-dev

WORKDIR /go/src/coupon-service

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o coupon_service ./cmd/coupon_service/  

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/coupon-service/coupon_service .
RUN chmod +x /app/coupon_service
ENTRYPOINT [ "/app/coupon_service" ]
# CMD [ "/app/main" ]
EXPOSE 8080
