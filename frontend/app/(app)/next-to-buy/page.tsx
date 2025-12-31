export default function NextToBuyPage() {
  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
          Next To Buy
        </p>
        <h1 className="mt-2 font-[var(--font-display)] text-3xl">
          次に買う本
        </h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          お気に入りの次巻と手入力メモをまとめます。
        </p>
      </section>

      <section className="grid gap-4 lg:grid-cols-3">
        {[
          { title: "星街メトロ 5", note: "発売日を確認" },
          { title: "透明な約束 番外編", note: "電子で購入" },
          { title: "空間の縁 2", note: "古書店で探す" },
        ].map((item) => (
          <div
            key={item.title}
            className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 shadow-sm"
          >
            <p className="font-[var(--font-display)] text-xl">{item.title}</p>
            <p className="mt-2 text-sm text-[#5c5d63]">{item.note}</p>
          </div>
        ))}
      </section>

      <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
        <h2 className="font-[var(--font-display)] text-2xl">手動追加</h2>
        <div className="mt-4 grid gap-3 text-sm md:grid-cols-[1fr_1fr]">
          <input className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]" placeholder="タイトル" />
          <input className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]" placeholder="シリーズ名 (任意)" />
          <input className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]" placeholder="巻数" />
          <input className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]" placeholder="メモ" />
        </div>
        <button className="mt-4 rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white">
          追加
        </button>
      </section>
    </div>
  );
}
