package main

import "flag"

func main() {
	//execution parameters
	remoteIp := flag.String("addr", "localhost", "Server address")
	remotePort := flag.Int("port", 1501, "The port that server is listening on")
	interval := flag.Int("interval", 200, "Time interval between measurement packets in ms")
	duration := flag.Int("duration", 5, "Overall measurement duration in seconds")

	//packetSize?
	flag.Parse()

	runMeasurement(remoteIp, remotePort, interval, duration)

}
