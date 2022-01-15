#!/usr/bin/env bash
set -eufo pipefail

# Check required commands are in place
command -v docker > /dev/null 2>&1 || { echo 'please install docker or use image that has it'; exit 1; }

# Remove any existing container
docker rm -f "${DOCKER_PROMSCALE_CONTAINER_NAME}" || true

# Run Promscale
docker run --name "${DOCKER_PROMSCALE_CONTAINER_NAME}" \
    -p "${PROMSCALE_HOST_PORT}:${PROMSCALE_CONTAINER_PORT}" \
    --network "${DOCKER_NETWORK}" \
    "${DOCKER_PROMSCALE_IMAGE}" \
    -db-password="${POSTGRES_PASSWORD}" \
    -db-port="${TIMESCALEDB_CONTAINER_PORT}" \
    -db-name="${POSTGRES_DB_NAME}" \
    -db-host="${POSTGRES_DB_HOST}" \
    -db-ssl-mode="allow"
