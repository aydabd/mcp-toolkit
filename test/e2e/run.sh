#!/bin/bash
# E2E test runner - tests mcp-toolkit against simulated editor environments
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

IMAGE_NAME="mcp-toolkit-e2e"

echo "Building e2e test image..."
docker build -t "$IMAGE_NAME" -f "$SCRIPT_DIR/Dockerfile" "$PROJECT_ROOT"

echo ""
echo "Running e2e tests..."
echo ""

# Test 1: Version command
echo "Test 1: Version command"
docker run --rm "$IMAGE_NAME" version
echo "PASS: Version command"
echo ""

# Test 2: Help command
echo "Test 2: Help command"
docker run --rm "$IMAGE_NAME" --help | grep -q "MCP Toolkit"
echo "PASS: Help command"
echo ""

# Test 3: Env command shows editor support
echo "Test 3: Editor support"
docker run --rm "$IMAGE_NAME" env | grep -q "MCP_EDITOR"
echo "PASS: Editor support"
echo ""

# Test 4: Detect simulated installed editors
echo "Test 4: Installed editors detection"
docker run --rm --entrypoint /bin/bash "$IMAGE_NAME" -c 'ls -la $HOME/.config/ | grep -E "Code|Cursor|Windsurf|zed" && echo "All editor directories present"'
echo "PASS: Installed editors detection"
echo ""

# Test 5: VS Code config directory exists
echo "Test 5: VS Code config path"
docker run --rm --entrypoint /bin/bash "$IMAGE_NAME" -c 'test -d $HOME/.config/Code/User && echo "VS Code dir exists"'
echo "PASS: VS Code config path"
echo ""

# Test 6: Cursor config directory exists
echo "Test 6: Cursor config path"
docker run --rm --entrypoint /bin/bash "$IMAGE_NAME" -c 'test -d $HOME/.config/Cursor/User && echo "Cursor dir exists"'
echo "PASS: Cursor config path"
echo ""

# Test 7: Windsurf config directory exists
echo "Test 7: Windsurf config path"
docker run --rm --entrypoint /bin/bash "$IMAGE_NAME" -c 'test -d $HOME/.config/Windsurf/User && echo "Windsurf dir exists"'
echo "PASS: Windsurf config path"
echo ""

# Test 8: Zed config directory exists
echo "Test 8: Zed config path"
docker run --rm --entrypoint /bin/bash "$IMAGE_NAME" -c 'test -d $HOME/.config/zed && echo "Zed dir exists"'
echo "PASS: Zed config path"
echo ""

# Test 9: Environment variables override
echo "Test 9: Environment variable override"
docker run --rm -e MCP_ATLASSIAN_URL=https://mycompany.atlassian.net "$IMAGE_NAME" env | grep -q "mycompany.atlassian.net"
echo "PASS: Environment variable override"
echo ""

# Test 10: Env file isolation
echo "Test 10: Env files are isolated in container"
docker run --rm --entrypoint /bin/bash "$IMAGE_NAME" -c 'echo "TEST=value" > $HOME/.mcp-server-envs/test.env && cat $HOME/.mcp-server-envs/test.env'
echo "PASS: Env file isolation"
echo ""

# Test 11: Quickstart help
echo "Test 11: Quickstart command"
docker run --rm "$IMAGE_NAME" quickstart --help | grep -q "Interactive setup"
echo "PASS: Quickstart command"
echo ""

# Test 12: Setup command
echo "Test 12: Setup command"
docker run --rm "$IMAGE_NAME" setup --help | grep -q "atlassian"
echo "PASS: Setup command"
echo ""

echo "======================================="
echo "All e2e tests passed!"
echo ""
echo "Simulated editors tested:"
echo "  - VS Code   (~/.config/Code/User/)"
echo "  - Cursor    (~/.config/Cursor/User/)"
echo "  - Windsurf  (~/.config/Windsurf/User/)"
echo "  - Zed       (~/.config/zed/)"
echo "======================================="
