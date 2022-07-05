FROM golang:1.17-alpine as build-stage

RUN mkdir -p /app

WORKDIR /app

COPY . /app
RUN go mod download

RUN go build -o balance main.go


EXPOSE 9999

ENTRYPOINT [ "/app/balance" ]