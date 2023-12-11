package main

func main() {
	//// execution parameters
	//remoteIP := flag.String("addr", "localhost", "Server address")
	//remotePort := flag.Int("port", 1501, "The port that server is listening on")
	//interval := flag.Int("interval", 200, "Time interval between measurement packets in ms")
	//duration := flag.Int("duration", 5, "Overall measurement duration in seconds")
	//
	//// packetSize?
	//flag.Parse()
	//
	//runMeasurement(remoteIP, remotePort, interval, duration)

	metricsPusher := NewCALMMetricsPusher("localhost:9091", "myJob")
	metricsPusher.SetAvgRTTValue(20)
	metricsPusher.SetAvgOneWayLatency(233)
	metricsPusher.SetMaxRTT(699)

	metricsPusher.addMetrics()

}
