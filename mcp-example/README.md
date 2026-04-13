# MCP Integration Example

Shows how to connect an AgenticGoKit agent to an MCP (Model Context Protocol) server for tool discovery and execution. This example uses [`fetcher-mcp`](https://github.com/jae-jae/fetcher-mcp), a lightweight MCP server that exposes web-fetch tools.

## Prerequisites

- Go 1.21+
- Ollama running locally with a function-calling model (e.g. `qwen3.5:cloud`)
- `fetcher-mcp` running in Docker:

```bash
docker run -p 3000:3000 ghcr.io/jae-jae/fetcher-mcp:latest
```

## Quick Start

```powershell
# Run (PowerShell)
go run .\main.go
```

```bash
# Run (Bash)
go run main.go
```

### Authentication (Optional)

Some MCP servers require an auth token. Set it via the environment variable before running:

```powershell
$env:MCP_GATEWAY_AUTH_TOKEN = 'your-token-here'
```

```bash
export MCP_GATEWAY_AUTH_TOKEN='your-token-here'
```

When set, the token is passed as a Bearer header on every MCP request. If the server does not require auth, the variable can be left unset.

## How It Works

### 1. Import MCP plugins

```go
import (
    agk "github.com/agenticgokit/agenticgokit/v1beta"

    _ "github.com/agenticgokit/agenticgokit/plugins/mcp/default"
    _ "github.com/agenticgokit/agenticgokit/plugins/mcp/unified"
)
```

The blank imports register MCP transport support via `init()`, no further setup needed.

### 2. Configure the MCP server

```go
mcpServer := agk.MCPServer{
    Name:    "docker-mcp-gateway-streaming",
    Type:    "http_streaming",   // or: http_sse, websocket, stdio
    Address: "http://localhost:3000/mcp",
    Enabled: true,
}
```

The `Address` field holds the full endpoint URL for `http_streaming` and `http_sse` transports.

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
    WithTools(agk.WithMCP(mcpServer), agk.WithReasoningConfig(10, true)).
    WithObservability("mcp-example", "1.0").
    Build()
```

`WithObservability` enables built-in OpenTelemetry tracing. See the [Observability](#observability) section below.

### 4. Run a query

```go
result, err := agent.Run(ctx, "please fetch url www.example.com")
fmt.Println(result.Content)
fmt.Println("Tools used:", result.ToolsCalled)
```

The agent automatically discovers the available tools from the MCP server, picks the right one (`fetch_url`), calls it with the correct parameters, and returns the result.

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
| `maxIterations` (e.g. `2`) | How many LLM/tool cycles to allow |
| `continueOnToolUse` (`true`) | When `true`, the LLM always makes a second pass after tool execution to produce a natural sentence instead of raw data |

**Without reasoning:**
```
[{type:text text:{datetime:2026-04-07T15:30:00Z ...}}]
```

**With reasoning:**
```
The current time in UTC is 3:30 PM on April 7, 2026.
```

> **Note:** Reasoning adds extra latency due to the second LLM call. Keep `maxIterations` low (1-2) for simple tool use.

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
[MCP] Tool: fetch_url
[Streaming] Request JSON: {"method":"tools/call","params":{"name":"fetch_url","arguments":{"url":"http://www.example.com"}}}
[Streaming] Response: {"result":{"content":[{"type":"text","text":"..."}]}}
```

## Observability

Enable tracing to capture every LLM call and MCP tool execution:

```powershell
$env:AGK_TRACE = "true"
go run .\main.go
```

Traces are saved to `.agk/runs/<run-id>/trace.jsonl`. View them with the `agk` CLI:

```bash
# Install agk CLI
go install github.com/agenticgokit/agk@latest

# List runs
agk trace list

# Show full span tree
agk trace show <run-id>
```

To capture tool input/output content in the trace:

```powershell
$env:AGK_TRACE_LEVEL = "detailed"
$env:AGK_TRACE       = "true"
go run .\main.go
```

For a full walkthrough, see [blog-observability.md](./blog-observability.md).

## Troubleshooting

| Symptom | Likely cause | Fix |
|---------|-------------|-----|
| `No tools discovered` | MCP server not running | Check `docker ps` and confirm port 3000 is reachable |
| Auth token required error | Server needs auth | Set `MCP_GATEWAY_AUTH_TOKEN` |
| Tool called with wrong args | Small model guessing parameters | Use `qwen3.5:cloud` or `llama3.2:3b` |
| `context deadline exceeded` | Slow server or network | Increase `Timeout` in config |

## Multiple MCP Servers

Pass multiple servers to the `WithTools` call. All their tools are merged and available to the agent:

```go
WithTools(
    agk.WithMCP(fetcherServer),
    agk.WithMCP(localServer),
    agk.WithReasoningConfig(10, true),
)
```
