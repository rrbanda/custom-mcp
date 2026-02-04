# Marketplace Order API - MCP Server

Creates orders in the Marketplace API (SecureSSO Policy Modify, etc.) via MCP.

---

## Quick Test with Mock Backend

### Step 1: Start the mock server

```bash
cd marketplace-order
go run mock_server.go
```

You should see:
```
Mock Marketplace API running on http://localhost:9999 (STRICT array-only mode)
```

### Step 2: Start genmcp (in a new terminal)

```bash
cd marketplace-order
genmcp run -f mcpfile-test.yaml -s mcpserver.yaml
```

You should see:
```
Starting MCP server on port 3000
```

### Step 3: Test with curl (simple)

```bash
curl -s http://localhost:3000/mcp -X POST -H "Content-Type: application/json" -d '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "create_order",
    "arguments": {
      "orders": [{
        "orderByGeid": "1011035821",
        "orderForGeid": "1011035821",
        "productId": "57097_20447_GLOBAL",
        "productName": "SecureSSO - Policy Modify",
        "questionnaire": [
          {"dataKey": "57097_20447_9818373", "key": "env_var", "IsHidden": "true", "value": ["DEV"]}
        ],
        "orderForUsers": [
          {"locationIdOrderFor": "37750", "geoIdOrderFor": "USA", "orderForGeid": "1011035821"}
        ]
      }]
    }
  }
}'
```

Expected response:
```json
{"bodyFormat":"array","message":"Order created successfully","orderId":"ORD-2024-001234","orders":1,"success":true}
```

### Step 4: Test with full payload

```bash
curl -s http://localhost:3000/mcp -X POST -H "Content-Type: application/json" -d '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "create_order",
    "arguments": {
      "orders": [{
        "orderByGeid": "1011035821",
        "orderForGeid": "1011035821",
        "productId": "57097_20447_GLOBAL",
        "productName": "SecureSSO - Policy Modify",
        "questionnaire": [
          {"dataKey": "57097_20447_9818373", "key": "env_var", "IsHidden": "true", "value": ["DEV"]},
          {"dataKey": "37215_17629_31657", "key": "ENV", "IsHidden": "false", "value": ["DEV"]},
          {"dataKey": "48317_32822_967", "key": "appId", "IsHidden": "false", "value": ["158840||OneReset - 158840"], "searchInput": "158840", "autoPopulated": false},
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
          {"dataKey": "48317_32822_851110", "key": "policyList", "IsHidden": "false", "value": ["158840-OneReset||158840-OneReset-"], "searchInput": "", "autoPopulated": true},
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
          {"locationIdOrderFor": "37750", "geoIdOrderFor": "USA", "orderForGeid": "1011035821"}
        ]
      }]
    }
  }
}'
```

---

## Use with Real Backend

### Step 1: Edit mcpfile.yaml

Update the auth headers (lines 77-80):

```yaml
          Username: "YOUR_USERNAME_HERE"
          Token: "YOUR_TOKEN_HERE"
          SM_USER_DEV: "YOUR_SM_USER_HERE"
          Authorization: "YOUR_AUTH_HERE"
```

### Step 2: Run genmcp

```bash
cd marketplace-order
genmcp run
```

### Step 3: Connect AI Assistant

**Cursor**: Settings → Features → MCP Servers → Add → URL: `http://localhost:3000/mcp`

---

## How It Works

The `bodyRoot: orders` setting extracts the array from the MCP input:

| MCP Input | HTTP Body Sent to Backend |
|-----------|---------------------------|
| `{"orders": [{...}]}` | `[{...}]` |

This matches the backend's expected format exactly.

---

## Files

| File | Purpose |
|------|---------|
| `mcpfile.yaml` | Real API config (edit auth values) |
| `mcpfile-test.yaml` | Mock server config (localhost:9999) |
| `mcpserver.yaml` | Server config (port 3000) |
| `mock_server.go` | Mock backend for testing |
