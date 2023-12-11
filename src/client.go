package main

import (
	"bytes"
	"github.com/montanaflynn/stats"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// TODO przetwarzanie wynikow
var packetsSentCounter int = 0
var rawMetrics []RawMetric

type RawMetric struct {
	clientSendStamp    time.Time
	serverReceiveStamp time.Time
	clientReceiveStamp time.Time
}

type CALMMetric struct {
	avgRTT                   float64
	maxRTT                   float64
	avgClientToServerLatency float64
	packetLossPercentage     float64
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
			print("Time's up, measurement finished")
			return processRawMetrics()
		}
	}
}

func prepareConnection(remoteIP *string, remotePort *int) *net.UDPConn {
	udpServer, err := net.ResolveUDPAddr("udp", *remoteIP+":"+strconv.Itoa(*remotePort))
	if err != nil {
		println("ResolveUDPAddr failed:", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, udpServer)
	if err != nil {
		println("Listen failed:", err.Error())
		os.Exit(1)
	}
	return conn
}

func measure(conn *net.UDPConn) {
	send(conn)
	result, err := handleResponse(conn)
	if err != nil {
		rawMetrics = append(rawMetrics, result)
	}
}

func send(conn *net.UDPConn) {
	conn.SetReadDeadline(time.Now().Add(time.Second))
	_, err := conn.Write([]byte(time.Now().Format(time.StampMilli)))
	packetsSentCounter++

	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}
}

func handleResponse(conn *net.UDPConn) (RawMetric, error) {
	result := RawMetric{}
	received := make([]byte, 50)
	_, err := conn.Read(received)
	if err != nil {
		println("Read data failed:", err.Error())
		if strings.Contains(err.Error(), "timeout") {
			return result, err
		}
		os.Exit(1)
	}
	clientReceivedStamp := time.Now()
	received = bytes.Trim(received, "\x00")
	receivedString := string(received)

	stamps := strings.Split(receivedString, ";")

	clientSendStamp, err := time.Parse(time.StampMilli, stamps[0])
	if err != nil {
		println("client stamp parsing error", err.Error())
	}
	serverStamp, err := time.Parse(time.StampMilli, stamps[1])
	if err != nil {
		println("server stamp parsing error", err)
	}

	result.clientSendStamp = clientSendStamp
	result.serverReceiveStamp = serverStamp
	result.clientReceiveStamp = clientReceivedStamp
	return result, nil
}

func processRawMetrics() CALMMetric {
	calmMetric := CALMMetric{}
	var clientToServerDurations []float64
	var rtts []float64
	// calculate latencies
	for _, r := range rawMetrics {
		clientToServerDuration := r.serverReceiveStamp.Sub(r.clientSendStamp)
		clientToServerDurations = append(clientToServerDurations, float64(clientToServerDuration.Milliseconds()))

		rtt := r.clientReceiveStamp.Sub(r.clientSendStamp)
		rtts = append(rtts, float64(rtt.Milliseconds()))
	}

	avgClientToServerLatency, err := stats.Mean(clientToServerDurations)
	if err != nil {
		// TODO handle err
	}

	avgRTT, err := stats.Mean(rtts)
	if err != nil {
		// TODO handle err
	}

	maxRTT, err := stats.Max(rtts)
	if err != nil {
		// TODO handle err
	}

	// packet loss percentage
	packetsReceived := cap(rawMetrics)
	calmMetric.packetLossPercentage = float64(packetsReceived/packetsSentCounter) * 100
	calmMetric.avgRTT = avgRTT
	calmMetric.maxRTT = maxRTT
	calmMetric.avgClientToServerLatency = avgClientToServerLatency

	return calmMetric
}
