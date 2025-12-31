# アーキテクチャ設計

## 構成
- フロント: Next.js (Vercel)
- バック: Go API (Railway)
- DB: PostgreSQL (Railway)

## 認証
- Email/Password
- JWT + Refresh Token

## 外部連携
- Google Books API: ISBN→書誌情報

## ログ
- DB操作ログを audit_logs に保存
- 90日保持（定期ジョブで削除）

## メール
- 書誌情報の誤り報告は product@rikut0904.site へ送信
