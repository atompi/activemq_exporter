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

	"github.com/atompi/activemq_exporter/collector/utils"
	"github.com/atompi/activemq_exporter/options"
	"github.com/go-kit/log"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	brokerSubsystem = "broker"
)

func init() {
	registerCollector(brokerSubsystem, defaultEnabled, NewBrokerCollector)
}

type TopicProducer struct {
	ObjectName string `json:"objectName"`
}

type QueueProducer struct {
	ObjectName string `json:"objectName"`
}

type BrokerTopic struct {
	ObjectName string `json:"objectName"`
}

type BrokerQueue struct {
	ObjectName string `json:"objectName"`
}

type TopicSubscriber struct {
	ObjectName string `json:"objectName"`
}

type DurableTopicSubscriber struct {
	ObjectName string `json:"objectName"`
}

type InactiveDurableTopicSubscriber struct {
	ObjectName string `json:"objectName"`
}

type QueueSubscriber struct {
	ObjectName string `json:"objectName"`
}

type BrokerValue struct {
	MemoryLimit                     int                              `json:"MemoryLimit"`
	MemoryPercentUsage              int                              `json:"MemoryPercentUsage"`
	StoreLimit                      int                              `json:"StoreLimit"`
	StorePercentUsage               int                              `json:"StorePercentUsage"`
	TempLimit                       int                              `json:"TempLimit"`
	TempPercentUsage                int                              `json:"TempPercentUsage"`
	CurrentConnectionsCount         int                              `json:"CurrentConnectionsCount"`
	TotalProducerCount              int                              `json:"TotalProducerCount"`
	TotalConsumerCount              int                              `json:"TotalConsumerCount"`
	TopicProducers                  []TopicProducer                  `json:"TopicProducers"`
	QueueProducers                  []QueueProducer                  `json:"QueueProducers"`
	Topics                          []BrokerTopic                    `json:"Topics"`
	Queues                          []BrokerQueue                    `json:"Queues"`
	TopicSubscribers                []TopicSubscriber                `json:"TopicSubscribers"`
	DurableTopicSubscribers         []DurableTopicSubscriber         `json:"DurableTopicSubscribers"`
	InactiveDurableTopicSubscribers []InactiveDurableTopicSubscriber `json:"InactiveDurableTopicSubscribers"`
	QueueSubscribers                []QueueSubscriber                `json:"QueueSubscribers"`
}

type Broker struct {
	Status    int         `json:"status"`
	Timestamp int         `json:"timestamp"`
	Value     BrokerValue `json:"value"`
}

type brokerCollector struct {
	status                              typedDesc
	timestamp                           typedDesc
	memoryLimit                         typedDesc
	memoryPercentUsage                  typedDesc
	storeLimit                          typedDesc
	storePercentUsage                   typedDesc
	tempLimit                           typedDesc
	tempPercentUsage                    typedDesc
	currentConnectionsCount             typedDesc
	totalProducerCount                  typedDesc
	totalConsumerCount                  typedDesc
	topicProducerCount                  typedDesc
	queueProducerCount                  typedDesc
	topicCount                          typedDesc
	queueCount                          typedDesc
	topicSubscriberCount                typedDesc
	durableTopicSubscriberCount         typedDesc
	inactiveDurableTopicSubscriberCount typedDesc
	queueSubscriberCount                typedDesc
	logger                              log.Logger
}

func NewBrokerCollector(logger log.Logger) (Collector, error) {
	return &brokerCollector{
		status: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "status"),
			"Status of the broker.",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		timestamp: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "timestamp"),
			"System time in seconds since epoch (1970).",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		memoryLimit: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "memory_limit"),
			"Memory limit in bytes.",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		memoryPercentUsage: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "memory_percent_usage"),
			"Memory usage percentage.",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		storeLimit: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "store_limit"),
			"Store limit in bytes.",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		storePercentUsage: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "store_percent_usage"),
			"Store usage percentage.",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		tempLimit: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "temp_limit"),
			"Temp limit in bytes.",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		tempPercentUsage: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "temp_percent_usage"),
			"Temp usage percentage.",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		currentConnectionsCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "current_connections_count"),
			"Current connections count.",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		totalProducerCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "total_producer_count"),
			"Total producer count.",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		totalConsumerCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "total_consumer_count"),
			"Total consumer count.",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		topicProducerCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "topic_producer_count"),
			"Topic producer count.",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		queueProducerCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "queue_producer_count"),
			"Queue producer count.",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		topicCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "topic_count"),
			"Topic count.",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		queueCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "queue_count"),
			"Queue count.",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		topicSubscriberCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "topic_subscriber_count"),
			"Topic subscriber count.",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		durableTopicSubscriberCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "durable_topic_subscriber_count"),
			"Durable topic subscriber count.",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		inactiveDurableTopicSubscriberCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "inactive_durable_topic_subscriber_count"),
			"Inactive durable topic subscriber count.",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		queueSubscriberCount: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "queue_subscriber_count"),
			"Queue subscriber count.",
			[]string{"broker"}, nil,
		), prometheus.GaugeValue},
		logger: logger,
	}, nil
}

