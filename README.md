# PukiWiki MCP Server

## Build

```bash
task build
```

## How to setting

edit `~/Library/Application\ Support/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "pukiwiki": {
      "command": "/path/to/your/pukiwiki-mcp-binary"
    }
  }
}
```

## Go Pukiwiki package

```bash
go get github.com/moriT958/pukiwiki-mcp/pukiwiki
```

### Usage

[examples](/examples)
