# ActiveMQ Exporter

The ActiveMQ exporter is a Prometheus exporter for ActiveMQ. It exposes metrics about the broker, queues and topics.

The exporter is written in Go and is based on the [Prometheus exporter template](https://github.com/prometheus/client_golang/blob/master/examples/random/main.go).

## Building

### Via Git

```
git clone https://github.com/atompi/activemq_exporter.git
cd activemq_exporter
go build -ldflags '-s -w'
```

### Via Go get

```
go get github.com/atompi/activemq_exporter
```

### Via Docker

```
git clone https://github.com/atompi/activemq_exporter.git
cd activemq_exporter
docker build -t atompi/activemq-exporter .
```

## Usage

The `activemq_exporter` listens on HTTP port 9168 by default. See the `--help` output for more options.

### Via Systemd

[systemd README](./examples/systemd/README.md)

### Via Docker Compose

```
docker compose up -d
```
