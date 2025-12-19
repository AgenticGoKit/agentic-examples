package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kunalkushwaha/agenticgokit/core/vnext"
)

// Config holds application configuration
type Config struct {
	APIKey   string
	Port     string
	Provider string // e.g., "huggingface"
	Model    string // e.g., "Qwen/Qwen2.5-72B-Instruct"
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	// Check for API key from LLM_API_KEY environment variable
	apiKey := os.Getenv("LLM_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("LLM_API_KEY environment variable not set\nPlease enable LLM_API_KEY in .env")
	}

	// Require LLM_PROVIDER to be explicitly set
	provider := os.Getenv("LLM_PROVIDER")
	if provider == "" {
		return nil, fmt.Errorf("LLM_PROVIDER environment variable not set\nPlease set it with: $env:LLM_PROVIDER=\"your-provider\" (e.g., \"openrouter\", \"huggingface\")")
	}

	// Require LLM_MODEL to be explicitly set
	model := os.Getenv("LLM_MODEL")
	if model == "" {
		return nil, fmt.Errorf("LLM_MODEL environment variable not set\nPlease set it with: $env:LLM_MODEL=\"your-model\" (e.g., \"openai/gpt-4o-mini\", \"Qwen/Qwen2.5-72B-Instruct\")")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		APIKey:   apiKey,
		Port:     port,
		Provider: provider,
		Model:    model,
	}, nil
}

// ValidateAPIConnection verifies the API key works by making a test request
func ValidateAPIConnection(apiKey string, provider string, model string) error {
	log.Println("🔍 Validating API connection...")

	testAgent, err := vnext.QuickChatAgentWithConfig("ValidationTest", &vnext.Config{
		Name:    "validation_test",
		Timeout: 15 * time.Second,
		LLM: vnext.LLMConfig{
			Provider: provider,
			Model:    model,
			APIKey:   apiKey,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create test agent: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err = testAgent.Run(ctx, "test")
	if err != nil {
		return fmt.Errorf("API connection test failed: %w", err)
	}

	return nil
}
