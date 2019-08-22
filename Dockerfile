FROM golang:1.12.5-alpine as builder

ARG COMMIT_REF
ARG BUILD_DATE

ENV APP_COMMIT_REF=${COMMIT_REF} \
    APP_BUILD_DATE=${BUILD_DATE}

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh g++ glide ca-certificates

RUN adduser -D -g '' appuser

RUN mkdir -p /go/src/github.com/brunoksato/golang-boilerplate
ADD . /go/src/github.com/brunoksato/golang-boilerplate
WORKDIR /go/src/github.com/brunoksato/golang-boilerplate
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o golang-boilerplate .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/src/github.com/brunoksato/golang-boilerplate/golang-boilerplate /app/
WORKDIR /app
EXPOSE 8080
CMD ["./golang-boilerplate"]
