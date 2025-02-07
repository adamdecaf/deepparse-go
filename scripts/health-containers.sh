#!/bin/bash

timeout=300  # 5 minutes
start_time=$(date +%s)

until docker ps -q | xargs -n1 docker inspect --format '{{.State.Health.Status}}' | grep -q healthy; do
  current_time=$(date +%s)
  elapsed=$((current_time - start_time))

  if [ $elapsed -gt $timeout ]; then
    echo "Timeout waiting for container health check"
    exit 1
  fi

  echo "Waiting for container to be healthy... ($elapsed seconds elapsed)"
  sleep 5
done

echo "Container is healthy!"
