# Systemd Unit

The unit files (`*.service` and `*.socket`) in this directory are to be put into `/etc/systemd/system`.

It needs a user named `activemq_exporter`, whose shell should be `/sbin/nologin` and should not have any special privileges.

It needs a sysconfig file in `/etc/sysconfig/activemq_exporter`.

It needs a config file in `/etc/activemq_exporter/activemq_exporter.yaml`.

A sample config file can be found in `../activemq_exporter.yaml.example`.

A sample file can be found in `sysconfig.activemq_exporter`.
