# Marketplace Order API - GenMCP Example

This example wraps the Marketplace Order API as an MCP server, allowing AI assistants to create orders.

## Key Feature: `bodyRoot` for Array Bodies

The Marketplace API expects a JSON array `[...]` as the request body. Since MCP requires `inputSchema` to be an object, we use the `bodyRoot` configuration to extract the `orders` property and send just the array:

```yaml
invocation:
  http:
    method: POST
    url: https://marketplace.dev.citigroup.net/api/orderextws/public/Orders/createOrder
    bodyRoot: orders  # <-- Extracts just the 'orders' array as the HTTP body
```

This ensures the backend receives `[{...}]` instead of `{"orders": [{...}]}`.

## Setup

### 1. Set Environment Variables

The API requires authentication headers. Set these environment variables before running:

```bash
export MARKETPLACE_USERNAME="your-username"
export MARKETPLACE_TOKEN="your-api-token"
export MARKETPLACE_SM_USER="your-sm-user-dev"
export MARKETPLACE_AUTH="Bearer your-auth-token"
```

### 2. Run the MCP Server

```bash
cd marketplace-order
genmcp run
```

The server will start on `http://localhost:8080/mcp`.

### 3. Connect to AI Assistants

**Cursor**: Settings → Features → MCP Servers → Add Server → URL: `http://localhost:8080/mcp`

**Claude Desktop**: Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "marketplace-order": {
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

## Tool: create_order

Creates a new order in the Marketplace API.

### MCP Input Format

When calling via MCP, wrap the orders in an `orders` property:

```json
{
  "orders": [
    {
      "orderByGeid": "1011035821",
      "orderForGeid": "1011035821", 
      "productId": "57097_20447_GLOBAL",
      "productName": "SecureSSO - Policy Modify",
      "questionnaire": [...],
      "orderForUsers": [...]
    }
  ]
}
```

### HTTP Body Sent to Backend

GenMCP extracts the `orders` array and sends it directly:

```json
[
  {
    "orderByGeid": "1011035821",
    "orderForGeid": "1011035821",
    "productId": "57097_20447_GLOBAL",
    "productName": "SecureSSO - Policy Modify",
    "questionnaire": [...],
    "orderForUsers": [...]
  }
]
```

### Example: SecureSSO Policy Modify Order

Here's a complete example for calling the tool via MCP:

```json
{
  "orders": [
    {
      "orderByGeid": "1011035821",
      "orderForGeid": "1011035821",
      "productId": "57097_20447_GLOBAL",
      "productName": "SecureSSO - Policy Modify",
      "questionnaire": [
        {"dataKey": "57097_20447_9818373", "key": "env_var", "IsHidden": "true", "value": ["DEV"]},
        {"dataKey": "37215_17629_31657", "key": "ENV", "IsHidden": "false", "value": ["DEV"]},
        {"dataKey": "48317_32822_967", "key": "appId", "IsHidden": "false", "value": ["158840||OneReset - 158840"]},
        {"dataKey": "48317_32822_693238", "key": "appName", "IsHidden": "false", "value": ["OneReset"]},
        {"dataKey": "48317_32822_37068", "key": "applicationManager", "IsHidden": "false", "value": ["1000311258"]},
        {"dataKey": "48317_32822_27849", "key": "ProjectManagerApprover", "IsHidden": "true", "value": ["1000311258"]},
        {"dataKey": "48317_32822_865829", "key": "applicationStatus", "IsHidden": "false", "value": ["Production"]},
        {"dataKey": "48317_32822_2642", "key": "Business Sector", "IsHidden": "false", "value": ["CISO"]},
        {"dataKey": "48317_32822_879827", "key": "businessOwnerGEID", "IsHidden": "false", "value": ["1000176492"]},
        {"dataKey": "48317_32822_741426", "key": "applicationISRisk", "IsHidden": "false", "value": ["MEDIUM"]},
        {"dataKey": "48317_32822_734231", "key": "applicationBusinessCriticality", "IsHidden": "false", "value": ["No"]},
        {"dataKey": "48317_32822_751911", "key": "applicationSecurityClassification", "IsHidden": "false", "value": ["Confidential"]},
        {"dataKey": "66705_11305_12971", "key": "sox", "IsHidden": "false", "value": ["Not SOX-critical"]},
        {"dataKey": "57097_20447_64263", "key": "OperatingEnv", "IsHidden": "false", "value": ["Extranet"]},
        {"dataKey": "57097_20447_15299", "key": "MFARequired", "IsHidden": "false", "value": ["No"]},
        {"dataKey": "37215_17629_10607", "key": "hostingModel", "IsHidden": "false", "value": ["Internal-CTI"]},
        {"dataKey": "48317_32822_16957", "key": "location", "IsHidden": "false", "value": ["NA"]},
        {"dataKey": "48317_32822_851110", "key": "policyList", "IsHidden": "false", "value": ["158840-OneReset||158840-OneReset-"]},
        {"dataKey": "48317_32822_947177", "key": "Please describe modifications needed", "IsHidden": "false", "value": ["158840-OneReset-Agent"]},
        {"dataKey": "48317_32822_711399", "key": "requestType", "IsHidden": "true", "value": ["NON-WEBBANK"]},
        {"dataKey": "48317_32822_756456", "key": "requestEnv", "IsHidden": "true", "value": ["%requestEnv%"]},
        {"dataKey": "48317_32822_880551", "key": "comments", "IsHidden": "false", "value": ["158840-OneReset-Agent"]},
        {"dataKey": "48317_32822_920901", "key": "CMPEnv", "IsHidden": "true", "value": ["%CMPEnv%"]},
        {"dataKey": "81574_11064_50255", "key": "SSOApproverHiddenflag", "IsHidden": "true", "value": ["true"]},
        {"dataKey": "66705_11305_76706", "key": "p1", "IsHidden": "true", "value": ["true"]},
        {"dataKey": "65939_71351_54343", "key": "I acknowledge that, this request is for application that is meant to be used by internal Citi users only and is available only on the intranet.", "IsHidden": "false", "value": ["Acknowledged"]},
        {"dataKey": "57097_20447_2137703", "key": "compliance", "IsHidden": "false", "value": ["I Acknowledge"]}
      ],
      "orderForUsers": [
        {
          "locationIdOrderFor": "37750",
          "geoIdOrderFor": "USA",
          "orderForGeid": "1011035821"
        }
      ]
    }
  ]
}
```

## Testing with curl

You can test the MCP server directly:

```bash
# Initialize session
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test"}}}'

# List available tools
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "id": 2, "method": "tools/list"}'

# Call create_order
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Mcp-Session-Id: test" \
  -d '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "create_order", "arguments": {"orders": [{"orderByGeid": "123", "orderForGeid": "123", "productId": "test", "productName": "Test", "questionnaire": [], "orderForUsers": []}]}}}'
```

## Files

- `mcpfile.yaml` - Tool definitions (uses env vars for auth)
- `mcpserver.yaml` - Server runtime configuration
- `README.md` - This documentation
