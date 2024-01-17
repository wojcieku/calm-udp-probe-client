package main

import (
	"bytes"
	"github.com/montanaflynn/stats"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var packetsSentCounter int
var rawMetrics []RawMetric

type RawMetric struct {
	clientSendStamp    time.Time
	serverReceiveStamp time.Time
	clientReceiveStamp time.Time
}

type CALMMetric struct {
	avgRTT               float64
	maxRTT               float64
	percentile95thRTT    float64
	packetLossPercentage float64
}

func runMeasurement(remoteIP *string, remotePort *int, interval *int, duration *int) CALMMetric {
	conn := prepareConnection(remoteIP, remotePort)
	// close the connection
	defer conn.Close()

	ticker := time.NewTicker(time.Duration(*interval) * time.Millisecond)
	defer ticker.Stop()
	timer := time.After(time.Duration(*duration) * time.Second)

	for {
		select {
		case <-ticker.C:
			measure(conn)
		case <-timer:
			log.Info("Time's up, measurement finished")
			return processRawMetrics()
		}
	}
}

func prepareConnection(remoteIP *string, remotePort *int) *net.UDPConn {
	udpServer, err := net.ResolveUDPAddr("udp", *remoteIP+":"+strconv.Itoa(*remotePort))
	if err != nil {
		log.Error("ResolveUDPAddr failed:", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, udpServer)
	if err != nil {
		log.Error("Listen failed:", err.Error())
		os.Exit(1)
	}
	return conn
}

func measure(conn *net.UDPConn) {
	send(conn)
	result, err := handleResponse(conn)
	if err == nil {
		rawMetrics = append(rawMetrics, result)
	}
}

func send(conn *net.UDPConn) {
	conn.SetReadDeadline(time.Now().Add(time.Second))
	_, err := conn.Write([]byte(time.Now().Format(time.StampMilli)))
	packetsSentCounter++

	if err != nil {
		log.Error("Write data failed:", err.Error())
		os.Exit(1)
	}
}

func handleResponse(conn *net.UDPConn) (RawMetric, error) {
	result := RawMetric{}
	received := make([]byte, 50)
	_, err := conn.Read(received)
	if err != nil {
		log.Error("Read data failed:", err.Error())
		if strings.Contains(err.Error(), "timeout") {
			return result, err
		}
		os.Exit(1)
	}
	clientReceivedStamp, err := time.Parse(time.StampMilli, time.Now().Format(time.StampMilli))
	if err != nil {
		log.Error("received client stamp parsing error", err.Error())
	}
	received = bytes.Trim(received, "\x00")
	receivedString := string(received)

	stamps := strings.Split(receivedString, ";")

	clientSendStamp, err := time.Parse(time.StampMilli, stamps[0])
	if err != nil {
		log.Error("client stamp parsing error", err.Error())
	}
	serverStamp, err := time.Parse(time.StampMilli, stamps[1])
	if err != nil {
		log.Error("server stamp parsing error", err)
	}

	result.clientSendStamp = clientSendStamp
	result.serverReceiveStamp = serverStamp
	result.clientReceiveStamp = clientReceivedStamp
	return result, nil
}

func processRawMetrics() CALMMetric {
	calmMetric := CALMMetric{}
	var rtts []float64
	var clientToServerDurations []float64
	var serverToClientDurations []float64

	// calculate latencies
	for _, r := range rawMetrics {
		rtt := r.clientReceiveStamp.Sub(r.clientSendStamp)
		rtts = append(rtts, float64(rtt.Milliseconds()))

		clientToServerDuration := r.serverReceiveStamp.Sub(r.clientSendStamp)
		clientToServerDurations = append(clientToServerDurations, float64(clientToServerDuration.Milliseconds()))

		serverToClientDuration := r.clientReceiveStamp.Sub(r.serverReceiveStamp)
		serverToClientDurations = append(serverToClientDurations, float64(serverToClientDuration.Milliseconds()))
	}

	avgRTT, err := stats.Mean(rtts)
	if err != nil {
		log.Error("Failed to calculate RTT mean:", err.Error())
	}
	avgRTT = truncateResult(avgRTT)

	maxRTT, err := stats.Max(rtts)
	if err != nil {
		log.Error("Failed to calculate max RTT:", err.Error())
	}
	maxRTT = truncateResult(maxRTT)

	percentile95thRTT, err := stats.Percentile(rtts, 95)
	if err != nil {
		log.Error("Failed to calculate 95th percentile of RTT:", err.Error())
	}
	truncateResult(percentile95thRTT)

	calmMetric.avgRTT = avgRTT
	calmMetric.maxRTT = maxRTT
	calmMetric.percentile95thRTT = percentile95thRTT
	calmMetric.packetLossPercentage = float64((cap(rawMetrics)-packetsSentCounter)/packetsSentCounter) * 100

	return calmMetric
}

func truncateResult(result float64) float64 {
	result, _ = decimal.NewFromFloat(result).Truncate(2).Float64()
	return result
}
