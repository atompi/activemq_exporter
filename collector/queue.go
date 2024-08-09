/*
Copyright © 2024 Atom Pi <coder.atompi@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package collector

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/atompi/activemq_exporter/collector/utils"
	"github.com/atompi/activemq_exporter/options"
	"github.com/go-kit/log"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	queueSubsystem = "queue"
)

func init() {
	registerCollector(queueSubsystem, defaultDisabled, NewQueueCollector)
}

type QueueValue struct {
	MemoryUsageByteCount int    `json:"MemoryUsageByteCount"`
	MemoryPercentUsage   int    `json:"MemoryPercentUsage"`
	InFlightCount        int    `json:"InFlightCount"`
	ForwardCount         int    `json:"ForwardCount"`
	Name                 string `json:"Name"`
	EnqueueCount         int    `json:"EnqueueCount"`
	ConsumerCount        int    `json:"ConsumerCount"`
	MemoryLimit          int    `json:"MemoryLimit"`
	DequeueCount         int    `json:"DequeueCount"`
	ProducerCount        int    `json:"ProducerCount"`
}

type Queue struct {
	Status    int        `json:"status"`
	Timestamp int        `json:"timestamp"`
	Value     QueueValue `json:"value"`
}

type queueCollector struct {
	status               typedDesc
	timestamp            typedDesc
	memoryUsageByteCount typedDesc
	memoryPercentUsage   typedDesc
	inFlightCount        typedDesc
	forwardCount         typedDesc
	dequeueCount         typedDesc
	enqueueCount         typedDesc
	consumerCount        typedDesc
	memoryLimit          typedDesc
	producerCount        typedDesc
	logger               log.Logger
}

func NewQueueCollector(logger log.Logger) (Collector, error) {
	return &queueCollector{
		status: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, queueSubsystem, "status"),
			"The status of the queue.",
			[]string{"queue"}, nil,
		), prometheus.GaugeValue},
		timestamp: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, queueSubsystem, "timestamp"),
			"System time in seconds since epoch (1970).",
			[]string{"queue"}, nil,
		), prometheus.GaugeValue},
		memoryUsageByteCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, queueSubsystem, "memory_usage_byte_count"),
			"The memory usage byte count of the queue.",
			[]string{"queue"}, nil,
		), prometheus.GaugeValue},
		memoryPercentUsage: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, queueSubsystem, "memory_percent_usage"),
			"The memory percent usage of the queue.",
			[]string{"queue"}, nil,
		), prometheus.GaugeValue},
		inFlightCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, queueSubsystem, "in_flight_count"),
			"The in flight count of the queue.",
			[]string{"queue"}, nil,
		), prometheus.GaugeValue},
		forwardCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, queueSubsystem, "forward_count"),
			"The forward count of the queue.",
			[]string{"queue"}, nil,
		), prometheus.GaugeValue},
		dequeueCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, queueSubsystem, "dequeue_count"),
			"The dequeue count of the queue.",
			[]string{"queue"}, nil,
		), prometheus.GaugeValue},
		enqueueCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, queueSubsystem, "enqueue_count"),
			"The enqueue count of the queue.",
			[]string{"queue"}, nil,
		), prometheus.GaugeValue},
		consumerCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, queueSubsystem, "consumer_count"),
			"The consumer count of the queue.",
			[]string{"queue"}, nil,
		), prometheus.GaugeValue},
		memoryLimit: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, queueSubsystem, "memory_limit"),
			"The memory limit of the queue.",
			[]string{"queue"}, nil,
		), prometheus.GaugeValue},
		producerCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, queueSubsystem, "producer_count"),
			"The producer count of the queue.",
			[]string{"queue"}, nil,
		), prometheus.GaugeValue},
		logger: logger,
	}, nil
}

func (c *queueCollector) Update(ch chan<- prometheus.Metric) error {
	opts := options.Opts

	url := fmt.Sprintf("%s/org.apache.activemq:type=Broker,brokerName=%s", opts.Jolokia.URL, opts.Broker.Name)
	data := []byte{}
	brokerData := &Broker{}
	err := utils.FetchData(url, &data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, brokerData)
	if err != nil {
		return err
	}

	for _, queueObj := range brokerData.Value.Queues {
		if strings.Contains(queueObj.ObjectName, "destinationName=ActiveMQ.") {
			continue
		}

		url := fmt.Sprintf("%s/%s", opts.Jolokia.URL, queueObj.ObjectName)
		data := []byte{}
		queueData := &Queue{}
		err := utils.FetchData(url, &data)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, queueData)
		if err != nil {
			return err
		}

		status := queueData.Status
		timestamp := queueData.Timestamp
		memoryUsageByteCount := queueData.Value.MemoryUsageByteCount
		memoryPercentUsage := queueData.Value.MemoryPercentUsage
		inFlightCount := queueData.Value.InFlightCount
		forwardCount := queueData.Value.ForwardCount
		dequeueCount := queueData.Value.DequeueCount
		enqueueCount := queueData.Value.EnqueueCount
		consumerCount := queueData.Value.ConsumerCount
		memoryLimit := queueData.Value.MemoryLimit
		producerCount := queueData.Value.ProducerCount
		ch <- c.status.mustNewConstMetric(float64(status), queueData.Value.Name)
		ch <- c.timestamp.mustNewConstMetric(float64(timestamp), queueData.Value.Name)
		ch <- c.memoryUsageByteCount.mustNewConstMetric(float64(memoryUsageByteCount), queueData.Value.Name)
		ch <- c.memoryPercentUsage.mustNewConstMetric(float64(memoryPercentUsage), queueData.Value.Name)
		ch <- c.inFlightCount.mustNewConstMetric(float64(inFlightCount), queueData.Value.Name)
		ch <- c.forwardCount.mustNewConstMetric(float64(forwardCount), queueData.Value.Name)
		ch <- c.dequeueCount.mustNewConstMetric(float64(dequeueCount), queueData.Value.Name)
		ch <- c.enqueueCount.mustNewConstMetric(float64(enqueueCount), queueData.Value.Name)
		ch <- c.consumerCount.mustNewConstMetric(float64(consumerCount), queueData.Value.Name)
		ch <- c.memoryLimit.mustNewConstMetric(float64(memoryLimit), queueData.Value.Name)
		ch <- c.producerCount.mustNewConstMetric(float64(producerCount), queueData.Value.Name)
	}
	return nil
}
