#!/bin/bash
set -euo pipefail

timeout=1800  # 30 minutes
start_time=$(date +%s)

healthy() {
  local ids status
  ids=$(docker compose ps -q 2>/dev/null || true)
  if [ -z "$ids" ]; then
    return 1
  fi
  status=$(echo "$ids" | xargs -n1 docker inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}' 2>/dev/null)
  [ "$status" = "healthy" ]
}

until healthy; do
  current_time=$(date +%s)
  elapsed=$((current_time - start_time))

  if [ "$elapsed" -gt "$timeout" ]; then
    echo "Timeout waiting for container health check"
    docker compose ps || true
    docker compose logs || true
    exit 1
  fi

  echo "Waiting for container to be healthy... ($elapsed seconds elapsed)"
  sleep 5
done

echo "Container is healthy!"
