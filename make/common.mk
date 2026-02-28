# =============================================================================
# MCP Toolkit - Common Make Configuration
# =============================================================================
# Shared variables and helpers used across all makefiles.

SHELL := /bin/bash

# =============================================================================
# Project Configuration
# =============================================================================
BINARY := mcp-toolkit
GOMOD := ./...
GOMOD_COVERAGE := ./internal/...

# =============================================================================
# Version Information
# =============================================================================
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags with version info
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)

# =============================================================================
# Platform Configuration
# =============================================================================
HOST_OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
HOST_ARCH := $(shell uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/')
TARGET_OS ?= $(HOST_OS)
TARGET_ARCH ?= $(HOST_ARCH)

PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# =============================================================================
# Directory Structure
# =============================================================================
BUILD_DIR := bin
REPORTS_DIR := reports
TMP_DIR := .tmp

# =============================================================================
# Coverage Configuration
# =============================================================================
COVERAGE_THRESHOLD := 90
COVERAGE_FILE := $(REPORTS_DIR)/coverage.out
COVERAGE_HTML := $(REPORTS_DIR)/coverage.html

# =============================================================================
# Micromamba Environment (auto-managed)
# =============================================================================
MAMBA_ENV := mcp-toolkit
MAMBA_SPEC := environment.yml
MAMBA_RUN := micromamba run -n $(MAMBA_ENV)
ENV_MARKER := $(TMP_DIR)/.env-marker

# Check if environment exists
define mamba_env_exists
	micromamba env list 2>/dev/null | grep -q "$(MAMBA_ENV)"
endef

# =============================================================================
# Container Runtime Detection (Docker/Podman/Rancher)
# =============================================================================
CONTAINER_CMD ?= $(shell command -v docker >/dev/null 2>&1 && echo docker || \
	(command -v podman >/dev/null 2>&1 && echo podman || \
	(command -v nerdctl >/dev/null 2>&1 && echo nerdctl || echo docker)))

# =============================================================================
# Auto Environment Setup (transparent to user)
# =============================================================================
# This target ensures the environment is created/updated as needed.
# - Creates env if it doesn't exist
# - Updates env if environment.yml changed
# - Does nothing if env is up-to-date

.PHONY: ensure-env
ensure-env:
	@mkdir -p $(TMP_DIR)
	@if ! $(call mamba_env_exists); then \
		echo "🔧 Creating micromamba environment '$(MAMBA_ENV)'..."; \
		micromamba create -y -f $(MAMBA_SPEC); \
		touch $(ENV_MARKER); \
	elif [ ! -f $(ENV_MARKER) ] || [ $(MAMBA_SPEC) -nt $(ENV_MARKER) ]; then \
		echo "🔄 Updating micromamba environment '$(MAMBA_ENV)'..."; \
		micromamba env update -n $(MAMBA_ENV) -f $(MAMBA_SPEC); \
		touch $(ENV_MARKER); \
	fi

# =============================================================================
# Manual Environment Targets
# =============================================================================

.PHONY: setup-env
setup-env: ## Create/update micromamba environment
	@mkdir -p $(TMP_DIR)
	@if ! $(call mamba_env_exists); then \
		echo "🔧 Creating micromamba environment '$(MAMBA_ENV)'..."; \
		micromamba create -y -f $(MAMBA_SPEC); \
	else \
		echo "🔄 Updating micromamba environment '$(MAMBA_ENV)'..."; \
		micromamba env update -n $(MAMBA_ENV) -f $(MAMBA_SPEC); \
	fi
	@touch $(ENV_MARKER)
	@echo "✓ Environment '$(MAMBA_ENV)' ready"

.PHONY: clean-env
clean-env: ## Remove micromamba environment
	micromamba env remove -n $(MAMBA_ENV) -y 2>/dev/null || true
	rm -f $(ENV_MARKER)
	@echo "✓ Environment removed"
