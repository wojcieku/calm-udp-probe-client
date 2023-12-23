package main

import "flag"

func main() {
	// execution parameters
	remoteIP := flag.String("addr", "localhost", "Server address")
	remotePort := flag.Int("port", 1501, "The port that server is listening on")
	interval := flag.Int("interval", 200, "Time interval between measurement packets in ms")
	duration := flag.Int("duration", 5, "Overall measurement duration in seconds")
	metricsAggregator := flag.String("aggr", "localhost:9091", "Metrics aggregator URL")
	measurementID := flag.String("id", "", "Measurement ID")
	sourceNode := flag.String("srcNode", "", "Source (client) node name")
	targetNode := flag.String("targetNode", "", "Target (server) node name")
	sourceCluster := flag.String("srcCluster", "", "Source (client) cluster name")
	targetCluster := flag.String("targetCluster", "", "Target (server) cluster name")
	flag.Parse()

	calmMetric := runMeasurement(remoteIP, remotePort, interval, duration)

	metricsPusher := NewCALMMetricsPusher(*metricsAggregator, "calm-push-gateway")
	setMetricsValues(metricsPusher, calmMetric)
	metricsPusher.addMetrics(*measurementID, *sourceNode, *targetNode, *sourceCluster, *targetCluster)
}

func setMetricsValues(metricsPusher *CALMMetricsPusher, calmMetric CALMMetric) {
	metricsPusher.SetAvgRTTValue(calmMetric.avgRTT)
	metricsPusher.SetMaxRTTValue(calmMetric.maxRTT)
	metricsPusher.SetPercentile95thRTTValue(calmMetric.percentile95thRTT)
	metricsPusher.SetPacketLossPercentageValue(calmMetric.packetLossPercentage)
}
