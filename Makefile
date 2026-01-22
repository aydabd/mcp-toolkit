# =============================================================================
# MCP Toolkit - Main Makefile
# =============================================================================
# Entry point that includes all modular makefiles.
#
# Structure:
#   make/common.mk  - Shared variables and environment setup
#   make/build.mk   - Build and install targets
#   make/test.mk    - Unit tests and coverage
#   make/lint.mk    - Code formatting and linting
#   make/dev.mk     - Development workflow targets
#
# Usage:
#   make help       - Show all available commands
#   make ci         - Run complete CI pipeline
#   make quickstart - Run interactive MCP setup
# =============================================================================

include make/common.mk
include make/build.mk
include make/test.mk
include make/lint.mk
include make/dev.mk

# =============================================================================
# Help Target
# =============================================================================

.PHONY: help
help: ## Show available commands
	@echo "MCP Toolkit - Available Commands"
	@echo ""
	@echo "Environment:"
	@grep -hE '^[a-zA-Z0-9_-]+:.*?## .*$$' make/common.mk | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "Build:"
	@grep -hE '^[a-zA-Z0-9_-]+:.*?## .*$$' make/build.mk | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "Test:"
	@grep -hE '^[a-zA-Z0-9_-]+:.*?## .*$$' make/test.mk | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "Lint:"
	@grep -hE '^[a-zA-Z0-9_-]+:.*?## .*$$' make/lint.mk | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "Dev:"
	@grep -hE '^[a-zA-Z0-9_-]+:.*?## .*$$' make/dev.mk | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
