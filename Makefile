# Engineering Foundation — Sprint 0
# Stack-agnostic. Stack-specific targets slot in after a stack is selected
# post-C4-spike. See CONTRIBUTING.md.

.DEFAULT_GOAL := help

PRECOMMIT := $(shell command -v pre-commit 2>/dev/null)
GITLEAKS  := $(shell command -v gitleaks 2>/dev/null)
DOCKER    := $(shell command -v docker 2>/dev/null)

.PHONY: help init lint format docs-lint secrets check-frozen frozen-baseline devshell test ci upgrade

help: ## list available targets
	@awk 'BEGIN {FS = ":.*##"; printf "Engineering Foundation targets:\n\n"} /^[a-zA-Z][a-zA-Z0-9_-]*:.*##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

init: ## first-time setup — install git hooks (pre-commit + commit-msg)
	@if [ -z "$(PRECOMMIT)" ]; then echo "pre-commit not found; install: pip install pre-commit (or brew/conda)"; exit 1; fi
	pre-commit install --install-hooks
	@echo "hooks installed (pre-commit + commit-msg)."

lint: ## run pre-commit on all files
	@if [ -z "$(PRECOMMIT)" ]; then echo "lint: pre-commit not found"; exit 1; fi
	pre-commit run --all-files

format: ## run auto-fixing hooks
	@if [ -z "$(PRECOMMIT)" ]; then echo "format: pre-commit not found"; exit 1; fi
	pre-commit run --all-files trailing-whitespace end-of-file-fixer mixed-line-ending || true

docs-lint: ## markdown lint (markdownlint-cli2 via npx)
	@if ! command -v npx >/dev/null 2>&1; then echo "docs-lint: npx not installed — SKIPPED (run in CI)"; else npx --yes markdownlint-cli2 '**/*.md' '!claude-skills/**' '!lab/evidence/**'; fi

secrets: ## scan for secrets (gitleaks)
	@if [ -z "$(GITLEAKS)" ]; then echo "secrets: gitleaks not installed — SKIPPED (run in CI)"; else gitleaks detect --redact --config .gitleaks.toml --source .; fi

check-frozen: ## verify frozen planning docs unchanged (SHA-256 vs FROZEN.sha256)
	@bash scripts/check-frozen-docs.sh

frozen-baseline: ## (re)capture frozen-doc hashes — ONLY after a documented amendment
	@sha256sum $$(grep -vE '^[[:space:]]*(#|$$)' scripts/frozen-docs.list) > FROZEN.sha256
	@echo "FROZEN.sha256 refreshed — commit it alongside the amendment + journal entry."

devshell: ## run the stack-agnostic dev-toolbox container (builds if needed)
	@if [ -z "$(DOCKER)" ]; then echo "devshell: docker not found"; exit 1; fi
	docker build -t eng-toolbox .
	docker run --rm -it -v "$$PWD":/ws -w /ws eng-toolbox

test: ## run tests (DEFERRED — no stack selected post-C4-spike; see CONTRIBUTING.md)
	@echo "test: no language stack selected yet (Sprint 0 is stack-agnostic)."
	@echo "      Stack-specific tests slot in once a stack is chosen after the C4 spike."

ci: ## run the local equivalent of CI (lint + docs-lint + check-frozen + secrets + test)
ci: lint docs-lint check-frozen secrets test
	@echo "ci: all gates passed."

upgrade: ## bump pre-commit hooks to latest (pre-commit autoupdate)
	@if [ -z "$(PRECOMMIT)" ]; then echo "pre-commit not found"; exit 1; fi
	pre-commit autoupdate
