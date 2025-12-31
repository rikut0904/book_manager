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
- Favorites/Next-to-buy/Tags/Recommendations: インメモリ実装
- そのほか: echo 返却

## ISBN lookup の挙動
- `/isbn/lookup` は取得した書誌を `/books` に自動登録します
- ISBN 取得結果はキャッシュします（`ISBN_CACHE_TTL_MINUTES`）

## シリーズ上書き
- `/user-series/override` は既存の user-book を探して seriesId と volumeNumber を更新します

## 書誌報告
- `/book-reports` はメール送信の代わりにログ出力します

## 環境変数
- `PORT`: APIのポート（default: 8080）
- `APP_ENV`: 実行環境名（default: local）
- `GOOGLE_BOOKS_API_KEY`: Google Books APIキー（未設定でも動作）
- `GOOGLE_BOOKS_BASE_URL`: APIベースURL（default: https://www.googleapis.com/books/v1/volumes）
- `DATABASE_URL`: PostgreSQL 接続URL（未設定時はメモリ実装）

## ヘルスチェック
```bash
curl http://localhost:8080/healthz
```
