export default function RecommendationsPage() {
  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
          Recommendations
        </p>
        <h1 className="mt-2 font-[var(--font-display)] text-3xl">
          おすすめ一覧
        </h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          みんなのおすすめがタイムライン形式で表示されます。
        </p>
      </section>

      <section className="grid gap-4">
        {[
          {
            user: "mizu",
            title: "夜更けの灯台",
            comment: "静かな夜に読みたくなる一冊。",
          },
          {
            user: "kent",
            title: "星街メトロ 4",
            comment: "シリーズ最高潮。世界観の広がりが良い。",
          },
        ].map((item) => (
          <div
            key={item.title}
            className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm"
          >
            <p className="text-xs text-[#c86b3c]">@{item.user}</p>
            <p className="mt-2 font-[var(--font-display)] text-xl">
              {item.title}
            </p>
            <p className="mt-2 text-sm text-[#5c5d63]">{item.comment}</p>
          </div>
        ))}
      </section>

      <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
        <h2 className="font-[var(--font-display)] text-2xl">おすすめ投稿</h2>
        <div className="mt-4 flex flex-col gap-3 text-sm">
          <input className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]" placeholder="書籍名 or ISBN" />
          <textarea
            className="min-h-[120px] rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
            placeholder="コメント"
          />
          <button className="rounded-full bg-[#c86b3c] px-5 py-3 text-sm font-medium text-white">
            投稿する
          </button>
        </div>
      </section>
    </div>
  );
}
