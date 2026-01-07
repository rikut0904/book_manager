# API設計（主要エンドポイント）

## 認証
- POST /auth/signup
- POST /auth/login
- POST /auth/refresh
- POST /auth/logout

## ISBN
- GET /isbn/lookup?isbn=978...
  - 書誌マスタ未登録なら取得 → books に登録

## 書誌（マスタ）
- POST /books
  - ISBNがあれば Google Books 取得 → 取得不可なら手入力
- GET /books/{id}
  - PATCH/DELETE は不可（マスタは更新不可）

## 所蔵（ユーザー中心）
- POST /user-books
- GET /user-books?query=&series=&page=
- PATCH /user-books/{id}
- DELETE /user-books/{id}

## シリーズ上書き
- PATCH /user-series/override
  - req: {bookId, seriesId, volumeNumber}

## お気に入り
- POST /favorites
- GET /favorites
- DELETE /favorites/{id}

## 次に買う本
- GET /next-to-buy
- POST /next-to-buy/manual
- PATCH /next-to-buy/manual/{id}
- DELETE /next-to-buy/manual/{id}

## おすすめ（全体公開）
- GET /recommendations
- POST /recommendations
- DELETE /recommendations/{id}

## ユーザー/プロフィール
- GET /users?query=
- GET /users/{id}
- PATCH /users/me
- PATCH /users/me/settings
- DELETE /users/me
  - recommendations も削除

## フォロー
- POST /follows/{userId}
- DELETE /follows/{userId}

## 書誌報告
- POST /book-reports
  - req: {bookId, suggestion, note?}
  - メール送信先: product@rikut0904.site
  - 本文: ISBN + 現在の書誌情報 + 修正提案 + 備考
