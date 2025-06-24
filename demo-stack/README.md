# Workshop observability
* Service to service
  * HTTP
  * RabbitMQ
  * MongoDB
* Golang
  * Go workspace
  * Go modules

## Start RabbitMQ
```
$docker compose up -d rabbitmq
$docker compose ps
```

Access to managment UI
* http://localhost:15672/
  * username=guest
  * password=guest

## Build and Run service-a
```
$docker compose build service_a
$docker compose up -d service_a
$docker compose ps
```

List of APIs
* POST /user
```
{
    "name": "somkiat",
    "email": "somkiat@xxx.com"
}
```

### Error message of Service-a
```
traces export: Post "http://otel-collector:4318/v1/traces": dial tcp: lookup otel-collector on 127.0.0.11:53: no such host
```

## Start Otel Collector server
```
$docker compose up -d jaeger
$docker compose up -d otel-collector
$docker compose ps
```

Go to Jaeger Server
* http://localhost:16686/search


## Build and Run service-b :: Consumer
* Receive data from RabitMQ server

```
$docker compose build service_b
$docker compose up -d service_b
$docker compose ps
```

## Build and Run service-c :: Product service
* Connect to MongoDB

### Start MongoDB
```
$docker compose up -d mongo
$docker compose ps
```

## Start service-c
```
$docker compose build service_c
$docker compose up -d service_c
$docker compose ps
```

List of URLs
* http://localhost:8083/products

## Alert rules
* [Alert manager](https://prometheus.io/docs/alerting/latest/alertmanager/)
* Prometheus
* Grafana

### Start Alert manager server
```
$docker compose up -d alertmanager
$docker compose ps
```

Alert manager dashboard
* http://localhost:9093/

### Start Prometheus server
```
$docker compose up -d prometheus
$docker compose ps
```

Prometheus dashboard
* http://localhost:9090/
