# Custom MCP Servers

This repository contains custom MCP (Model Context Protocol) server configurations built with [GenMCP](https://github.com/genmcp/gen-mcp). These configurations expose backend APIs as tools for AI assistants - **zero code required**.

---

## Available MCP Servers

| Server | Description | Directory |
|--------|-------------|-----------|
| **CMP** | Certificate Management Platform tools | [`cmp/`](./cmp/) |

---

## Prerequisites

Before using any MCP server in this repo, you need:

| Requirement | Purpose | Installation |
|-------------|---------|--------------|
| **Go 1.21+** | Build GenMCP CLI | `brew install go` (macOS) |
| **GenMCP CLI** | Run MCP servers | See below |
| **podman or docker** | Build container images | `brew install podman` |
| **oc CLI** | Deploy to OpenShift (optional) | [Download](https://mirror.openshift.com/pub/openshift-v4/clients/ocp/) |

### Install GenMCP CLI

```bash
# Clone GenMCP repository
git clone https://github.com/genmcp/gen-mcp.git
cd gen-mcp

# Build the CLI
make build-cli

# Move to a directory in your PATH
sudo mv genmcp /usr/local/bin/

# Verify installation
genmcp version
```

---

## Quick Start

### Run CMP MCP Server Locally

```bash
# Navigate to the CMP server
cd cmp

# Set your backend URL
export CMP_BACKEND_URL="https://your-cmp-api.company.com"

# Start the server
genmcp run

# Server is now running at http://localhost:8080/mcp
```

### Test the Server

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/list","id":1}'
```

---

## Documentation

For detailed documentation on each MCP server, see:

- **CMP Server**: [`cmp/README.md`](./cmp/README.md)

---

## Adding New MCP Servers

To add a new MCP server to this repository:

1. Create a new directory (e.g., `my-server/`)
2. Add `mcpfile.yaml` (tool definitions)
3. Add `mcpserver.yaml` (runtime config)
4. Add `openshift/config/` directory with deployment manifests
5. Add `README.md` with documentation

See the [`cmp/`](./cmp/) directory as a template.

---

## License

Apache 2.0 - See [LICENSE](./LICENSE) for details.
