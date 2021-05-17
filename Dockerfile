#FROM golang:1.14 as build
#
#WORKDIR /go/src/github.com/observerward
#
#COPY . .
#
#RUN go build -o /go/bin/observerward observer_ward.go

FROM nvidia/cuda:11.1-base

COPY bin/observerward /usr/bin/observerward

CMD [ "observerward" ]

