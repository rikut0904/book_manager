import Link from "next/link";

export default function Home() {
  return (
    <div className="min-h-screen">
      <header className="mx-auto flex w-full max-w-6xl items-center justify-between px-6 pt-8">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-full border border-[#e4d8c7] bg-white/70 text-lg font-semibold">
            B
          </div>
          <div>
            <p className="text-sm tracking-[0.2em] text-[#5c5d63]">BOOK</p>
            <p className="font-[var(--font-display)] text-xl font-semibold">
              Manager
            </p>
          </div>
        </div>
        <nav className="hidden items-center gap-6 text-sm text-[#5c5d63] md:flex">
          <Link className="hover:text-[#1b1c1f]" href="/books">
            所蔵一覧
          </Link>
          <Link className="hover:text-[#1b1c1f]" href="/books/new">
            登録
          </Link>
          <Link className="hover:text-[#1b1c1f]" href="/user/suggest">
            おすすめの投稿
          </Link>
        </nav>
        <div className="flex items-center gap-3">
          <Link
            className="text-sm text-[#5c5d63] hover:text-[#1b1c1f]"
            href="/login"
          >
            ログイン
          </Link>
          <Link
            className="rounded-full bg-[#c86b3c] px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-[#8f3d1f]"
            href="/signup"
          >
            新規登録
          </Link>
        </div>
      </header>

      <main className="mx-auto grid w-full max-w-6xl gap-10 px-6 pb-20 pt-16 lg:grid-cols-[1.1fr_0.9fr]">
        <section className="flex flex-col gap-6">
          <p className="text-xs uppercase tracking-[0.4em] text-[#c86b3c]">
            Personal Library Toolkit
          </p>
          <h1 className="font-[var(--font-display)] text-4xl leading-tight text-[#1b1c1f] md:text-5xl">
            ISBNから一瞬で登録。
            <br />
            迷わない本棚へ。
          </h1>
          <p className="max-w-xl text-base leading-7 text-[#5c5d63]">
            所蔵、シリーズ、次に買う本まで。読む・買う・残すを一本の流れで
            管理できる、あなた専用の本棚アプリ。
          </p>
          <div className="flex flex-wrap gap-3">
            <Link
              className="rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white shadow-lg shadow-[#1b1c1f]/10 transition hover:bg-black"
              href="/books"
            >
              さっそく使う
            </Link>
            <Link
              className="rounded-full border border-[#1b1c1f] px-5 py-3 text-sm font-medium text-[#1b1c1f] transition hover:bg-[#1b1c1f] hover:text-white"
              href="/books/new"
            >
              書籍を登録
            </Link>
          </div>
          <div className="grid gap-4 md:grid-cols-2">
            <div className="rounded-2xl border border-[#e4d8c7] bg-white/70 p-5 shadow-sm">
              <p className="text-xs text-[#c86b3c]">ISBN検索</p>
              <p className="mt-2 text-sm text-[#5c5d63]">
                ISBN入力で書誌を取得。失敗時は手動登録に切り替え。
              </p>
            </div>
            <div className="rounded-2xl border border-[#e4d8c7] bg-white/70 p-5 shadow-sm">
              <p className="text-xs text-[#c86b3c]">シリーズ管理</p>
              <p className="mt-2 text-sm text-[#5c5d63]">
                自動判定のシリーズ情報を確認し、必要なら上書き。
              </p>
            </div>
            <div className="rounded-2xl border border-[#e4d8c7] bg-white/70 p-5 shadow-sm">
              <p className="text-xs text-[#c86b3c]">ブックマーク</p>
              <p className="mt-2 text-sm text-[#5c5d63]">
                単巻・シリーズを登録して、次に買う本へ繋げる。
              </p>
            </div>
            <div className="rounded-2xl border border-[#e4d8c7] bg-white/70 p-5 shadow-sm">
              <p className="text-xs text-[#c86b3c]">次に買う本</p>
              <p className="mt-2 text-sm text-[#5c5d63]">
                手入力のメモもまとめて管理。買い忘れ防止。
              </p>
            </div>
            <div className="rounded-2xl border border-[#e4d8c7] bg-white/70 p-5 shadow-sm md:col-span-2">
              <p className="text-xs text-[#c86b3c]">AIシリーズ推定</p>
              <p className="mt-2 text-sm text-[#5c5d63]">
                シリーズ推定にOpenAI APIを利用します。設定画面から有効化できます。
              </p>
              <Link
                className="mt-3 inline-flex text-xs text-[#c86b3c] hover:text-[#8f3d1f]"
                href="/settings/ai"
              >
                AI設定を開く
              </Link>
            </div>
          </div>
        </section>

        <section className="flex flex-col gap-4">
          <div className="rounded-[32px] border border-[#e4d8c7] bg-[#1b1c1f] p-6 text-white shadow-xl">
            <div className="flex items-center justify-between text-sm text-white/70">
              <span>今日の本棚</span>
              <span>2024.11.07</span>
            </div>
            <div className="mt-6 space-y-4">
              <div className="rounded-2xl bg-white/10 p-4">
                <p className="text-sm text-white/60">所蔵数</p>
                <p className="mt-1 text-3xl font-semibold">128冊</p>
                <p className="mt-2 text-xs text-white/60">
                  シリーズ: 18 / ブックマーク: 7
                </p>
              </div>
              <div className="rounded-2xl bg-white/10 p-4">
                <p className="text-sm text-white/60">次に買う本</p>
                <p className="mt-1 text-lg font-medium">星街メトロ 5</p>
                <p className="mt-1 text-xs text-white/60">
                  メモ: 発売日をチェック
                </p>
              </div>
              <div className="rounded-2xl bg-white/10 p-4">
                <p className="text-sm text-white/60">直近の登録</p>
                <p className="mt-1 text-lg font-medium">
                  透明な約束 - 冬原 千景
                </p>
                <p className="mt-1 text-xs text-white/60">
                  ISBN: 978-4-00-123456-7
                </p>
              </div>
            </div>
          </div>
          <div className="rounded-[28px] border border-[#e4d8c7] bg-white/70 p-6 shadow-md">
            <p className="text-xs uppercase tracking-[0.3em] text-[#5c5d63]">
              MVP ROADMAP
            </p>
            <ul className="mt-4 space-y-3 text-sm text-[#5c5d63]">
              <li className="flex items-center justify-between">
                認証と所蔵管理
                <span className="rounded-full bg-[#efe5d4] px-3 py-1 text-xs text-[#1b1c1f]">
                  Now
                </span>
              </li>
              <li className="flex items-center justify-between">
                ブックマークと次に買う本
                <span className="text-xs text-[#c86b3c]">Next</span>
              </li>
              <li className="flex items-center justify-between">
                おすすめ・プロフィール
                <span className="text-xs text-[#5c5d63]">Later</span>
              </li>
            </ul>
          </div>
        </section>
      </main>
    </div>
  );
}
