# Atlas Lab — a real distributed-systems laboratory for the delegation
# primitive. The substrate-independent targets (bench, report) run anywhere;
# the substrate targets (up, experiment, capture) need a real Docker daemon.
.DEFAULT_GOAL := help

COMPOSE := docker compose -f docker/docker-compose.yml

.PHONY: help bench report up down config experiment run-all

help: ## list targets
	@awk 'BEGIN{FS=":.*##"; printf "Atlas Lab targets:\n\n"} /^[a-zA-Z_-]+:.*##/{printf "  \033[36m%-12s\033[0m %s\n",$$1,$$2}' $(MAKEFILE_LIST)

bench: ## run TrustPerf benchmarks (substrate-independent; runs anywhere)
	cd .. && go test ./atlas-lab/bench -bench . -benchmem

report: ## (re)generate atlas-lab/bench/RESULTS.md
	cd .. && go test ./atlas-lab/bench -run TestTrustPerfReport -report
	@echo "wrote atlas-lab/bench/RESULTS.md"

config: ## validate the docker topology syntax (no daemon needed)
	$(COMPOSE) config >/dev/null && echo "docker-compose.yml: valid"

up: ## bring up the two-domain SPIRE substrate (needs a Docker daemon)
	$(COMPOSE) up -d --build

down: ## tear down the substrate
	$(COMPOSE) down -v

experiment: ## run the flagship revocation-under-partition experiment
	./experiments/revocation-under-partition/run.sh

run-all: ## bench + substrate + experiment (skips substrate if no daemon)
	./scripts/run-all.sh
