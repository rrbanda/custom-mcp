# Process Automation MCP Server

An MCP (Model Context Protocol) server that exposes backend APIs as tools for AI assistants. Built with [GenMCP](https://github.com/rrbanda/gen-mcp) - zero code required.

## What This Server Does

This MCP server provides tools for **process automation** - any workflow that involves calling backend APIs. Current tools include:

- **Certificate Management**: Create certificate orders, check status
- **(Extensible)**: Add tools for any API-driven process (approvals, provisioning, service requests, etc.)

The server is designed to be easily extended with new tools without writing code - just add YAML configurations.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Complete Local Walkthrough](#complete-local-walkthrough)
3. [Quick Start](#quick-start)
4. [How It Works](#how-it-works)
5. [Deploy to OpenShift](#deploy-to-openshift)
6. [Testing Your Server](#testing-your-server)
7. [Customization](#customization)
8. [Troubleshooting](#troubleshooting)
9. [Reference](#reference)

---

## Prerequisites

| Requirement | Purpose | Installation |
|-------------|---------|--------------|
| **Go 1.21+** | Build GenMCP CLI | `brew install go` (macOS) |
| **GenMCP CLI** | Build and run MCP servers | See below |
| **podman or docker** | Build container images | `brew install podman` |
| **oc CLI** | Deploy to OpenShift (optional) | [Download](https://mirror.openshift.com/pub/openshift-v4/clients/ocp/) |

### Install GenMCP CLI

> **Important**: This project requires GenMCP with schemaVersion 0.2.0 support.  
> The released version (v0.1.x) does NOT work. You must build from source.

```bash
# Clone GenMCP repository (use this specific repo)
git clone https://github.com/rrbanda/gen-mcp.git
cd gen-mcp

# Build the CLI
make build-cli

# Move to a directory in your PATH (replaces any old version)
sudo mv genmcp /usr/local/bin/

# Verify installation - should show "development" version
genmcp version
# Expected: genmcp version development@<hash>
```

**Verify version works with 0.2.0 schema:**
```bash
genmcp version
# If you see "v0.1.x", you have the old version - rebuild from source!
```

---

## Complete Local Walkthrough

This section walks you through **every step from scratch** to get the MCP server running locally.

### Step 1: Install Prerequisites

```bash
# macOS
brew install go
brew install podman  # or docker

# Verify installations
go version    # Should show 1.21+
podman --version
```

### Step 2: Install GenMCP CLI

> **Important**: Must build from source to get schemaVersion 0.2.0 support.

```bash
# Clone GenMCP (this specific repo has 0.2.0 support)
git clone https://github.com/rrbanda/gen-mcp.git
cd gen-mcp
make build-cli

# Install to PATH (replaces any old version)
sudo mv genmcp /usr/local/bin/

# Verify - should show "development" version, NOT "v0.1.x"
genmcp version
```

### Step 3: Clone This Repository

```bash
git clone https://github.com/rrbanda/custom-mcp.git
cd custom-mcp/cmp
```

### Step 4: Review the Configuration Files

| File | Purpose | Key Contents |
|------|---------|--------------|
| `mcpfile.yaml` | Defines the MCP tools | Tool name, input schema, backend URL pattern |
| `mcpserver.yaml` | Runtime configuration | Port (8080), transport protocol, logging |

**No code changes needed** - these files define everything.

### Step 5: Set Your Backend URL

```bash
# Set your CMP backend URL (REQUIRED)
export CMP_BACKEND_URL="https://your-cmp-api.company.com"

# Verify it's set
echo $CMP_BACKEND_URL
```

### Step 6: Start the MCP Server

```bash
# Start the server
genmcp run
```

**Expected output:**
```
{"level":"info","msg":"Loaded 1 tools from .../mcpfile.yaml"}
{"level":"info","msg":"Starting MCP server","port":8080}
{"level":"info","msg":"Starting HTTP server"}
```

Keep this terminal running. Open a **new terminal** for testing.

### Step 7: Test the Server

In a new terminal:

```bash
# List available tools
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/list","id":1}'
```

**Expected:** Response containing `"name":"create_certificate_order"`

### Step 8: Stop the Server

Press `Ctrl+C` in the terminal running the server.

### Step 9: (Optional) Build Container Image

```bash
# Build the container image
genmcp build \
  --tag quay.io/YOUR-ORG/cmp-mcp:v1.0.0 \
  -f mcpfile.yaml \
  --platform linux/amd64

# Login to your registry
podman login quay.io

# Push the image
podman push quay.io/YOUR-ORG/cmp-mcp:v1.0.0
```

---

### Summary: Complete Local Flow

```
┌─────────────────────────────────────────────────────────────────┐
│ Step 1: Install Go + Podman                                     │
│    ↓                                                            │
│ Step 2: Install GenMCP CLI                                      │
│    ↓                                                            │
│ Step 3: Clone this repo, cd custom-mcp/cmp                      │
│    ↓                                                            │
│ Step 4: Review mcpfile.yaml and mcpserver.yaml                 │
│    ↓                                                            │
│ Step 5: export CMP_BACKEND_URL="https://..."                   │
│    ↓                                                            │
│ Step 6: genmcp run                                              │
│    ↓                                                            │
│ Step 7: curl -X POST http://localhost:8080/mcp ...             │
│    ↓                                                            │
│ Step 8: Ctrl+C to stop                                          │
│    ↓                                                            │
│ Step 9: (Optional) genmcp build + podman push                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Quick Start

> **Already have GenMCP CLI installed?** Use this quick start.  
> **New to this?** Use the [Complete Local Walkthrough](#complete-local-walkthrough) above.

```bash
cd cmp
export CMP_BACKEND_URL="https://your-cmp-api.company.com"
genmcp run
```

---

## How It Works

### What is GenMCP?

GenMCP generates MCP servers from YAML config files. **You write NO code.**

You provide:
- `mcpfile.yaml` - Defines tools (what the server can do)
- `mcpserver.yaml` - Defines runtime settings (port, logging)

GenMCP provides the server binary that reads these configs.

### Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         MCP Client                              │
│              (Claude, Cursor, MCP Inspector, etc.)              │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ MCP Protocol (JSON-RPC over HTTP)
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      MCP Server Container                       │
│                                                                 │
│   mcpfile.yaml → Defines: create_certificate_order tool        │
│                                                                 │
│   When tool is called:                                          │
│   1. Receives parameters (username, token, questionnaire, etc.) │
│   2. Adds auth headers (Username, Token, SM_USER_DEV)           │
│   3. Forwards to: ${CMP_BACKEND_URL}/api/orderextws/...         │
│   4. Returns response to client                                 │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### File Structure

```
cmp/
├── mcpfile.yaml              # Tool definitions
├── mcpserver.yaml            # Runtime config (port 8080, /mcp path)
├── README.md                 # This file
└── openshift/
    └── config/
        ├── configmap.yaml    # Tool definitions for k8s
        ├── deployment.yaml   # Pod specification
        ├── service.yaml      # Internal service
        └── route.yaml        # External HTTPS route
```

---

## Deploy to OpenShift

This section walks through deploying the MCP server to OpenShift/Kubernetes.

> **No Dockerfile needed!** The `genmcp build` command creates a container image directly from your YAML config files. It embeds `mcpfile.yaml` and `mcpserver.yaml` into the image along with the GenMCP server binary.

### Step 1: Build container image

```bash
cd cmp

# Build the container image (genmcp generates everything internally)
genmcp build \
  --tag quay.io/YOUR-ORG/cmp-mcp:v1.0.0 \
  -f mcpfile.yaml \
  --platform linux/amd64
```

**What `genmcp build` does:**
- Creates a minimal container with the GenMCP server binary
- Embeds your `mcpfile.yaml` and `mcpserver.yaml`
- No Dockerfile required - it's all handled for you

### Step 2: Push image to registry

```bash
# Login to your container registry
podman login quay.io

# Push the image
podman push quay.io/YOUR-ORG/cmp-mcp:v1.0.0
```

> **Alternative**: Use the pre-built image `quay.io/rbrhssa/mcp-techx:v1.0.1` to skip Steps 1-2.

### Step 3: Update deployment configuration

Edit `openshift/config/deployment.yaml` and update these TWO values:

```yaml
# Line 29 - Your container image
image: quay.io/YOUR-ORG/cmp-mcp:v1.0.0

# Line 43 - Your CMP backend URL
- name: CMP_BACKEND_URL
  value: "https://your-cmp-api.company.com"
```

### Step 4: Login to OpenShift

```bash
oc login --token=<your-token> --server=<your-cluster-api>
```

### Step 5: Create project and deploy

```bash
oc new-project cmp-mcp

# Deploy all resources (order matters!)
oc apply -f openshift/config/configmap.yaml
oc apply -f openshift/config/deployment.yaml
oc apply -f openshift/config/service.yaml
oc apply -f openshift/config/route.yaml
```

### Step 6: Verify deployment

```bash
oc get pods
oc get route cmp-mcp-server -o jsonpath='{.spec.host}'
```

### Step 7: Test deployed server

```bash
ROUTE=$(oc get route cmp-mcp-server -o jsonpath='{.spec.host}')

curl -X POST "https://${ROUTE}/mcp" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/list","id":1}'
```

---

## Testing Your Server

### Test with curl

```bash
# List available tools
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/list","id":1}'

# Call a tool
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc":"2.0",
    "method":"tools/call",
    "id":2,
    "params":{
      "name":"create_certificate_order",
      "arguments":{
        "username":"testuser",
        "token":"test-token",
        "smUserDev":"0000000003",
        "orderByGeid":"1011035821",
        "orderForGeid":"1011035821",
        "productxd":"1234_1234_GLOBAL",
        "productName":"Test Product",
        "questionnaire":[{"dataKey":"123","key":"Environment","isHidden":"false","value":["DEV"]}],
        "locationIdOrderFor":"37750",
        "geoIdOrderFor":"USA"
      }
    }
  }'
```

### Test with MCP Inspector

```bash
npx @modelcontextprotocol/inspector
```

Connect to: `http://localhost:8080/mcp`

---

## Customization

### Understanding What's Configurable

#### 1. Deployment-time Configuration (set in deployment.yaml)

| Value | Where | Purpose |
|-------|-------|---------|
| `image` | deployment.yaml line 29 | Your container image |
| `CMP_BACKEND_URL` | deployment.yaml line 43 | Base URL for your CMP API |

#### 2. Runtime Parameters (passed when tool is called)

These are **NOT hardcoded** - they're passed by the AI/user each time:

| Parameter | Purpose | Passed As |
|-----------|---------|-----------|
| `username` | CMP username | Tool input → `Username` header |
| `token` | Auth token | Tool input → `Token` header |
| `smUserDev` | SM_USER_DEV ID | Tool input → `SM_USER_DEV` header |

### Adding New Tools

Edit `mcpfile.yaml` and add to the `tools` array. See the existing `create_certificate_order` tool as a template.

After editing, update the ConfigMap on OpenShift:

```bash
oc delete configmap cmp-mcp-config
oc apply -f openshift/config/configmap.yaml
oc rollout restart deployment/cmp-mcp-server
```

---

## Troubleshooting

### "invalid mcp file version, expected 0.1.0"

**Cause**: You have an old version of genmcp (v0.1.x) that doesn't support schemaVersion 0.2.0.

**Fix**: Rebuild genmcp from source:
```bash
cd gen-mcp  # wherever you cloned it
git pull
make build-cli
sudo mv genmcp /usr/local/bin/
genmcp version  # Should show "development@..." not "v0.1.x"
```

### Pod keeps restarting

**Fix**: Ensure probes use TCP socket (MCP only accepts POST):
```yaml
livenessProbe:
  tcpSocket:
    port: 8080
```

### "dial tcp: lookup ... i/o timeout"

**Fix**: Backend URL not reachable from cluster. Check `CMP_BACKEND_URL`.

### Tool not appearing

**Fix**: Update ConfigMap and restart:
```bash
oc delete configmap cmp-mcp-config
oc apply -f openshift/config/configmap.yaml
oc rollout restart deployment/cmp-mcp-server
```

### Check logs

```bash
# OpenShift
oc logs deployment/cmp-mcp-server

# Local (logs in terminal running genmcp)
```

---

## Reference

### Environment Variables

| Variable | Description |
|----------|-------------|
| `CMP_BACKEND_URL` | Base URL for CMP API (required) |

### Ports

| Environment | Port | URL |
|-------------|------|-----|
| Local | 8080 | `http://localhost:8080/mcp` |
| OpenShift | 443 | `https://<route-host>/mcp` |

### Pre-built Image

```
quay.io/rbrhssa/mcp-techx:v1.0.1
```

---

## License

Apache 2.0 - See [LICENSE](../LICENSE) for details.
