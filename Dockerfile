FROM golang:1.22.5 AS builder
ENV GOPROXY='https://proxy.golang.com.cn,direct'
WORKDIR /app
COPY . .
RUN go build -ldflags '-s -w'

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/activemq_exporter /app/activemq_exporter
COPY --from=builder /app/examples/activemq_exporter.yaml.example /app/activemq_exporter.yaml
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

ENTRYPOINT ["/app/activemq_exporter", "--collector.config.file=/app/activemq_exporter.yaml"]
