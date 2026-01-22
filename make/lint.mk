# =============================================================================
# MCP Toolkit - Lint Targets
# =============================================================================

.PHONY: lint fmt vet tidy

lint: ensure-env fmt vet tidy ## Run all linters
	@echo "✓ Lint complete"

fmt: ensure-env ## Format Go code
	$(MAMBA_RUN) gofmt -w .
	@if command -v goimports >/dev/null 2>&1; then \
		$(MAMBA_RUN) go install golang.org/x/tools/cmd/goimports@latest 2>/dev/null || true; \
		$(MAMBA_RUN) goimports -w . 2>/dev/null || true; \
	fi
	@echo "✓ Formatted"

vet: ensure-env ## Run go vet
	$(MAMBA_RUN) go vet $(GOMOD)
	@echo "✓ Vet passed"

tidy: ensure-env ## Tidy go.mod and go.sum
	$(MAMBA_RUN) go mod tidy
	@echo "✓ Tidied"
