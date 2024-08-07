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
	queueSubsystem = "queue"
)

func init() {
	registerCollector(queueSubsystem, defaultEnabled, NewQueueCollector)
}

func (c *queueCollector) Update(ch chan<- prometheus.Metric) error {
	var metricType prometheus.ValueType
	return nil
}

type queueCollector struct {
	logger log.Logger
}

func NewQueueCollector(logger log.Logger) (Collector, error) {
	return &queueCollector{}, nil
}
