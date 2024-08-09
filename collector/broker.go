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
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/atompi/activemq_exporter/options"
	"github.com/go-kit/log"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	brokerSubsystem = "broker"
)

var broker = &Broker{}

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
	now    typedDesc
	logger log.Logger
}

func NewBrokerCollector(logger log.Logger) (Collector, error) {
	return &brokerCollector{
		now: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, brokerSubsystem, "now"),
			"System time in seconds since epoch (1970).",
			nil, nil,
		), prometheus.GaugeValue},
		logger: logger,
	}, nil
}

func (c *brokerCollector) Update(ch chan<- prometheus.Metric) error {
	opts := options.Opts

	hc := &http.Client{}
	jolokiaApiPath := "/read/org.apache.activemq:type=Broker,brokerName=" + options.Opts.Broker.Name
	url := opts.Jolokia.URL + jolokiaApiPath
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	r.Close = true
	if strings.HasPrefix(opts.Jolokia.URL, "https://") {
		clientTLSCert, err := tls.LoadX509KeyPair(opts.Jolokia.Cert, opts.Jolokia.Key)
		if err != nil {
			return err
		}
		RootCAPool, err := x509.SystemCertPool()
		if err != nil {
			return err
		}
		if caCert, err := os.ReadFile(opts.Jolokia.CA); err != nil {
			return err
		} else if ok := RootCAPool.AppendCertsFromPEM(caCert); !ok {
			return err
		}
		tlsConfig := &tls.Config{
			RootCAs: RootCAPool,
			Certificates: []tls.Certificate{
				clientTLSCert,
			},
			InsecureSkipVerify: opts.Jolokia.InsecureSkipVerify,
		}
		tr := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		hc = &http.Client{Transport: tr}
	}
	if opts.Jolokia.BasicAuth {
		auth := opts.Jolokia.Username + ":" + opts.Jolokia.Password
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		r.Header.Set("Authorization", basicAuth)
	}

	resp, err := hc.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("get jolokia api failed, status code: %d", resp.StatusCode)
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	err = json.Unmarshal(buf.Bytes(), broker)
	if err != nil {
		return err
	}

	now := broker.Timestamp
	ch <- c.now.mustNewConstMetric(float64(now))
	return nil
}
