# discord-hands-up-for-sq

## 使い方

1. https://discord.com/api/oauth2/authorize?client_id=1168004102282813512&scope=bot から招待

## 開発

### 単体テスト

```bash
gcloud emulators firestore start --host-port=localhost:5000
# 別ターミナル
go test ./...
```

### テスト実行

```bash
gcloud emulators firestore start --host-port=localhost:5000
# 別ターミナル
## 環境変数の設定 (windowsの場合)
$env:BOT_TOKEN="[YOUR BOT TOKEN]"
$env:FIREBASE_PROJECT_ID="test"
$env:FIRESTORE_EMULATOR_HOST="localhost:5000"
go run main.go -guild GUILD_ID
```

### デプロイ

#### Fly.io

fly.tomlの用意

```
flyctl deploy -a YOUR_APP_NAME
```
