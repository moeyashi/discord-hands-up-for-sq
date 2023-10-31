# discord-hands-up-for-sq

## 使い方

1. https://discord.com/api/oauth2/authorize?client_id=1168004102282813512&scope=bot から招待

## 開発

### 準備

```bash
npm install -g firebase-tools
firebase login --no-localhost
```

### 単体テスト

```bash
go test ./...
```

### テスト実行

```bash
firebase emulators:start
# 別ターミナル
go run main.go -guild GUILD_ID -token DISCORD_BOT_TOKEN
```

### デプロイ

#### Fly.io

fly.tomlの用意

```
app = "YOUR APP NAME"
kill_signal = "SIGINT"
kill_timeout = 5
processes = []

[env]
```