func (c *brokerCollector) Update(ch chan<- prometheus.Metric) error {
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

	status := brokerData.Status
	timestamp := brokerData.Timestamp
	memoryLimit := brokerData.Value.MemoryLimit
	memoryPercentUsage := brokerData.Value.MemoryPercentUsage
	storeLimit := brokerData.Value.StoreLimit
	storePercentUsage := brokerData.Value.StorePercentUsage
	tempLimit := brokerData.Value.TempLimit
	tempPercentUsage := brokerData.Value.TempPercentUsage
	currentConnectionsCount := brokerData.Value.CurrentConnectionsCount
	totalProducerCount := brokerData.Value.TotalProducerCount
	totalConsumerCount := brokerData.Value.TotalConsumerCount
	topicProducerCount := len(brokerData.Value.TopicProducers)
	queueProducerCount := len(brokerData.Value.QueueProducers)
	topicCount := len(brokerData.Value.Topics)
	queueCount := len(brokerData.Value.Queues)
	topicSubscriberCount := len(brokerData.Value.TopicSubscribers)
	durableTopicSubscriberCount := len(brokerData.Value.DurableTopicSubscribers)
	inactiveDurableTopicSubscriberCount := len(brokerData.Value.InactiveDurableTopicSubscribers)
	queueSubscriberCount := len(brokerData.Value.QueueSubscribers)
	ch <- c.status.mustNewConstMetric(float64(status), opts.Broker.Name)
	ch <- c.timestamp.mustNewConstMetric(float64(timestamp), opts.Broker.Name)
	ch <- c.memoryLimit.mustNewConstMetric(float64(memoryLimit), opts.Broker.Name)
	ch <- c.memoryPercentUsage.mustNewConstMetric(float64(memoryPercentUsage), opts.Broker.Name)
	ch <- c.storeLimit.mustNewConstMetric(float64(storeLimit), opts.Broker.Name)
	ch <- c.storePercentUsage.mustNewConstMetric(float64(storePercentUsage), opts.Broker.Name)
	ch <- c.tempLimit.mustNewConstMetric(float64(tempLimit), opts.Broker.Name)
	ch <- c.tempPercentUsage.mustNewConstMetric(float64(tempPercentUsage), opts.Broker.Name)
	ch <- c.currentConnectionsCount.mustNewConstMetric(float64(currentConnectionsCount), opts.Broker.Name)
	ch <- c.totalProducerCount.mustNewConstMetric(float64(totalProducerCount), opts.Broker.Name)
	ch <- c.totalConsumerCount.mustNewConstMetric(float64(totalConsumerCount), opts.Broker.Name)
	ch <- c.topicProducerCount.mustNewConstMetric(float64(topicProducerCount), opts.Broker.Name)
	ch <- c.queueProducerCount.mustNewConstMetric(float64(queueProducerCount), opts.Broker.Name)
	ch <- c.topicCount.mustNewConstMetric(float64(topicCount), opts.Broker.Name)
	ch <- c.queueCount.mustNewConstMetric(float64(queueCount), opts.Broker.Name)
	ch <- c.topicSubscriberCount.mustNewConstMetric(float64(topicSubscriberCount), opts.Broker.Name)
	ch <- c.durableTopicSubscriberCount.mustNewConstMetric(float64(durableTopicSubscriberCount), opts.Broker.Name)
	ch <- c.inactiveDurableTopicSubscriberCount.mustNewConstMetric(float64(inactiveDurableTopicSubscriberCount), opts.Broker.Name)
	ch <- c.queueSubscriberCount.mustNewConstMetric(float64(queueSubscriberCount), opts.Broker.Name)
	return nil
}
