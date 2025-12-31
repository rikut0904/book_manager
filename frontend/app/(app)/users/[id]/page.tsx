export default function UserDetailPage() {
  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
          Profile
        </p>
        <h1 className="mt-2 font-[var(--font-display)] text-3xl">mizu</h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          公開範囲: フォロワー限定
        </p>
      </section>

      <section className="grid gap-4 md:grid-cols-3">
        {[
          { label: "所蔵数", value: "82冊" },
          { label: "シリーズ", value: "14件" },
          { label: "フォロワー", value: "5人" },
        ].map((stat) => (
          <div
            key={stat.label}
            className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 shadow-sm"
          >
            <p className="text-xs text-[#5c5d63]">{stat.label}</p>
            <p className="mt-2 font-[var(--font-display)] text-2xl">
              {stat.value}
            </p>
          </div>
        ))}
      </section>

      <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
        <h2 className="font-[var(--font-display)] text-2xl">
          おすすめした本
        </h2>
        <div className="mt-4 grid gap-3 md:grid-cols-2">
          {["夜更けの灯台", "空間の縁 1"].map((title) => (
            <div
              key={title}
              className="rounded-2xl border border-[#e4d8c7] bg-white p-4 text-sm text-[#5c5d63]"
            >
              <p className="font-medium text-[#1b1c1f]">{title}</p>
              <p className="mt-1 text-xs text-[#5c5d63]">
                コメント: 読後感がとても良い
              </p>
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
