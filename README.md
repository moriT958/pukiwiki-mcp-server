# PukiWiki MCP Server

- Support 
  - MacOS (認証情報の保存に MacOS Keychain を使用します)
  - リモート MCP は未対応なため、ChatGPT ではまだ利用出ません。

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

## TODO

- MacOS 以外もサポート
- リモート MCP
