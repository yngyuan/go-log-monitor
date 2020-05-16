# go-log-monitor
a log monitor tool written in go, connected to influxDB and grafana.

### project outline
**nginx log data --> log process(realtime reading, parsing, and wrting) -> save to influxDB --> show in grafana**

![pipeline](https://github.com/yngyuan/go-log-monitor/blob/master/grafana.png?raw=true)

### influxDB
InfluxDB is an open-source time series database (TSDB) developed by InfluxData. It is written in Go and optimized for fast, high-availability storage and retrieval of time series data in fields such as operations monitoring, application metrics, Internet of Things sensor data, and real-time analytics. It also has support for processing data from Graphite.

https://github.com/influxdata/influxdb

### grafana
The tool for beautiful monitoring and metric analytics & dashboards for Graphite, InfluxDB & Prometheus & More 

https://github.com/grafana/grafana
