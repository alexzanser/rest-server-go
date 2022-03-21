FROM golang:latest as builder

RUN mkdir /build
WORKDIR /build
ADD . /build

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o server cmd/api/main.go

FROM scratch

ENV SERVERPORT=4112 
EXPOSE 4112

COPY --from=builder /build/server .

CMD ["./server"]