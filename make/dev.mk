# =============================================================================
# MCP Toolkit - Development Targets
# =============================================================================

.PHONY: run quickstart deps pre-commit ci update

run: ensure-env ## Run CLI without building
	$(MAMBA_RUN) go run ./cmd/mcp-toolkit $(ARGS)

quickstart: ensure-env ## Run quickstart command
	$(MAMBA_RUN) go run ./cmd/mcp-toolkit quickstart

deps: ensure-env ## Download and verify dependencies
	$(MAMBA_RUN) go mod download
	$(MAMBA_RUN) go mod verify
	@echo "✓ Dependencies ready"

pre-commit: ensure-env ## Install pre-commit hooks
	$(MAMBA_RUN) pre-commit install
	@echo "✓ Pre-commit hooks installed"

update: ## Update all dependencies (micromamba + Go)
	@echo "📦 Updating micromamba environment..."
	micromamba update -y -n $(MAMBA_ENV) --all
	@touch $(ENV_MARKER)
	@echo "📦 Updating Go dependencies..."
	$(MAMBA_RUN) go get -u ./...
	$(MAMBA_RUN) go mod tidy
	@echo "✓ All dependencies updated"

ci: ensure-env lint test-cover ## Run complete CI pipeline
	@echo "✓ CI passed"
