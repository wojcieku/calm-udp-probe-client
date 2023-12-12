package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	log "github.com/sirupsen/logrus"
)

type CALMMetricsPusher struct {
	gatewayPusher            *push.Pusher
	avgRTT                   prometheus.Gauge
	maxRTT                   prometheus.Gauge
	percentile95thRTT        prometheus.Gauge
	avgClientToServerLatency prometheus.Gauge
	avgServerToClientLatency prometheus.Gauge
	packetLossPercentage     prometheus.Gauge
}

func NewCALMMetricsPusher(pushGatewayURL string, job string) *CALMMetricsPusher {
	calmPusher := CALMMetricsPusher{
		gatewayPusher: push.New(pushGatewayURL, job),
		avgRTT: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "calm_average_RTT",
			Help: "Average Round Trip Time measured between particular nodes",
		}),
		maxRTT: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "calm_max_RTT",
			Help: "Maximal Round Trip Time measured between particular nodes",
		}),
		percentile95thRTT: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "calm_95th_percentile_RTT",
			Help: "95th percentile of RTT measured between particular nodes",
		}),
		avgClientToServerLatency: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "calm_avg_one_way_latency_client_to_server",
			Help: "Average one way latency between particular nodes measured from client to server",
		}),
		avgServerToClientLatency: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "calm_avg_one_way_latency_server_to_client",
			Help: "Average one way latency between particular nodes measured from server to client",
		}),
		packetLossPercentage: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "calm_packet_loss_percentage",
			Help: "Percentage of lost packet sent",
		}),
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(calmPusher.avgClientToServerLatency, calmPusher.avgRTT, calmPusher.maxRTT)
	calmPusher.gatewayPusher.Gatherer(registry)

	return &calmPusher
}

func (m *CALMMetricsPusher) addMetrics(measurementID string, sourceNode string,
	targetNode string, sourceCluster string, targetCluster string) {
	m.gatewayPusher.
		Grouping("measurement_ID", measurementID).
		Grouping("source_node", sourceNode).
		Grouping("target_node", targetNode).
		Grouping("source_cluster", sourceCluster).
		Grouping("target_cluster", targetCluster)

	err := m.gatewayPusher.Add()
	if err != nil {
		log.Error("Could not push to Pushgateway:", err.Error())
	}
}

func (c *CALMMetricsPusher) SetAvgRTTValue(value float64) {
	c.avgRTT.Set(value)
}

func (c *CALMMetricsPusher) SetMaxRTTValue(value float64) {
	c.maxRTT.Set(value)
}

func (c *CALMMetricsPusher) SetPercentile95thRTTValue(value float64) {
	c.percentile95thRTT.Set(value)
}

func (c *CALMMetricsPusher) SetAvgClientToServerLatencyValue(value float64) {
	c.avgClientToServerLatency.Set(value)
}

func (c *CALMMetricsPusher) SetAvgServerToClientLatencyValue(value float64) {
	c.avgServerToClientLatency.Set(value)
}

func (c *CALMMetricsPusher) SetPacketLossPercentageValue(value float64) {
	c.packetLossPercentage.Set(value)
}
