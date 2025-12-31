export default function BookNewPage() {
  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
          Register
        </p>
        <h1 className="mt-2 font-[var(--font-display)] text-3xl">
          書籍登録
        </h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          ISBN検索で書誌を自動取得。見つからない場合は手動で登録できます。
        </p>
      </section>

      <div className="grid gap-6 lg:grid-cols-[1.2fr_0.8fr]">
        <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
          <h2 className="font-[var(--font-display)] text-2xl">ISBN検索</h2>
          <p className="mt-2 text-sm text-[#5c5d63]">
            13桁のISBNを入力してください。バーコード入力は後から追加予定です。
          </p>
          <div className="mt-4 flex flex-wrap gap-3">
            <input
              className="flex-1 rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              placeholder="978-4-00-123456-7"
            />
            <button className="rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white">
              取得
            </button>
          </div>
          <div className="mt-6 rounded-2xl border border-dashed border-[#e4d8c7] bg-white/70 p-4 text-sm text-[#5c5d63]">
            <p className="font-medium text-[#1b1c1f]">取得結果</p>
            <p className="mt-2">タイトル: 透明な約束</p>
            <p>著者: 冬原 千景</p>
            <p>出版社: 北風文庫 / 2024-01-01</p>
          </div>
        </section>

        <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
          <h2 className="font-[var(--font-display)] text-2xl">手動登録</h2>
          <div className="mt-4 flex flex-col gap-3 text-sm">
            <label className="text-[#1b1c1f]">
              タイトル
              <input className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]" />
            </label>
            <label className="text-[#1b1c1f]">
              著者
              <input className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]" />
            </label>
            <label className="text-[#1b1c1f]">
              ISBN (任意)
              <input className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]" />
            </label>
            <button className="rounded-full bg-[#c86b3c] px-5 py-3 text-sm font-medium text-white">
              所蔵に追加
            </button>
          </div>
        </section>
      </div>
    </div>
  );
}
