package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var packetsSentCounter int = 0 //packet lost percentage?

func runMeasurement(remoteIp *string, remotePort *int, interval *int, duration *int) {
	conn := prepareConnection(remoteIp, remotePort)
	//close the connection
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
			return
		}
	}
}

func prepareConnection(remoteIp *string, remotePort *int) *net.UDPConn {
	udpServer, err := net.ResolveUDPAddr("udp", *remoteIp+":"+strconv.Itoa(*remotePort))
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

// metric format: sentClientTimestamp;serverReceivedTimestamp
func measure(conn *net.UDPConn) {
	send(conn)
	handleResponse(conn)

	//TODO: RTT, loss percentage, aggregation (struct?)

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

func handleResponse(conn *net.UDPConn) {
	received := make([]byte, 50)
	_, err := conn.Read(received)
	if err != nil {
		println("Read data failed:", err.Error())
		if strings.Contains(err.Error(), "timeout") {
			return
		} else {
			os.Exit(1)
		}
	}
	received = bytes.Trim(received, "\x00")
	receivedString := string(received)

	stamps := strings.Split(receivedString, ";")

	clientSendStamp, err := time.Parse(time.StampMilli, stamps[0])
	if err != nil {
		println("client stamp parsing error", err.Error())
		fmt.Println(err)
	}
	serverStamp, err := time.Parse(time.StampMilli, stamps[1])
	if err != nil {
		println("server stamp parsing error", err)
	}

	clientToServerDuration := serverStamp.Sub(clientSendStamp)
	println(clientToServerDuration.Milliseconds())
}
