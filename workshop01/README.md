# Workshop :: LGTM stack with Go

- Observability with LGTM stack
  - Application metric
  - Distributed tracing
  - Centralized log
- Go and [Gin web framework](https://gin-gonic.com/)
- Docker compose

## Install Loki logging plugin in Docker

```
$docker plugin install grafana/loki-docker-driver:3.3.2-arm64 --alias loki --grant-all-permissions
$docker plugin ls
ID             NAME          DESCRIPTION           ENABLED
2f1dbe207d07   loki:latest   Loki Logging Driver   true
```

## Build and run LGTM stack

```
docker-compose up -d prometheus
docker-compose up -d loki
docker-compose up -d tempo
docker-compose up -d grafana
```

List of URLs

- Grafana :: http://localhost:3000
- Prometheus :: http://localhost:9090

## Build order service

- Go + Gin
  - Metric with [go-gin-prometheus](https://github.com/zsais/go-gin-prometheus)
  - Loggin with [sloggin](https://github.com/samber/slog-gin)
  - [OpenTelemetry with go](https://opentelemetry.io/docs/languages/go/)
    - https://github.com/open-telemetry/opentelemetry-go-contrib

```
docker-compose build order
docker-compose up -d order
```

List of URLs

- http://localhost:8080/ping
- http://localhost:8080/metrics
