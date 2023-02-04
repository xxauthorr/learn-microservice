# building the base image

# FROM golang:1.19.3-alpine as builder

# RUN mkdir /app

# COPY . /app

# WORKDIR /app

# RUN CGO_ENABLED=0  go build -o brokerApp ./cmd/api

# RUN chmod +x /app/brokerApp

# build a tiny docker image 

# all the above steps are skipped because
# the makefile is building the buid format of the app

FROM alpine:latest

RUN mkdir /app

COPY brokerApp /app

CMD [ "/app/brokerApp" ]

 
