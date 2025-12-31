export default function SettingsPage() {
  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
          Settings
        </p>
        <h1 className="mt-2 font-[var(--font-display)] text-3xl">設定</h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          タグ管理とプロフィール公開範囲を設定します。
        </p>
      </section>

      <section className="grid gap-6 lg:grid-cols-[1fr_1fr]">
        <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
          <h2 className="font-[var(--font-display)] text-2xl">タグ管理</h2>
          <div className="mt-4 flex flex-wrap gap-2 text-xs">
            {["積読", "貸出中", "推し"].map((tag) => (
              <span
                key={tag}
                className="rounded-full border border-[#e4d8c7] bg-white px-3 py-2 text-[#5c5d63]"
              >
                {tag}
              </span>
            ))}
            <button className="rounded-full bg-[#efe5d4] px-3 py-2 text-[#1b1c1f]">
              + 新規タグ
            </button>
          </div>
        </div>

        <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
          <h2 className="font-[var(--font-display)] text-2xl">
            プロフィール公開
          </h2>
          <div className="mt-4 flex flex-col gap-3 text-sm">
            <label className="flex items-center justify-between rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3">
              <span>公開</span>
              <input type="radio" name="visibility" defaultChecked />
            </label>
            <label className="flex items-center justify-between rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3">
              <span>フォロワー限定</span>
              <input type="radio" name="visibility" />
            </label>
          </div>
          <button className="mt-4 rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white">
            変更を保存
          </button>
        </div>
      </section>
    </div>
  );
}
