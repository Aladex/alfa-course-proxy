FROM golang:1.15-alpine

RUN mkdir /app

COPY main.go /app/.

RUN cd /app && go build -o app && rm main.go

WORKDIR /app

EXPOSE 8090