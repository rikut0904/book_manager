import Link from "next/link";

export default function SignupPage() {
  return (
    <div className="mx-auto flex min-h-screen w-full max-w-4xl items-center justify-center px-6 py-16">
      <div className="grid w-full gap-8 rounded-[32px] border border-[#e4d8c7] bg-white/80 p-8 shadow-lg md:grid-cols-[1.1fr_0.9fr]">
        <div className="flex flex-col justify-between gap-6">
          <div>
            <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
              Create account
            </p>
            <h1 className="mt-3 font-[var(--font-display)] text-3xl">
              新しい本棚をはじめる
            </h1>
            <p className="mt-3 text-sm leading-6 text-[#5c5d63]">
              ISBN検索、シリーズ自動判定、次に買う本まで。最初の登録から一緒に準備しましょう。
            </p>
          </div>
          <div className="rounded-2xl border border-[#e4d8c7] bg-[#f6f1e7] p-4 text-xs text-[#5c5d63]">
            <p className="font-medium text-[#1b1c1f]">無料ではじめる</p>
            <p className="mt-2">メールアドレスと表示名だけで開始できます。</p>
          </div>
        </div>

        <form className="flex flex-col gap-4">
          <label className="text-sm text-[#1b1c1f]">
            表示名
            <input
              className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              name="username"
              placeholder="rikut"
              type="text"
            />
          </label>
          <label className="text-sm text-[#1b1c1f]">
            メールアドレス
            <input
              className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              name="email"
              placeholder="user@example.com"
              type="email"
            />
          </label>
          <label className="text-sm text-[#1b1c1f]">
            パスワード
            <input
              className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              name="password"
              placeholder="********"
              type="password"
            />
          </label>
          <button
            className="mt-2 rounded-full bg-[#c86b3c] px-5 py-3 text-sm font-medium text-white shadow-lg shadow-[#c86b3c]/20 transition hover:bg-[#8f3d1f]"
            type="button"
          >
            アカウント作成
          </button>
          <div className="flex items-center justify-between text-xs text-[#5c5d63]">
            <span>登録すると利用規約に同意したことになります。</span>
            <Link className="hover:text-[#1b1c1f]" href="/login">
              ログインへ
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
}
