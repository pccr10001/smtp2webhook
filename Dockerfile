FROM golang:1.19.1-alpine AS builder
LABEL maintainer="pccr10001@gmail.com"
WORKDIR /build
ADD . /build
RUN apk add build-base
RUN go build -o smtp2webhook

FROM alpine
WORKDIR /app
COPY --from=builder /build/smtp2webhook /app/
RUN chmod +x /app/smtp2webhook
EXPOSE 2525
ENTRYPOINT ["/app/smtp2webhook"]