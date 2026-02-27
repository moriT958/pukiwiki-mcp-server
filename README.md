# PukiWiki MCP Server

## How to setting

edit `~/Library/Application\ Support/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "pukiwiki": {
      "command": "/path/to/pukiwiki-mcp-binary",
      "env": {
        "PUKIWIKI_URL": "https://your-wiki.example.jp",
        "PUKIWIKI_USER": "user",
        "PUKIWIKI_PASS": "pass"
      }
    }
  }
}
```
