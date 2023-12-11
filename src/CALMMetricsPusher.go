package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

// TODO o co chodzi z gatherer i collector itd., raczej użyj Add()
// tutaj chyba jakiś przykład jak to działa https://pkg.go.dev/github.com/prometheus/client_golang/prometheus/push#example-Pusher.Add

// TODO packet loss percentage
type CALMMetricsPusher struct {
	gatewayPusher    *push.Pusher
	avgRTT           prometheus.Gauge
	maxRTT           prometheus.Gauge
	avgOneWayLatency prometheus.Gauge
}

func NewCALMMetricsPusher(pushGatewayURL string, job string) *CALMMetricsPusher {
	calmPusher := CALMMetricsPusher{
		gatewayPusher: push.New(pushGatewayURL, job),
		avgRTT: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "calm_average_RTT",
			Help: "Average Round Trip Time measured between particular nodes between clusters",
		}),
		maxRTT: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "calm_max_node_to_node_latency",
			Help: "Maximum latency measured between particular nodes between clusters",
		}),
		avgOneWayLatency: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "calm_avg_one_way_latency_client_to_server",
			Help: "Average one way latency between particular nodes between clusters measured from client to server",
		}),
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(calmPusher.avgOneWayLatency, calmPusher.avgRTT, calmPusher.maxRTT)
	calmPusher.gatewayPusher.Gatherer(registry)

	return &calmPusher
}

func (m *CALMMetricsPusher) addMetrics() {
	// tutaj w grouping labelki, TODO uzmiennic
	m.gatewayPusher.
		Grouping("measurement_ID", "DRUGI POMIAR").
		Grouping("source_node", "client-side-worker").
		Grouping("target_node", "server-side-worker").
		Grouping("itd", "itp")

	err := m.gatewayPusher.Add()
	if err != nil {
		fmt.Println("Could not push to Pushgateway:", err)
	}
}

// var avgRTT = prometheus.NewGauge(prometheus.GaugeOpts{
//	Name: "calm_average_RTT",
//	Help: "Average Round Trip Time measured between particular nodes between clusters",
//})
//
//var avgOneWayLatency = prometheus.NewGauge(prometheus.GaugeOpts{
//	Name: "calm_avg_one_way_latency_client_to_server",
//	Help: "Average one way latency between particular nodes between clusters measured from client to server",
//})
//
//var maxRTT = prometheus.NewGauge(prometheus.GaugeOpts{
//	Name: "calm_max_node_to_node_latency",
//	Help: "Maximum latency measured between particular nodes between clusters",
//})

func (c *CALMMetricsPusher) SetAvgRTTValue(value float64) {
	c.avgRTT.Set(value)
}

func (c *CALMMetricsPusher) SetMaxRTT(value float64) {
	c.maxRTT.Set(value)
}

func (c *CALMMetricsPusher) SetAvgOneWayLatency(value float64) {
	c.avgOneWayLatency.Set(value)
}
