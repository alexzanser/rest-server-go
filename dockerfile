FROM golang:latest as builder
RUN mkdir /rest_server

ADD . /rest_server/
WORKDIR  /rest_server
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
        go build -o main 

FROM scratch

CMD CMD ["/rest_server/main"]