# MCP Integration Example

Shows how to connect an AgenticGoKit agent to an MCP (Model Context Protocol) server for tool discovery and execution.

## Prerequisites

- Go 1.21+
- Ollama running locally
- Docker MCP Gateway running on port 8012

```bash
# Start the gateway
docker mcp gateway run --port 8012 --transport sse
# Note the auth token printed in the output
```

## Quick Start

```powershell
# Set auth token (PowerShell)
$env:MCP_GATEWAY_AUTH_TOKEN='your-token-here'

# Run
cd examples/mcp-example
go run main.go
```

## How It Works

### 1. Import MCP plugins

```go
import (
    agk "github.com/agenticgokit/agenticgokit/v1beta"

    _ "github.com/agenticgokit/agenticgokit/plugins/mcp/default"
    _ "github.com/agenticgokit/agenticgokit/plugins/mcp/unified"
)
```

The blank imports register MCP transport support via `init()`  no further setup needed.

### 2. Configure the MCP server

```go
mcpServer := agk.MCPServer{
    Name:    "docker-mcp-gateway",
    Type:    "http_sse",   // or: websocket, stdio, tcp
    Address: "localhost",
    Port:    8012,
    Enabled: true,
}
```

### 3. Create an agent with MCP tools

```go
agent, err := agk.NewBuilder("mcp-example").
    WithConfig(&agk.Config{
        Name:         "mcp-agent",
        SystemPrompt: "You are a helpful assistant with access to MCP tools.",
        Timeout:      60 * time.Second,
        LLM: agk.LLMConfig{
            Provider:    "ollama",
            Model:       "qwen3.5:cloud",
            Temperature: 0.2,
            MaxTokens:   800,
        },
    }).
    WithTools(agk.WithMCP(mcpServer)).
    Build()
```

### 4. Run a query

```go
result, err := agent.Run(ctx, "What time is it in UTC right now?")
fmt.Println(result.Content)
fmt.Println("Tools used:", result.ToolsCalled)
```

The agent automatically discovers the available tools, picks the right one (`get_current_time`), calls it with the correct parameters, and returns the result.

## Reasoning Mode

By default the agent returns the raw tool output. Enable reasoning to have the LLM synthesize a natural-language response:

```go
WithTools(
    agk.WithMCP(mcpServer),
    agk.WithReasoningConfig(2, true),
)
```

| Parameter | Description |
|-----------|-------------|
| `maxIterations` (e.g. `2`) | How many LLM  tool cycles to allow |
| `continueOnToolUse` (`true`) | When `true`, the LLM always makes a second pass after tool execution to produce a natural sentence instead of raw data |

**Without reasoning:**
```
[{type:text text:{datetime:2026-04-07T15:30:00Z ...}}]
```

**With reasoning:**
```
The current time in UTC is 3:30 PM on April 7, 2026.
```

> **Note:** Reasoning adds ~515s latency due to the extra LLM call.

## Inspecting Discovered Tools

```go
mgr := agk.GetMCPManager()

// Check health
health := mgr.HealthCheck(ctx)

// List tools
for _, tool := range mgr.GetAvailableTools() {
    fmt.Printf("%s (%s): %s\n", tool.Name, tool.ServerName, tool.Description)
}
```

## Debugging

Set `DebugMode: true` in the config or `AGENTICGOKIT_LOG_LEVEL=debug` to see the full MCP protocol exchange:

```
[MCP] Tool: get_current_time - Get current time in a specific timezone
[SSE] Request JSON: {"method":"tools/call","params":{"name":"get_current_time","arguments":{"timezone":"UTC"}}}
[SSE] Response: {"result":{"content":[{"type":"text","text":"{...}"}]}}
```

## Troubleshooting

| Symptom | Likely cause | Fix |
|---------|-------------|-----|
| `No tools discovered` | Auth token wrong or server not running | Check `MCP_GATEWAY_AUTH_TOKEN` and `docker ps` |
| Tool called with wrong args | Small model guessing parameters | Use `qwen3.5:cloud` or `llama3.2:3b` |
| `context deadline exceeded` | Slow server or network | Increase `Timeout` in config |

## Multiple MCP Servers

Pass multiple servers  all their tools are merged and available to the agent:

```go
WithTools(
    agk.WithMCP(gatewayServer),
    agk.WithMCP(localServer),
)
```
