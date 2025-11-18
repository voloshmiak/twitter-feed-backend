#!/bin/bash

set -e

BOOTSTRAP_SERVER="kafka:29092"

echo "Kafka is ready. Proceeding with topic creation..."

echo "Creating topic 'events-to-process'..."
kafka-topics \
  --create \
  --if-not-exists \
  --bootstrap-server $BOOTSTRAP_SERVER \
  --topic events-to-process \
  --partitions 3 \
  --replication-factor 1

echo "Creating topic 'events-processed'..."
kafka-topics \
  --create \
  --if-not-exists \
  --bootstrap-server $BOOTSTRAP_SERVER \
  --topic events-processed \
  --partitions 3 \
  --replication-factor 1

echo "All topics created successfully."