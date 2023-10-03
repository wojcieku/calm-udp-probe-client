package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "1501"
	SERVER_TYPE = "udp"
)

func main() {
	for i := 0; i < 100; i++ {
		/*
				tu na przykład można sparametryzować liczbę wywołań (nieskończona np.), albo można wjebać sleepa żeby był interval
			do interwału: https://stackoverflow.com/questions/16466320/is-there-a-way-to-do-repetitive-tasks-at-intervals
		*/
		udpServer, err := net.ResolveUDPAddr("udp", "127.0.0.1:1501") //ale co to robi? xd
		if err != nil {
			println("ResolveUDPAddr failed:", err.Error())
			os.Exit(1)
		}

		conn, err := net.DialUDP("udp", nil, udpServer)
		if err != nil {
			println("Listen failed:", err.Error())
			os.Exit(1)
		}

		//close the connection
		defer conn.Close()

		_, err = conn.Write([]byte(time.Now().Format(time.StampMilli)))
		if err != nil {
			println("Write data failed:", err.Error())
			os.Exit(1)
		}

		// buffer to get data

		received := make([]byte, 3000)
		_, err = conn.Read(received)
		if err != nil {
			println("Read data failed:", err.Error())
			os.Exit(1)
		}
		received = bytes.Trim(received, "\x00")
		receivedString := string(received)
		stamps := strings.Split(receivedString, ";")
		//print(stamps[0])
		//print(stamps[1])
		serverStamp, err := time.Parse(time.StampMilli, stamps[0])
		if err != nil {
			println("server stamp parsing error", err)
		}
		clientSendStamp, err := time.Parse(time.StampMilli, stamps[1])
		if err != nil {
			//println("client stamp parsing error", err.Error())
			fmt.Println(err)
		}

		duration := serverStamp.Sub(clientSendStamp)
		println(duration.Milliseconds())
	}

}
