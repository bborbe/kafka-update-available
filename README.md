# Kafka updates availabe

Print available update for all applications seen.

## Run

```bash
go run main.go \
-kafka-brokers=kafka:9092 \
-kafka-schema-registry-url=http://localhost:8081 \
-kafka-latest-version-topic=application-version-latest \
-kafka-installed-version-topic=application-version-installed \
-datadir=/tmp \
-v=2
```
