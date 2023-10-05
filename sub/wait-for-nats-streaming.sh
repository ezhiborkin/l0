#!/bin/sh
# wait-for-nats-streaming.sh
# Wait for NATS Streaming server to be ready

until nc -z -v -w30 nats-streaming 4222; do
  echo "Waiting for NATS Streaming server to start..."
  sleep 5
done
echo "NATS Streaming server is up - starting clients."

