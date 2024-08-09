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
	now    typedDesc
	logger log.Logger
}

func NewTopicCollector(logger log.Logger) (Collector, error) {
	return &topicCollector{
		now: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, topicSubsystem, "now"),
			"System time in seconds since epoch (1970).",
			nil, nil,
		), prometheus.GaugeValue},
		logger: logger,
	}, nil
}

func (c *topicCollector) Update(ch chan<- prometheus.Metric) error {
	now := broker.Timestamp
	ch <- c.now.mustNewConstMetric(float64(now))
	return nil
}
