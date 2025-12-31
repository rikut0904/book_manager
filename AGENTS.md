# Repository Guidelines

## プロジェクト構成とモジュール
- `frontend/`: Next.js App Router（TypeScript + Tailwind）。`app/` が画面、`public/` が静的資産。設定は `next.config.ts` と `eslint.config.mjs`。
- `backend/`: Go API（クリーンアーキテクチャ）用の領域。実装開始後に `cmd/` や `internal/` を配置。
- `doc/`: 要件・設計・運用ドキュメント。

## ビルド/テスト/開発コマンド
`frontend/` で実行:
- `npm run dev`: 開発サーバ起動
- `npm run build`: 本番ビルド
- `npm run start`: 本番サーバ起動
- `npm run lint`: ESLint 実行

バックエンドのコマンドは未定義。Go 側の雛形作成後に追記。

## コーディング規約・命名
- フロントは TypeScript + ESLint。Next.js App Router の慣習に従う。
- インデントは 2 スペース。
- ルートやページ名は機能名で明確に（例: `app/books/page.tsx`）。
- バックエンドは `gofmt` 準拠。クリーンアーキテクチャの層分離を厳守。

## テスト指針
- まだテスト基盤は未設定。導入時に実行方法と命名（例: `*.test.tsx`）を明記。

## コミット・PR 指針
- 既存履歴がないため、明確で短い形式を推奨（例: `feat: add isbn lookup`）。
- PR には要約、テスト方法、UI変更時のスクショを含める。

## セキュリティ・設定
- 秘密情報は `frontend/.env.local` / `backend/.env` で管理し、追跡しない。
- 監査ログは操作内容のみ記録し、秘匿情報は残さない。
