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
ENV METRICS_AGGREGATOR="localhost:9091"
ENV MEASUREMENT_ID = ""
ENV SRC_NODE = ""
ENV TARGET_NODE = ""
ENV SRC_CLUSTER = ""
ENV TARGET_CLUSTER = ""

ENTRYPOINT ["/bin/sh", "-c","./latencyClient -addr $ADDRESS -port $PORT -interval $INTERVAL -duration $DURATION \
 -aggr $METRICS_AGGREGATOR -id $MEASUREMENT_ID -srcNode $SRC_NODE -targetNode $TARGET_NODE -srcCluster $SRC_CLUSTER \
 -targetCluster $TARGET_CLUSTER"]