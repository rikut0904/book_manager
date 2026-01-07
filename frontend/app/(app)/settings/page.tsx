 "use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { fetchJSON } from "@/lib/api";
import { getAuthState } from "@/lib/auth";

type UserResponse = {
  isAdmin?: boolean;
};

export default function SettingsPage() {
  const [isAdmin, setIsAdmin] = useState(false);

  useEffect(() => {
    const auth = getAuthState();
    if (!auth?.userId) {
      return;
    }
    fetchJSON<UserResponse>(`/users/${auth.userId}`, { auth: true })
      .then((data) => {
        setIsAdmin(Boolean(data.isAdmin));
      })
      .catch(() => {
        setIsAdmin(false);
      });
  }, []);
  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
          Settings
        </p>
        <h1 className="mt-2 font-[var(--font-display)] text-3xl">設定</h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          シリーズ・プロフィール・AIの設定へ移動します。
        </p>
      </section>

      <section className="grid gap-6 lg:grid-cols-[1fr_1fr]">
        <Link
          className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm transition hover:-translate-y-1 hover:shadow-md"
          href="/settings/series"
        >
          <h2 className="font-[var(--font-display)] text-2xl">シリーズ管理</h2>
          <p className="mt-2 text-sm text-[#5c5d63]">
            シリーズの追加・確認を行います。
          </p>
        </Link>

        <Link
          className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm transition hover:-translate-y-1 hover:shadow-md"
          href="/settings/profile"
        >
          <h2 className="font-[var(--font-display)] text-2xl">
            プロフィール公開
          </h2>
          <p className="mt-2 text-sm text-[#5c5d63]">
            公開範囲の変更と表示設定。
          </p>
        </Link>

        {isAdmin ? (
          <Link
            className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm transition hover:-translate-y-1 hover:shadow-md"
            href="/settings/ai"
          >
            <h2 className="font-[var(--font-display)] text-2xl">AI設定</h2>
            <p className="mt-2 text-sm text-[#5c5d63]">
              OpenAIシリーズ推定の有効化とモデル選択。
            </p>
          </Link>
        ) : null}
      </section>
    </div>
  );
}
