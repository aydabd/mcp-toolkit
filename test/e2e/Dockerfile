# E2E Test Container for mcp-toolkit
# Simulates editor environments to test mcp.json generation
FROM golang:1.25-alpine

RUN apk add --no-cache git bash jq

WORKDIR /app

# Create isolated home directory with simulated editor configs
ENV HOME=/home/testuser

# Simulate installed editors by creating their config directories
RUN mkdir -p \
    $HOME/.config/Code/User \
    $HOME/.config/Cursor/User \
    $HOME/.config/Windsurf/User \
    $HOME/.config/zed \
    $HOME/.mcp-server-envs \
    $HOME/.kube

# Copy source and build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /usr/local/bin/mcp-toolkit ./cmd/mcp-toolkit

# Create test user
RUN adduser -D -h $HOME testuser && chown -R testuser:testuser $HOME /app
USER testuser

# Default entrypoint for running tests
ENTRYPOINT ["mcp-toolkit"]
