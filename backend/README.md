# Backend

## 起動方法
```bash
go run ./cmd/api
```

## 実装状況
- Auth: インメモリ実装（再起動で消える）
- ISBN lookup: Google Books API 連携
- 書誌マスタ（/books）: インメモリ実装
- 所蔵（/user-books）: インメモリ実装
- Users/Follows: インメモリ実装
- そのほか: echo 返却

## 環境変数
- `PORT`: APIのポート（default: 8080）
- `APP_ENV`: 実行環境名（default: local）
- `GOOGLE_BOOKS_API_KEY`: Google Books APIキー（未設定でも動作）
- `GOOGLE_BOOKS_BASE_URL`: APIベースURL（default: https://www.googleapis.com/books/v1/volumes）

## ヘルスチェック
```bash
curl http://localhost:8080/healthz
```
