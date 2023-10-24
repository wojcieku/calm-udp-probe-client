FROM golang:1.20 AS build-stage

WORKDIR /probeClient
COPY go.mod /probeClient
RUN go mod download

COPY /src/*.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /latencyClient

FROM alpine:3.14 AS build-release-stage

WORKDIR /

COPY --from=build-stage /latencyClient /latencyClient
RUN apk add -U tzdata
ENV TZ=Europe/Sarajevo
RUN cp /usr/share/zoneinfo/Europe/Sarajevo /etc/localtime
ENV ADDRESS="localhost"
ENV PORT=1501
ENV INTERVAL=100
ENV DURATION=10

ENTRYPOINT ["/bin/sh", "-c","./latencyClient -addr $ADDRESS -port $PORT -interval $INTERVAL -duration $DURATION"]