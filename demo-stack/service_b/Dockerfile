FROM golang:1.24.4-alpine3.22 AS builder
WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build cd service_b && go mod tidy 
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux go build -o consumer service_b

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/consumer .
CMD ["./consumer"]