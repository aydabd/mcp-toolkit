# =============================================================================
# MCP Toolkit - Build Targets
# =============================================================================

.PHONY: build build-all install clean

build: ensure-env ## Build binary for current OS/arch
	@mkdir -p $(BUILD_DIR)
	$(MAMBA_RUN) go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY) ./cmd/mcp-toolkit
	@echo "✓ Built $(BUILD_DIR)/$(BINARY)"

build-all: ensure-env ## Build binaries for all platforms
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; arch=$${platform#*/}; \
		ext=""; [ "$$os" = "windows" ] && ext=".exe"; \
		echo "Building $$os/$$arch..."; \
		$(MAMBA_RUN) env GOOS=$$os GOARCH=$$arch go build -ldflags "$(LDFLAGS)" \
			-o $(BUILD_DIR)/$(BINARY)-$$os-$$arch$$ext ./cmd/mcp-toolkit; \
	done
	@echo "✓ Built all platforms"

install: ensure-env build ## Install to GOPATH/bin
	$(MAMBA_RUN) go install -ldflags "$(LDFLAGS)" ./cmd/mcp-toolkit
	@echo "✓ Installed to GOPATH/bin"

clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR) $(REPORTS_DIR) coverage.out coverage.html
	@echo "✓ Cleaned"
