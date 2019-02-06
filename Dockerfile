FROM golang:alpine as build
ADD ./go-cron.go /tmp/goapp/
WORKDIR /tmp/goapp
RUN apk add --no-cache git \
&&	go get github.com/robfig/cron \
&&	go build -o ./multi-cron \
&&	apk del git


FROM alpine:latest
COPY --from=build /tmp/goapp/multi-cron /usr/local/bin/multi-cron
ENTRYPOINT ["multi-cron"]
