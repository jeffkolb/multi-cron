FROM golang:latest
ADD ./go-cron.go /tmp/goapp
WORKDIR /tmp/goapp
RUN go get github.com/robfig/cron \
&&	go build -o /multi-cron \
&&	cp multi-cron /usr/local/bin/test-cron
ENTRYPOINT {"/bin/bash"]
