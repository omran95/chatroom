# Chat-APP
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/omran95/chat-app?label=Version&sort=semver)

Real-time chat app in a highly scalable architecture. 


### System architecture

<img width="1028" alt="image" src="https://raw.githubusercontent.com/omran95/chat-app/main/architecture.png">


### Features
- Real-time chatting using websockets.
- Services **are stateless** and can be horizontally scaled. (Redundancy).
  - `room`: creates rooms and handles messages.
  - `subscriber`: maintains Kafka subscriber topics for each room in a Redis cluster.

- gRPC for low-latency and high-throughput inter-service communication.
  - with retry (Exponential backoff with jitter), timeout, and circuit breaker.
- Graceful shutdown.
- Observability using Prometheus for service monitoring and OpenTelemetry + Jaeger for distributed tracing.
- Pub/Sub using Kafka.
- Persist messages and rooms in Cassandra, A highly available and scalable NoSQL Database with tunable consistency.
- Protect create room API with distributed rate limiter using Token-Bucket Algorithm with Redis.
- Broadcasting seen, typing, joining, and leaving events to all room members.
