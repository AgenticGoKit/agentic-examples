# AgenticGoKit Examples & Demos

> **Production-ready examples showcasing AgenticGoKit's powerful agent orchestration capabilities**

This repository contains practical examples demonstrating how to build AI agent applications using [AgenticGoKit](https://github.com/kunalkushwaha/agenticgokit) - a Go framework for creating, orchestrating, and managing AI agent workflows.

## 📚 Examples

### Story Writer Chat v2

**[📖 Full Documentation](story-writer-chat-v2/README.md)** | **Path**: `story-writer-chat-v2/`

A collaborative AI story writing application with 3-agent workflow, real-time streaming, and React UI.

**What it demonstrates:**
- SubWorkflow composition and nested agent execution
- Conditional loop termination with dynamic exit conditions
- Real-time WebSocket streaming with lifecycle events
- Agent memory and context retrieval
- Modern React frontend with agent status tracking

**Tech Stack:** Go, AgenticGoKit, WebSocket, React + TypeScript, TailwindCSS

**Quick Start:**
```bash
cd story-writer-chat-v2
export LLM_PROVIDER=huggingface LLM_MODEL=Qwen/Qwen2.5-72B-Instruct LLM_API_KEY=your-key
go run main.go  # Backend on :8080
cd frontend && npm install && npm run dev  # Frontend on :5173
```

---

*More examples coming soon! Check back for additional patterns and use cases.*

## 🚀 Getting Started

### Prerequisites
- **Go 1.23+** - [Install Go](https://golang.org/dl/)
- **LLM API Key** - Get one from [HuggingFace](https://huggingface.co/settings/tokens), [OpenRouter](https://openrouter.ai/keys), or [OpenAI](https://platform.openai.com/api-keys)
- **Node.js 18+** (for frontend examples) - [Install Node.js](https://nodejs.org/)

### Basic Setup

```bash
# Clone the repository
git clone https://github.com/AgenticGoKit/agentic-examples.git
cd agentic-examples

# Navigate to an example
cd story-writer-chat-v2

# Set up environment variables
export LLM_PROVIDER=huggingface
export LLM_MODEL=Qwen/Qwen2.5-72B-Instruct
export LLM_API_KEY=your-api-key

# Run the example (see individual README for details)
go run main.go
```

**Windows PowerShell:**
```powershell
$env:LLM_PROVIDER="huggingface"
$env:LLM_MODEL="Qwen/Qwen2.5-72B-Instruct"
$env:LLM_API_KEY="your-api-key"
```

Each example has its own detailed README with specific setup instructions and architecture documentation.

---

## 🤝 Contributing

We welcome new examples! To contribute:

1. Create a new directory for your example
2. Include a comprehensive README with:
   - Overview and what it demonstrates
   - Quick start guide
   - Architecture details (if complex)
   - Troubleshooting section
3. Add your example to the catalog above with a brief description
4. Submit a Pull Request

See individual examples for code organization patterns and best practices.

---

**Built with ❤️ by the AgenticGoKit community**

*Empowering developers to build sophisticated AI agent applications in Go*
