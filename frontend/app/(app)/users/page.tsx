import Link from "next/link";

export default function UsersPage() {
  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
          Users
        </p>
        <h1 className="mt-2 font-[var(--font-display)] text-3xl">
          ユーザー検索
        </h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          キーワードでユーザーを検索し、プロフィールを閲覧できます。
        </p>
        <div className="mt-4 flex flex-wrap gap-3">
          <input className="flex-1 rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]" placeholder="ユーザー名で検索" />
          <button className="rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white">
            検索
          </button>
        </div>
      </section>

      <section className="grid gap-4 md:grid-cols-2">
        {[
          { id: "u1", name: "mizu", note: "本棚: 82冊" },
          { id: "u2", name: "kent", note: "お気に入り: 12冊" },
        ].map((user) => (
          <Link
            key={user.id}
            className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 shadow-sm transition hover:-translate-y-1 hover:shadow-md"
            href={`/users/${user.id}`}
          >
            <p className="text-xs text-[#c86b3c]">@{user.name}</p>
            <p className="mt-2 font-[var(--font-display)] text-xl">
              {user.name}
            </p>
            <p className="mt-2 text-sm text-[#5c5d63]">{user.note}</p>
          </Link>
        ))}
      </section>
    </div>
  );
}
