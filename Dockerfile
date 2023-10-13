FROM --platform=linux/amd64 golang:1.20-buster

RUN mkdir /app
WORKDIR /app

COPY ./setup /app
COPY ./server /app/server
COPY ./config /app/config

COPY ./migrations /app/migrations