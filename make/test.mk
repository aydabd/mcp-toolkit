# =============================================================================
# MCP Toolkit - Test Targets
# =============================================================================

.PHONY: test test-cover test-verbose coverage-check

test: ensure-env ## Run tests
	$(MAMBA_RUN) go test -v $(GOMOD)
	@echo "✓ Tests passed"

test-cover: ensure-env ## Run tests with coverage report
	@mkdir -p $(REPORTS_DIR)
	$(MAMBA_RUN) go test -coverprofile=$(COVERAGE_FILE) -covermode=atomic $(GOMOD_COVERAGE)
	$(MAMBA_RUN) go tool cover -func=$(COVERAGE_FILE)
	$(MAMBA_RUN) go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "✓ Coverage report: $(COVERAGE_HTML)"

test-verbose: ensure-env ## Run tests with verbose output
	$(MAMBA_RUN) go test -v -count=1 $(GOMOD)

coverage-check: test-cover ## Check coverage meets threshold
	@coverage=$$($(MAMBA_RUN) go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print $$3}' | tr -d '%'); \
	echo "Coverage: $${coverage}% (threshold: $(COVERAGE_THRESHOLD)%)"; \
	if [ $$(echo "$${coverage} < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage below threshold"; exit 1; \
	else \
		echo "✓ Coverage meets threshold"; \
	fi

# =============================================================================
# E2E Tests (Isolated)
# =============================================================================

.PHONY: e2e e2e-docker e2e-local

e2e: e2e-local ## Run e2e tests (isolated temp HOME)

e2e-local: ensure-env ## Run e2e tests with isolated temp HOME
	$(MAMBA_RUN) go test -tags=integration -v ./test/e2e/...
	@echo "✓ E2E tests passed (isolated)"

e2e-docker: ## Run e2e tests in Docker container
	@chmod +x test/e2e/run.sh
	@test/e2e/run.sh
