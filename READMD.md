# Book Manager 設計まとめ

本ドキュメントは、Web向け書籍管理アプリの全体像を簡潔にまとめたものです。
詳細は各設計/マニュアルドキュメントを参照してください。

## 目的
- 書籍を手動/ISBN/バーコードで登録し、所蔵・シリーズ・次に買う本を管理する
- お気に入りやおすすめ、プロフィール公開を通じて共有性も持たせる

## 主要機能（確定）
- 書籍登録: 手動入力/ISBN取得/バーコード
- 所蔵検索
- シリーズ自動判定 + 手動上書き
- お気に入り: 単巻/シリーズ
- 次に買う本: シリーズ次巻 + 手入力
- おすすめ投稿: 全体公開
- プロフィール: 公開範囲の切替（公開/フォロワー限定）
- ユーザー検索
- 誤った書誌情報の報告（メール送信）

## 技術スタック
- フロント: Next.js
- バック: Go
- DB: PostgreSQL (Railway)
- ホスティング: Vercel (FE), Railway (BE/DB)
- ISBN: Google Books API

## 設計ドキュメント
- 画面設計: doc/design-ui-flow.md
- データ設計: doc/design-db.md
- ER関係: doc/design-er.md
- API設計: doc/design-api.md
- APIサンプル: doc/design-api-samples.md
- アーキ設計: doc/design-architecture.md
- 実装順: doc/plan-mvp.md

## マニュアル
- 利用マニュアル: doc/manual-usage.md
- 運用/監視マニュアル: doc/manual-operations.md
