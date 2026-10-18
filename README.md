# TarsDify

MCP integration for Dify (LLM app development platform).

## Tools

- `dify_chat` — Send a chat message to a Dify app
- `dify_workflow` — Run a Dify workflow
- `dify_knowledge_retrieve` — Retrieve from Dify knowledge base
- `dify_files` — Upload a file to Dify
- `dify_conversations` — List conversations
- `dify_health` — Check Dify instance health
- `dify_status` — Integration status

## Quick Start

```bash
go build -o tarsdify
TARS_MCP_ENDPOINT=https://api.example.com TARS_MCP_API_KEY=xxx ./tarsdify
```

## Environment Variables

- `TARS_MCP_ENDPOINT` — upstream API endpoint
- `TARS_MCP_API_KEY` — upstream API key
- `TARS_MCP_MODE` — `stdio` (default) or `http`

## License

MIT
