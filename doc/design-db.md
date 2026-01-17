# データ設計（PostgreSQL）

## 方針
- 書誌/シリーズはマスタ保持（ユーザーから更新不可）
- 所蔵は user_books を中心に管理
- シリーズ判定は自動 + ユーザー上書き
- ログは 90 日保持

## テーブル

### users
- id (uuid, PK)
- email (unique)
- password_hash
- user_id (unique)
- status (active/deleted)
- deleted_at
- created_at, updated_at

### profile_settings
- user_id (PK, FK users)
- visibility (public/followers)

### books
- id (PK)
- isbn13 (unique, nullable)
- title
- authors (text[])
- publisher
- published_date
- thumbnail_url
- source (manual/google)
- source_payload (jsonb)

### series
- id (PK)
- name
- normalized_name (unique)

### book_series_auto
- book_id (FK books)
- series_id (FK series)
- volume_number (int)
- confidence (int)
- unique(book_id)

### user_books
- id (PK)
- user_id (FK users)
- book_id (FK books)
- note (text)
- acquired_at (date)
- unique(user_id, book_id)

### user_book_series_override
- user_id (FK users)
- book_id (FK books)
- series_id (FK series)
- volume_number (int)
- unique(user_id, book_id)

### favorites
- id (PK)
- user_id
- book_id (nullable)
- series_id (nullable)
- type (book/series)
- check: book_id xor series_id
- unique(user_id, book_id)
- unique(user_id, series_id)

### next_to_buy_manual
- id (PK)
- user_id
- title
- series_name (nullable)
- volume_number (nullable)
- note

### recommendations
- id (PK)
- user_id (FK users)
- book_id (FK books)
- comment
- created_at
- on user delete: cascade

### follows
- follower_id
- followee_id
- created_at
- unique(follower_id, followee_id)

### isbn_cache (任意)
- isbn13 (PK)
- payload (jsonb)
- fetched_at

### audit_logs
- id (PK)
- user_id
- action (create/update/delete/read)
- entity
- entity_id
- payload (jsonb)
- ip
- user_agent
- created_at
- retention: 90 days

## 主要インデックス
- user_books(user_id)
- recommendations(created_at)
- follows(followee_id)
- book_series_auto(series_id)

## 集計ルール
- series_count は自動判定 + ユーザー上書きの両方を含めて算出
  - series_id = coalesce(user_book_series_override.series_id, book_series_auto.series_id)
