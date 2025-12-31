import Link from "next/link";

const navItems = [
  { href: "/books", label: "所蔵一覧" },
  { href: "/books/new", label: "書籍登録" },
  { href: "/favorites", label: "お気に入り" },
  { href: "/next-to-buy", label: "次に買う本" },
  { href: "/recommendations", label: "おすすめ" },
  { href: "/users", label: "ユーザー検索" },
  { href: "/settings", label: "設定" },
];

export default function AppLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen">
      <header className="sticky top-0 z-10 border-b border-[#e4d8c7] bg-[#f6f1e7]/90 backdrop-blur">
        <div className="mx-auto flex w-full max-w-6xl items-center justify-between px-6 py-4">
          <div className="flex items-center gap-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-full border border-[#e4d8c7] bg-white text-sm font-semibold">
              B
            </div>
            <div>
              <p className="text-[10px] tracking-[0.25em] text-[#5c5d63]">
                BOOK
              </p>
              <p className="font-[var(--font-display)] text-lg font-semibold">
                Manager
              </p>
            </div>
          </div>
          <div className="flex items-center gap-4">
            <div className="hidden items-center gap-2 rounded-full border border-[#e4d8c7] bg-white/80 px-4 py-2 text-xs text-[#5c5d63] md:flex">
              <span className="h-2 w-2 rounded-full bg-[#c86b3c]" />
              同期済み
            </div>
            <div className="flex items-center gap-2 text-sm">
              <span className="text-[#5c5d63]">rikut</span>
              <div className="h-9 w-9 rounded-full bg-[#1b1c1f] text-center text-sm leading-9 text-white">
                R
              </div>
            </div>
          </div>
        </div>
      </header>
      <div className="mx-auto grid w-full max-w-6xl gap-6 px-6 py-8 lg:grid-cols-[220px_1fr]">
        <aside className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-4 shadow-sm">
          <nav className="flex flex-col gap-1 text-sm">
            {navItems.map((item) => (
              <Link
                key={item.href}
                className="rounded-2xl px-3 py-2 text-[#5c5d63] transition hover:bg-[#efe5d4] hover:text-[#1b1c1f]"
                href={item.href}
              >
                {item.label}
              </Link>
            ))}
          </nav>
          <div className="mt-6 rounded-2xl border border-[#e4d8c7] bg-[#f6f1e7] p-3 text-xs text-[#5c5d63]">
            <p className="font-medium text-[#1b1c1f]">書誌情報の修正</p>
            <p className="mt-2">
              誤り報告はフォームから送信。運営が確認後に反映します。
            </p>
          </div>
        </aside>
        <main className="flex min-h-[70vh] flex-col gap-6">
          {children}
        </main>
      </div>
    </div>
  );
}
