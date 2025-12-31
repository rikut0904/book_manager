export default function FavoritesPage() {
  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
          Favorites
        </p>
        <h1 className="mt-2 font-[var(--font-display)] text-3xl">
          お気に入り
        </h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          単巻とシリーズのお気に入りをまとめて管理します。
        </p>
      </section>
      <section className="grid gap-4 md:grid-cols-2">
        {[
          { name: "星街メトロ", type: "シリーズ", note: "次巻通知オン" },
          { name: "海辺の標本室", type: "単巻", note: "再読予定" },
        ].map((item) => (
          <div
            key={item.name}
            className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 shadow-sm"
          >
            <p className="text-xs text-[#c86b3c]">{item.type}</p>
            <p className="mt-2 font-[var(--font-display)] text-xl">
              {item.name}
            </p>
            <p className="mt-2 text-sm text-[#5c5d63]">{item.note}</p>
          </div>
        ))}
      </section>
    </div>
  );
}
