"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { fetchJSON } from "@/lib/api";

type UserItem = {
  id: string;
  email: string;
  userId: string;
};

export default function UsersPage() {
  const [items, setItems] = useState<UserItem[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let isMounted = true;
    fetchJSON<{ items: UserItem[] }>("/users", { auth: true })
      .then((data) => {
        if (!isMounted) {
          return;
        }
        setItems(data.items ?? []);
      })
      .catch(() => {
        if (!isMounted) {
          return;
        }
        setError("ユーザー一覧を取得できませんでした。");
      });
    return () => {
      isMounted = false;
    };
  }, []);

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
          <input className="flex-1 rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]" placeholder="ユーザーIDで検索" />
          <button className="rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white">
            検索
          </button>
        </div>
      </section>

      <section className="grid gap-4 md:grid-cols-2">
        {error ? (
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 text-sm text-red-600">
            {error}
          </div>
        ) : null}
        {!error && items.length === 0 ? (
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 text-sm text-[#5c5d63]">
            まだユーザーが登録されていません。
          </div>
        ) : null}
        {items.map((user) => (
          <Link
            key={user.id}
            className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 shadow-sm transition hover:-translate-y-1 hover:shadow-md"
            href={`/users/${user.id}`}
          >
            <p className="text-xs text-[#c86b3c]">@{user.userId}</p>
            <p className="mt-2 font-[var(--font-display)] text-xl">
              {user.userId}
            </p>
            <p className="mt-2 text-sm text-[#5c5d63]">{user.email}</p>
          </Link>
        ))}
      </section>
    </div>
  );
}
