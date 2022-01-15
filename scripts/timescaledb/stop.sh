#!/usr/bin/env bash
set -eufo pipefail

# Check required commands are in place
command -v docker > /dev/null 2>&1 || { echo 'please install docker or use image that has it'; exit 1; }

docker stop "${DOCKER_TIMESCALEDB_CONTAINER_NAME}"
