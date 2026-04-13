package main

import (
	"context"
	"fmt"
	"os"
	"time"

	agk "github.com/agenticgokit/agenticgokit/v1beta"

	// Import MCP plugins - these enable MCP server connectivity
	_ "github.com/agenticgokit/agenticgokit/plugins/mcp/default"
	_ "github.com/agenticgokit/agenticgokit/plugins/mcp/unified"
)

func main() {
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println("  MCP Integration Example - AgenticGoKit")
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println()

	// Step 1: Check for authentication token
	authToken := os.Getenv("MCP_GATEWAY_AUTH_TOKEN")
	/*	if authToken == "" {
			fmt.Println("❌ MCP_GATEWAY_AUTH_TOKEN not set")
			fmt.Println()
			fmt.Println("Please set your auth token:")
			fmt.Println("   PowerShell: $env:MCP_GATEWAY_AUTH_TOKEN='your-token-here'")
			fmt.Println("   Bash: export MCP_GATEWAY_AUTH_TOKEN='your-token-here'")
			fmt.Println()
			fmt.Println("Get your token by running:")
			fmt.Println("   docker mcp gateway run --port 8012 --transport sse")
			return
		}
	*/
	fmt.Printf("✓ Auth token configured (%d chars)\n\n", len(authToken))

	// Step 2: Configure the MCP server
	// This tells AgenticGoKit how to connect to your MCP server
	mcpServer := agk.MCPServer{
		Name:    "docker-mcp-gateway-streaming",
		Type:    "http_streaming",
		Address: "http://localhost:3000/mcp",
		Enabled: true,
	}

	fmt.Println("📡 MCP Server Configuration:")
	fmt.Printf("   Name: %s\n", mcpServer.Name)
	fmt.Printf("   Type: %s\n", mcpServer.Type)
	fmt.Printf("   Endpoint: %s\n", mcpServer.Address)
	fmt.Println()

	// Step 3: Create an agent with MCP tools
	ctx := context.Background()

	fmt.Println("🤖 Creating agent with MCP tools...")
	agent, err := agk.NewBuilder("mcp-example").
		WithConfig(&agk.Config{
			Name: "mcp-agent",

			// System prompt: guides the agent on how to use tools
			SystemPrompt: `You are a helpful AI assistant with access to MCP tools.

When users ask for current information (time, weather, searches):
- Use the appropriate tool from the available MCP tools
- Provide clear, natural responses based on the tool results
- If a tool returns an error, explain it to the user

Available tool categories:
- Time/timezone tools (get_current_time, convert_time)
- Web search (search)
- Content fetching (fetch_content)
- And more (see discovered tools)`,

			Timeout: 60 * time.Second,

			// Set DebugMode to true to see detailed MCP protocol exchange
			DebugMode: false,

			LLM: agk.LLMConfig{
				Provider:    "ollama",
				Model:       "qwen3.5:cloud", // Use a model good at function calling
				Temperature: 0.2,             // Lower temp = more consistent tool use
				MaxTokens:   800,
			},
		}).
		// This is where the magic happens - register MCP tools!
		WithTools(agk.WithMCP(mcpServer), agk.WithReasoningConfig(10, true)).
		WithObservability("mcp-example", "1.0").
		Build()

	if err != nil {
		fmt.Printf("❌ Failed to create agent: %v\n", err)
		return
	}
	defer agent.Cleanup(ctx)
	fmt.Println("   ✓ Agent created successfully\n")

	// Step 4: Get the MCP manager and check server health
	mgr := agk.GetMCPManager()
	if mgr == nil {
		fmt.Println("❌ MCP manager not initialized")
		return
	}

	fmt.Println("🏥 Checking server health...")
	health := mgr.HealthCheck(ctx)

	allHealthy := true
	for name, status := range health {
		if status.Status == "healthy" {
			fmt.Printf("   ✓ %s is healthy\n", name)
		} else {
			fmt.Printf("   ✗ %s: %s", name, status.Status)
			if status.Error != "" {
				fmt.Printf(" - %s", status.Error)
			}
			fmt.Println()
			allHealthy = false
		}
	}

	if !allHealthy {
		fmt.Println("\n⚠️  Some servers are unhealthy. Continuing anyway...")
	}
	fmt.Println()

	// Step 5: Discover available tools from MCP servers
	fmt.Println("🔍 Discovering tools...")
	tools := mgr.GetAvailableTools()

	if len(tools) == 0 {
		fmt.Println("   ⚠️  No tools discovered")
		fmt.Println()
		fmt.Println("Troubleshooting:")
		fmt.Println("   • Is the MCP server running?")
		fmt.Println("   • Is the auth token correct?")
		fmt.Println("   • Check the server logs for errors")
		fmt.Println("   • Try enabling DebugMode: true to see details")
		return
	}

	fmt.Printf("   ✓ Found %d tools\n\n", len(tools))

	// Display all discovered tools with their details
	fmt.Println("📋 Available Tools:")
	for i, tool := range tools {
		fmt.Printf("   %2d. %s\n", i+1, tool.Name)
		if tool.Description != "" {
			// Limit description length for cleaner display
			desc := tool.Description
			if len(desc) > 100 {
				desc = desc[:97] + "..."
			}
			fmt.Printf("       Description: %s\n", desc)
		}
		fmt.Printf("       Server: %s\n", tool.ServerName)
	}

	fmt.Println()
	fmt.Println("✓ Discovery complete")

	// Step 6: Test the agent with a sample query
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println("  Testing Agent with MCP Tools")
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println()

	// Example query that should use the get_current_time tool
	query := "please fetch url www.plankeya.com"
	fmt.Printf("💬 Query: %s\n\n", query)

	result, err := agent.Run(ctx, query)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Println("📊 Result:")
		fmt.Printf("   %s\n", result.Content)
		fmt.Printf("   ⏱️  Duration: %v\n", result.Duration)

		if len(result.ToolsCalled) > 0 {
			fmt.Printf("   🔧 Tools used: %v\n", result.ToolsCalled)
		} else {
			fmt.Println("   ⚠️  No tools were called")
			fmt.Println("   Tip: Try a more explicit query or check LLM model capabilities")
		}
	}

	fmt.Println()
	fmt.Println("✓ Example complete!")
	fmt.Println()
	fmt.Println("Try modifying the query above to test other tools:")
	fmt.Println("  - \"What time is it in Tokyo?\"")
	fmt.Println("  - \"Search for AgenticGoKit documentation\"")
	fmt.Println("  - \"Convert 3pm EST to Tokyo time\"")
}
