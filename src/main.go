package main

import "flag"

func main() {
	// execution parameters
	remoteIP := flag.String("addr", "localhost", "Server address")
	remotePort := flag.Int("port", 1501, "The port that server is listening on")
	interval := flag.Int("interval", 200, "Time interval between measurement packets in ms")
	duration := flag.Int("duration", 5, "Overall measurement duration in seconds")
	metricsAggregator := flag.String("aggr", "localhost:9091", "Metrics aggregator URL")
	// TODO: server node, client node names; measurement ID,
	flag.Parse()

	calmMetric := runMeasurement(remoteIP, remotePort, interval, duration)

	metricsPusher := NewCALMMetricsPusher(*metricsAggregator, "calm-push-gateway")
	metricsPusher.SetAvgRTTValue(calmMetric.avgRTT)
	metricsPusher.SetAvgClientToServerLatencyValue(calmMetric.avgClientToServerLatency)
	metricsPusher.SetMaxRTTValue(calmMetric.maxRTT)
	metricsPusher.addMetrics()

}
