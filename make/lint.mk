# =============================================================================
# MCP Toolkit - Lint Targets
# =============================================================================

.PHONY: lint lint-check fmt vet tidy

lint: ensure-env ## Run all linters (fix mode)
	$(MAMBA_RUN) bash -c '\
		export PATH="$$(go env GOPATH)/bin:$$PATH" && \
		pre-commit run --all-files \
	'
	@echo "✓ Lint complete"

lint-check: ensure-env ## Run linters in check mode (CI)
	$(MAMBA_RUN) bash -c '\
		export PATH="$$(go env GOPATH)/bin:$$PATH" && \
		export LINT_MODE=check && \
		pre-commit run --all-files \
	'
	@echo "✓ Lint check complete"

fmt: ensure-env ## Format Go code
	$(MAMBA_RUN) gofmt -w .
	@if command -v goimports >/dev/null 2>&1; then \
		$(MAMBA_RUN) goimports -w . 2>/dev/null || true; \
	fi
	@echo "✓ Formatted"

vet: ensure-env ## Run go vet
	$(MAMBA_RUN) go vet $(GOMOD)
	@echo "✓ Vet passed"

tidy: ensure-env ## Tidy go.mod and go.sum
	$(MAMBA_RUN) go mod tidy
	@echo "✓ Tidied"
