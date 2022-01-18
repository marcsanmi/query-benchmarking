#!/usr/bin/env bash
set -eufo pipefail

# Check required commands are in place
command -v docker > /dev/null 2>&1 || { echo 'please install docker or use image that has it'; exit 1; }

# Remove any existing network
docker network rm "${DOCKER_NETWORK}" || true

# Remove any existing container
docker rm -f "${DOCKER_TIMESCALEDB_CONTAINER_NAME}" || true

# Create a specific network
docker network create --driver bridge "${DOCKER_NETWORK}" || true

# Run TimescaleDB with Promscale extension
docker run --name "${DOCKER_TIMESCALEDB_CONTAINER_NAME}" -d \
    -e POSTGRES_PASSWORD="${POSTGRES_PASSWORD}" \
    -e POSTGRES_USER="${POSTGRES_USER}" \
    -p "${TIMESCALEDB_HOST_PORT}:${TIMESCALEDB_CONTAINER_PORT}" \
    --network "${DOCKER_NETWORK}" \
    "${DOCKER_TIMESCALEDB_IMAGE}" "${POSTGRES_DB_NAME}" "${POSTGRES_ARGS}"
