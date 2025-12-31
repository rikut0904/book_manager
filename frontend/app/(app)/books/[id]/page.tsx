import Link from "next/link";

export default function BookDetailPage() {
  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
              Book Detail
            </p>
            <h1 className="mt-2 font-[var(--font-display)] text-3xl">
              透明な約束
            </h1>
            <p className="mt-1 text-sm text-[#5c5d63]">冬原 千景</p>
          </div>
          <Link
            className="rounded-full border border-[#e4d8c7] px-4 py-2 text-xs text-[#5c5d63]"
            href="/books"
          >
            一覧へ戻る
          </Link>
        </div>
        <div className="mt-6 grid gap-4 md:grid-cols-3">
          <div className="rounded-2xl border border-[#e4d8c7] bg-white/70 p-4 text-sm text-[#5c5d63]">
            <p className="text-xs text-[#c86b3c]">ISBN</p>
            <p className="mt-2 text-[#1b1c1f]">978-4-00-123456-7</p>
          </div>
          <div className="rounded-2xl border border-[#e4d8c7] bg-white/70 p-4 text-sm text-[#5c5d63]">
            <p className="text-xs text-[#c86b3c]">出版社</p>
            <p className="mt-2 text-[#1b1c1f]">北風文庫</p>
          </div>
          <div className="rounded-2xl border border-[#e4d8c7] bg-white/70 p-4 text-sm text-[#5c5d63]">
            <p className="text-xs text-[#c86b3c]">所蔵メモ</p>
            <p className="mt-2 text-[#1b1c1f]">サイン本</p>
          </div>
        </div>
      </section>

      <section className="grid gap-6 lg:grid-cols-[1.2fr_0.8fr]">
        <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
          <h2 className="font-[var(--font-display)] text-2xl">
            シリーズ判定
          </h2>
          <p className="mt-2 text-sm text-[#5c5d63]">
            自動判定されたシリーズ情報を確認し、必要なら上書きしてください。
          </p>
          <div className="mt-4 rounded-2xl border border-[#e4d8c7] bg-white p-4 text-sm text-[#5c5d63]">
            <div className="flex items-center justify-between">
              <span className="text-[#1b1c1f]">星街メトロ</span>
              <span>Vol.4</span>
            </div>
            <p className="mt-2 text-xs text-[#c86b3c]">自動判定: 信頼度 86%</p>
          </div>
          <div className="mt-4 grid gap-3 text-sm">
            <label className="text-[#1b1c1f]">
              シリーズ名
              <input className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]" />
            </label>
            <label className="text-[#1b1c1f]">
              巻数
              <input className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]" />
            </label>
            <button className="rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white">
              上書き保存
            </button>
          </div>
        </div>

        <div className="flex flex-col gap-6">
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
            <h2 className="font-[var(--font-display)] text-2xl">タグ</h2>
            <div className="mt-4 flex flex-wrap gap-2 text-xs">
              {["推し", "再読", "貸出中"].map((tag) => (
                <span
                  key={tag}
                  className="rounded-full border border-[#e4d8c7] bg-white px-3 py-2 text-[#5c5d63]"
                >
                  {tag}
                </span>
              ))}
              <button className="rounded-full bg-[#efe5d4] px-3 py-2 text-[#1b1c1f]">
                + タグ追加
              </button>
            </div>
          </div>

          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
            <h2 className="font-[var(--font-display)] text-2xl">お気に入り</h2>
            <p className="mt-2 text-sm text-[#5c5d63]">
              この巻またはシリーズをお気に入りに登録できます。
            </p>
            <div className="mt-4 flex gap-3">
              <button className="rounded-full border border-[#e4d8c7] px-4 py-2 text-xs text-[#5c5d63]">
                単巻で登録
              </button>
              <button className="rounded-full bg-[#c86b3c] px-4 py-2 text-xs text-white">
                シリーズで登録
              </button>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}
