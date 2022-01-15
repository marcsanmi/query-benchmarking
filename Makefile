.PHONY: check run-promscale run-timescaledb stop-promscale stop-timescaledb test help

check: ## Run checks
	bash scripts/lint.sh
	bash scripts/tidy.sh

run-promscale: ## Run Promscale container
	@export $$(cat .env | grep -v "#" | xargs) && \
	bash scripts/promscale/run.sh

run-timescaledb: ## Install and run TimescaleDB with Promscale extension
	@export $$(cat .env | grep -v "#" | xargs) && \
	bash scripts/timescaledb/run.sh

stop-promscale: ## Stop Promscale container
	@export $$(cat .env | grep -v "#" | xargs) && \
	bash scripts/promscale/stop.sh

stop-timescaledb: ## Stop timescaledb container
	@export $$(cat .env | grep -v "#" | xargs) && \
	bash scripts/timescaledb/stop.sh

test: ## Run tests
	bash scripts/test.sh

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
