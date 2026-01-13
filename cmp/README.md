# Process Automation MCP Server

An MCP (Model Context Protocol) server that exposes backend APIs as tools for AI assistants. Built with [GenMCP](https://github.com/rrbanda/gen-mcp) - zero code required.

## What This Server Does

This MCP server provides tools for **process automation** - any workflow that involves calling backend APIs. Current tools include:

- **Certificate Management**: Create certificate orders, check status
- **(Extensible)**: Add tools for any API-driven process (approvals, provisioning, service requests, etc.)

The server is designed to be easily extended with new tools without writing code - just add YAML configurations.

---

## What is MCP?

**MCP (Model Context Protocol)** is an open standard that allows AI assistants (like Claude, Cursor, etc.) to securely connect to external tools and data sources.

Think of it as a "plugin system" for AI:
- **Without MCP**: AI can only respond based on its training data
- **With MCP**: AI can call real APIs, query databases, trigger workflows

This server turns your backend APIs into MCP tools that AI assistants can use. When you ask the AI to "create a certificate order," it calls the tool defined in this server, which forwards the request to your backend API.

---

## Table of Contents

1. [What is MCP?](#what-is-mcp)
2. [Prerequisites](#prerequisites)
3. [Complete Local Walkthrough](#complete-local-walkthrough)
4. [Quick Start](#quick-start)
5. [How It Works](#how-it-works)
6. [Deploy to OpenShift](#deploy-to-openshift)
7. [Testing Your Server](#testing-your-server)
8. [Connecting to AI Assistants](#connecting-to-ai-assistants)
9. [Getting Tool Information from Your Backend API](#getting-tool-information-from-your-backend-api)
   - [Auto-Convert from OpenAPI Spec](#option-1-auto-convert-from-openapi-spec-recommended)
   - [Manual Definition from API Documentation](#option-2-manual-definition-from-api-documentation)
   - [Converting curl to mcpfile.yaml](#example-converting-a-curl-command-to-mcpfileyaml)
10. [Customization](#customization)
    - [Changing the Backend API URL](#changing-the-backend-api-url)
    - [Changing the Container Image](#changing-the-container-image)
    - [Adding a New Tool: Step-by-Step Guide](#adding-a-new-tool-step-by-step-guide)
    - [Adding a New API Base](#adding-a-new-api-base-for-different-endpoints)
    - [Common Patterns](#common-patterns)
11. [Troubleshooting](#troubleshooting)
12. [Reference](#reference)

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

# Install to PATH (choose one option):
# Option A: System-wide (requires sudo)
sudo mv genmcp /usr/local/bin/

# Option B: User-local (no sudo required)
mkdir -p ~/bin && mv genmcp ~/bin/
export PATH="$HOME/bin:$PATH"  # Add this to ~/.zshrc or ~/.bashrc

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

# Linux (Ubuntu/Debian)
sudo apt update && sudo apt install -y golang podman

# Linux (Fedora/RHEL)
sudo dnf install -y golang podman

# Windows
# Install Go from https://go.dev/dl/
# Install Podman Desktop from https://podman-desktop.io/

# Verify installations (all platforms)
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

# Install to PATH (choose one option):
# Option A: System-wide (requires sudo)
sudo mv genmcp /usr/local/bin/

# Option B: User-local (no sudo required)
mkdir -p ~/bin && mv genmcp ~/bin/
export PATH="$HOME/bin:$PATH"  # Add this to ~/.zshrc or ~/.bashrc

# Verify - should show "development" version, NOT "v0.1.x"
genmcp version
```

### Step 3: Clone This Repository

```bash
git clone https://github.com/rrbanda/custom-mcp-1.git
cd custom-mcp-1/cmp
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
{"level":"info","msg":"Loaded 2 tools from .../mcpfile.yaml"}
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

> **Note**: The MCP server uses Server-Sent Events (SSE) format. The response will be prefixed with `event: message` and `data:`. To extract just the JSON:
> ```bash
> curl -s -X POST http://localhost:8080/mcp \
>   -H "Content-Type: application/json" \
>   -d '{"jsonrpc":"2.0","method":"tools/list","id":1}' \
>   | grep "^data:" | sed 's/^data: //' | python3 -m json.tool
> ```

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
│ Step 3: Clone this repo, cd custom-mcp-1/cmp                    │
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
│   mcpfile.yaml → Defines tools:                                 │
│     • create_certificate_order (POST)                           │
│     • get_certificate_status (GET)                              │
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

## Connecting to AI Assistants

Once your MCP server is running, you can connect it to AI assistants like Cursor or Claude Desktop.

### Configure in Cursor

1. Open Cursor Settings (`Cmd+,` on macOS, `Ctrl+,` on Windows/Linux)
2. Navigate to **Features** → **MCP Servers**
3. Click **Add MCP Server**
4. Configure:
   - **Name**: `cmp-mcp` (or any name you prefer)
   - **Type**: `sse`
   - **URL**: `http://localhost:8080/mcp` (local) or `https://<your-route>/mcp` (OpenShift)

Alternatively, add to your `~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "cmp-mcp": {
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

### Configure in Claude Desktop

Add to your Claude Desktop config file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`  
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "cmp-mcp": {
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

Restart Claude Desktop after saving.

### Verify Connection

Once configured, you should see the `create_certificate_order` tool available in your AI assistant. Try asking:

> "What tools do you have for certificate management?"

The AI should respond with information about the certificate ordering tool.

---

## Getting Tool Information from Your Backend API

Before adding a new tool, you need to understand your backend API. GenMCP supports two approaches:

### Option 1: Auto-Convert from OpenAPI Spec (Recommended)

If your backend API has an OpenAPI/Swagger specification, GenMCP can automatically generate the mcpfile.yaml:

```bash
# From a URL
genmcp convert https://your-api.company.com/openapi.json

# From a local file
genmcp convert ./api-spec.yaml

# With custom output filename
genmcp convert ./api-spec.yaml -o my-tools.yaml
```

This generates a complete mcpfile.yaml with all endpoints as tools. You can then:
- Edit the generated file to keep only the tools you need
- Customize descriptions for better AI understanding
- Add authentication headers

### Option 2: Manual Definition from API Documentation

If you don't have an OpenAPI spec, gather this information from your API docs or curl examples:

| Information Needed | Where to Find It | Maps to mcpfile.yaml |
|--------------------|------------------|---------------------|
| API endpoint URL | API docs, curl examples | `invocationBases.*.http.url` |
| HTTP method | API docs (GET, POST, etc.) | `invocationBases.*.http.method` |
| Required headers | Auth docs, curl -H flags | `invocationBases.*.http.headers` |
| Request parameters | API docs, curl -d payload | `tools.*.inputSchema.properties` |
| Required vs optional | API docs | `tools.*.inputSchema.required` |
| Parameter types | API docs (string, number, array) | `tools.*.inputSchema.properties.*.type` |

### Example: Converting a curl Command to mcpfile.yaml

**Your curl command:**
```bash
curl -X POST "https://api.company.com/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "productId": "PROD-123",
    "quantity": 5,
    "notes": "Rush order"
  }'
```

**Translates to mcpfile.yaml:**
```yaml
invocationBases:
  baseOrders:
    http:
      method: POST
      url: https://api.company.com/orders
      headers:
        Content-Type: "application/json"
        Authorization: "Bearer {token}"  # {token} comes from tool input

tools:
  - name: create_order
    title: "Create Order"
    description: "Creates a new order for a product"
    inputSchema:
      type: object
      properties:
        token:
          type: string
          description: "Your API bearer token"
        productId:
          type: string
          description: "The product ID to order"
        quantity:
          type: integer
          description: "Number of items to order"
        notes:
          type: string
          description: "Optional order notes"
      required:
        - token
        - productId
        - quantity
    invocation:
      extends:
        from: baseOrders
```

### Key Mapping Rules

| curl/API | mcpfile.yaml | Notes |
|----------|--------------|-------|
| `-X POST` | `method: POST` | HTTP method |
| URL path | `url: https://...` | Base URL in invocationBases |
| `-H "Header: value"` | `headers: Header: "value"` | Static headers |
| `-H "Header: $VAR"` | `headers: Header: "{param}"` | Dynamic from input |
| `-d '{"field": ...}'` | `inputSchema.properties.field` | Request body fields |
| Path params `/items/{id}` | `url: "/items/{id}"` | `{id}` comes from inputSchema |

### Syncing mcpfile.yaml with configmap.yaml

**Important**: When you update `mcpfile.yaml`, you must also update `openshift/config/configmap.yaml` to match. The ConfigMap contains a copy of your tool definitions for Kubernetes deployment.

After editing mcpfile.yaml:
```bash
# Manually sync: Copy the tools section from mcpfile.yaml to configmap.yaml
# Then redeploy:
oc delete configmap cmp-mcp-config
oc apply -f openshift/config/configmap.yaml
oc rollout restart deployment/cmp-mcp-server
```

---

## Customization

This section explains how to modify the MCP server configuration, update backend API settings, and add new tools.

### Configuration Files Overview

| File | Purpose | When to Modify |
|------|---------|----------------|
| `mcpfile.yaml` | Defines tools, their inputs, and API endpoints | Adding/modifying tools |
| `mcpserver.yaml` | Server runtime settings (port, logging) | Changing port or log level |
| `openshift/config/deployment.yaml` | Kubernetes deployment | Changing image or backend URL |
| `openshift/config/configmap.yaml` | Tool definitions for OpenShift | After modifying mcpfile.yaml |

---

### Changing the Backend API URL

The backend API URL is configured in **two places** depending on your environment:

#### Local Development

Set the environment variable before running:

```bash
export CMP_BACKEND_URL="https://your-api-server.company.com"
genmcp run
```

#### OpenShift Deployment

Edit `openshift/config/deployment.yaml` line 44-45:

```yaml
        - name: CMP_BACKEND_URL
          value: "https://your-api-server.company.com"
```

Then redeploy:

```bash
oc apply -f openshift/config/deployment.yaml
oc rollout restart deployment/cmp-mcp-server
```

---

### Changing the Container Image

Edit `openshift/config/deployment.yaml` line 29:

```yaml
        image: quay.io/YOUR-ORG/your-image:v1.0.0
```

To build and push your own image:

```bash
genmcp build --tag quay.io/YOUR-ORG/your-image:v1.0.0 -f mcpfile.yaml --platform linux/amd64
podman push quay.io/YOUR-ORG/your-image:v1.0.0
```

---

### Adding a New Tool: Step-by-Step Guide

This walkthrough shows you how to add a new tool from scratch.

#### Step 1: Understand the Tool Structure

Each tool in `mcpfile.yaml` has these components:

```yaml
tools:
  - name: tool_name              # Unique identifier (snake_case)
    title: "Human Readable Name" # Display name
    description: |               # What the tool does (shown to AI)
      Detailed description of what this tool does.
      Include examples of when to use it.
    inputSchema:                 # JSON Schema for inputs
      type: object
      properties:
        param1:
          type: string
          description: "What this parameter is for"
        param2:
          type: integer
          description: "Another parameter"
      required: [param1, param2]
    invocation:                  # How to call the backend API
      extends:
        from: baseCertOrders     # Inherit from a base config
        extend:
          url: "/endpoint"       # Append to base URL
        override:
          method: GET            # Override HTTP method (optional)
```

#### Step 2: Add the Tool to mcpfile.yaml

Open `mcpfile.yaml` and add your tool to the `tools` array. Here's a complete example adding a "get certificate status" tool:

```yaml
  # Add this after the existing create_certificate_order tool
  - name: get_certificate_status
    title: "Get Certificate Status"
    description: |
      Retrieves the current status of a certificate order.
      
      Use this tool when you need to:
      - Check if a certificate order has been approved
      - Get the current processing state
      - Retrieve order details
      
      Returns the order status, creation date, and any error messages.
    inputSchema:
      type: object
      properties:
        # --- AUTH FIELDS (required for all tools) ---
        username:
          type: string
          description: "Your CMP username"
        token:
          type: string
          description: "Authentication token for the CMP API"
        smUserDev:
          type: string
          description: "SM_USER_DEV identifier"
        # --- TOOL-SPECIFIC FIELDS ---
        orderId:
          type: string
          description: "The order ID to check status for (e.g., ORD-12345)"
      required:
        - username
        - token
        - smUserDev
        - orderId
    invocation:
      extends:
        from: baseCertOrders
        extend:
          url: "/{orderId}/status"
        override:
          method: GET
```

#### Step 3: Test Locally

```bash
# Start the server
export CMP_BACKEND_URL="https://your-api.company.com"
genmcp run

# In another terminal, list tools to verify
curl -s -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/list","id":1}' \
  | grep -o '"name":"[^"]*"'

# Expected output should include:
# "name":"create_certificate_order"
# "name":"get_certificate_status"
```

#### Step 4: Update the ConfigMap for OpenShift

Copy your updated tool definitions to the ConfigMap:

```bash
# Option A: Manually edit openshift/config/configmap.yaml
# Copy the new tool from mcpfile.yaml into the configmap's mcpfile.yaml section

# Option B: Regenerate the entire ConfigMap (recommended)
# The ConfigMap contains a copy of mcpfile.yaml - keep them in sync!
```

> **Important**: The `mcpfile.yaml` in the repo root and the one embedded in 
> `configmap.yaml` must match! After editing `mcpfile.yaml`, update the ConfigMap.

#### Step 5: Deploy to OpenShift

```bash
# Delete old ConfigMap and apply new one
oc delete configmap cmp-mcp-config
oc apply -f openshift/config/configmap.yaml

# Restart the deployment to pick up changes
oc rollout restart deployment/cmp-mcp-server

# Wait for rollout to complete
oc rollout status deployment/cmp-mcp-server

# Verify the new tool appears
ROUTE=$(oc get route cmp-mcp-server -o jsonpath='{.spec.host}')
curl -s -X POST "https://${ROUTE}/mcp" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/list","id":1}' \
  | grep -o '"name":"[^"]*"'
```

#### Step 6: Test in Cursor

After deploying, refresh MCP servers in Cursor:
1. Open Command Palette (`Cmd+Shift+P`)
2. Run "MCP: Refresh Servers"
3. Ask the AI: "What tools do you have for checking certificate status?"

---

### Adding a New API Base (for Different Endpoints)

If your new tool uses a different API path structure, add a new base in `mcpfile.yaml`:

```yaml
invocationBases:
  # Existing base for certificate orders
  baseCertOrders:
    http:
      method: POST
      url: ${CMP_BACKEND_URL}/api/orderextws/public/Orders
      headers:
        Content-Type: "application/json"
        Username: "{username}"
        Token: "{token}"
        SM_USER_DEV: "{smUserDev}"

  # NEW: Base for inventory API
  baseInventory:
    http:
      method: GET
      url: ${CMP_BACKEND_URL}/api/inventory/v2
      headers:
        Content-Type: "application/json"
        Username: "{username}"
        Token: "{token}"
        SM_USER_DEV: "{smUserDev}"
```

Then reference it in your tool:

```yaml
  - name: list_inventory
    title: "List Inventory"
    description: "Lists available inventory items"
    inputSchema:
      # ... your schema ...
    invocation:
      extends:
        from: baseInventory  # Use the new base
        extend:
          url: "/items"
```

---

### Common Patterns

#### Pattern 1: GET Request with Path Parameter

```yaml
  - name: get_order_by_id
    # ...
    invocation:
      extends:
        from: baseCertOrders
        extend:
          url: "/{orderId}"
        override:
          method: GET
```

#### Pattern 2: POST Request with JSON Body

```yaml
  - name: create_order
    # ...
    invocation:
      extends:
        from: baseCertOrders
        extend:
          url: "/create"
      # method: POST is inherited from base
```

#### Pattern 3: DELETE Request

```yaml
  - name: cancel_order
    # ...
    invocation:
      extends:
        from: baseCertOrders
        extend:
          url: "/{orderId}"
        override:
          method: DELETE
```

---

### Quick Reference: What to Update Where

| Change | File(s) to Update | Commands to Run |
|--------|-------------------|-----------------|
| Add new tool | `mcpfile.yaml`, `configmap.yaml` | `oc delete configmap...`, `oc apply...`, `oc rollout restart...` |
| Change backend URL | `deployment.yaml` (line 45) | `oc apply -f deployment.yaml`, `oc rollout restart...` |
| Change container image | `deployment.yaml` (line 29) | `oc apply -f deployment.yaml` |
| Change server port | `mcpserver.yaml`, `deployment.yaml`, `service.yaml` | Rebuild image + redeploy all |
| Add new API base | `mcpfile.yaml`, `configmap.yaml` | Same as "Add new tool" |

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
