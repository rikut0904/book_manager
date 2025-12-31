# Backend

## 起動方法
```bash
go run ./cmd/api
```

## 実装状況
- Auth: インメモリ実装（再起動で消える）
- ISBN lookup: 仮データのみ（`9780000000000`）
- そのほか: echo 返却

## 環境変数
- `PORT`: APIのポート（default: 8080）
- `APP_ENV`: 実行環境名（default: local）

## ヘルスチェック
```bash
curl http://localhost:8080/healthz
```
