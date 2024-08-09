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
	topicSubsystem = "topic"
)

func init() {
	registerCollector(topicSubsystem, defaultDisabled, NewTopicCollector)
}

type TopicValue struct {
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

type Topic struct {
	Status    int        `json:"status"`
	Timestamp int        `json:"timestamp"`
	Value     TopicValue `json:"value"`
}

type topicCollector struct {
	status               typedDesc
	timestamp            typedDesc
	memoryUsageByteCount typedDesc
	memoryPercentUsage   typedDesc
	inFlightCount        typedDesc
	forwardCount         typedDesc
	enqueueCount         typedDesc
	consumerCount        typedDesc
	memoryLimit          typedDesc
	dequeueCount         typedDesc
	producerCount        typedDesc
	logger               log.Logger
}

func NewTopicCollector(logger log.Logger) (Collector, error) {
	return &topicCollector{
		status: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, topicSubsystem, "status"),
			"The status of the topic.",
			[]string{"topic"}, nil,
		), prometheus.GaugeValue},
		timestamp: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, topicSubsystem, "timestamp"),
			"System time in seconds since epoch (1970).",
			[]string{"topic"}, nil,
		), prometheus.GaugeValue},
		memoryUsageByteCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, topicSubsystem, "memory_usage_byte_count"),
			"The memory usage of the topic.",
			[]string{"topic"}, nil,
		), prometheus.GaugeValue},
		memoryPercentUsage: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, topicSubsystem, "memory_percent_usage"),
			"The memory usage of the topic.",
			[]string{"topic"}, nil,
		), prometheus.GaugeValue},
		inFlightCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, topicSubsystem, "in_flight_count"),
			"The in flight count of the topic.",
			[]string{"topic"}, nil,
		), prometheus.GaugeValue},
		forwardCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, topicSubsystem, "forward_count"),
			"The forward count of the topic.",
			[]string{"topic"}, nil,
		), prometheus.GaugeValue},
		enqueueCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, topicSubsystem, "enqueue_count"),
			"The enqueue count of the topic.",
			[]string{"topic"}, nil,
		), prometheus.GaugeValue},
		consumerCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, topicSubsystem, "consumer_count"),
			"The consumer count of the topic.",
			[]string{"topic"}, nil,
		), prometheus.GaugeValue},
		memoryLimit: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, topicSubsystem, "memory_limit"),
			"The memory limit of the topic.",
			[]string{"topic"}, nil,
		), prometheus.GaugeValue},
		dequeueCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, topicSubsystem, "dequeue_count"),
			"The dequeue count of the topic.",
			[]string{"topic"}, nil,
		), prometheus.GaugeValue},
		producerCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, topicSubsystem, "producer_count"),
			"The producer count of the topic.",
			[]string{"topic"}, nil,
		), prometheus.GaugeValue},
		logger: logger,
	}, nil
}

func (c *topicCollector) Update(ch chan<- prometheus.Metric) error {
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

	for _, topicObj := range brokerData.Value.Topics {
		if strings.Contains(topicObj.ObjectName, "destinationName=ActiveMQ.") {
			continue
		}

		url := fmt.Sprintf("%s/%s", opts.Jolokia.URL, topicObj.ObjectName)
		data := []byte{}
		topicData := &Topic{}
		err := utils.FetchData(url, &data)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, topicData)
		if err != nil {
			return err
		}

		status := topicData.Status
		timestamp := topicData.Timestamp
		memoryUsageByteCount := topicData.Value.MemoryUsageByteCount
		memoryPercentUsage := topicData.Value.MemoryPercentUsage
		inFlightCount := topicData.Value.InFlightCount
		forwardCount := topicData.Value.ForwardCount
		enqueueCount := topicData.Value.EnqueueCount
		consumerCount := topicData.Value.ConsumerCount
		memoryLimit := topicData.Value.MemoryLimit
		dequeueCount := topicData.Value.DequeueCount
		producerCount := topicData.Value.ProducerCount
		ch <- c.status.mustNewConstMetric(float64(status), topicData.Value.Name)
		ch <- c.timestamp.mustNewConstMetric(float64(timestamp), topicData.Value.Name)
		ch <- c.memoryUsageByteCount.mustNewConstMetric(float64(memoryUsageByteCount), topicData.Value.Name)
		ch <- c.memoryPercentUsage.mustNewConstMetric(float64(memoryPercentUsage), topicData.Value.Name)
		ch <- c.inFlightCount.mustNewConstMetric(float64(inFlightCount), topicData.Value.Name)
		ch <- c.forwardCount.mustNewConstMetric(float64(forwardCount), topicData.Value.Name)
		ch <- c.enqueueCount.mustNewConstMetric(float64(enqueueCount), topicData.Value.Name)
		ch <- c.consumerCount.mustNewConstMetric(float64(consumerCount), topicData.Value.Name)
		ch <- c.memoryLimit.mustNewConstMetric(float64(memoryLimit), topicData.Value.Name)
		ch <- c.dequeueCount.mustNewConstMetric(float64(dequeueCount), topicData.Value.Name)
		ch <- c.producerCount.mustNewConstMetric(float64(producerCount), topicData.Value.Name)
	}
	return nil
}
