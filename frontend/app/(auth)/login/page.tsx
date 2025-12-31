import Link from "next/link";

export default function LoginPage() {
  return (
    <div className="mx-auto flex min-h-screen w-full max-w-4xl items-center justify-center px-6 py-16">
      <div className="grid w-full gap-8 rounded-[32px] border border-[#e4d8c7] bg-white/80 p-8 shadow-lg md:grid-cols-[1.1fr_0.9fr]">
        <div className="flex flex-col justify-between gap-6">
          <div>
            <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
              Sign in
            </p>
            <h1 className="mt-3 font-[var(--font-display)] text-3xl">
              本棚に戻る
            </h1>
            <p className="mt-3 text-sm leading-6 text-[#5c5d63]">
              ログインして所蔵の更新、シリーズの確認、次に買う本の管理を続けましょう。
            </p>
          </div>
          <div className="rounded-2xl border border-[#e4d8c7] bg-[#f6f1e7] p-4 text-xs text-[#5c5d63]">
            <p className="font-medium text-[#1b1c1f]">テスト用アカウント</p>
            <p className="mt-2">email: demo@book.local / password: password</p>
          </div>
        </div>

        <form className="flex flex-col gap-4">
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
            className="mt-2 rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white shadow-lg shadow-[#1b1c1f]/10 transition hover:bg-black"
            type="button"
          >
            ログイン
          </button>
          <div className="flex items-center justify-between text-xs text-[#5c5d63]">
            <Link className="hover:text-[#1b1c1f]" href="#">
              パスワードを忘れた
            </Link>
            <Link className="hover:text-[#1b1c1f]" href="/signup">
              新規登録へ
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
}
