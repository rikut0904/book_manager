"use client";

import { useParams } from "next/navigation";
import { useEffect, useState } from "react";

import { fetchJSON } from "@/lib/api";

type UserDetail = {
  user: { id: string; email: string; userId: string };
  settings: { visibility: string };
  stats: { ownedCount: number; seriesCount: number; followers: number; following: number };
  recommendations: string[];
};

export default function UserDetailPage() {
  const params = useParams<{ id: string }>();
  const [detail, setDetail] = useState<UserDetail | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!params?.id) {
      return;
    }
    let isMounted = true;
    fetchJSON<UserDetail>(`/users/${params.id}`, { auth: true })
      .then((data) => {
        if (!isMounted) {
          return;
        }
        setDetail(data);
      })
      .catch(() => {
        if (!isMounted) {
          return;
        }
        setError("プロフィールを取得できませんでした。");
      });
    return () => {
      isMounted = false;
    };
  }, [params?.id]);

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
          Profile
        </p>
        <h1 className="mt-2 font-[var(--font-display)] text-3xl">
          {detail?.user.userId || "ユーザー"}
        </h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          公開範囲: {detail?.settings.visibility || "public"}
        </p>
      </section>

      <section className="grid gap-4 md:grid-cols-3">
        {error ? (
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 text-sm text-red-600">
            {error}
          </div>
        ) : null}
        {[
          {
            label: "所蔵数",
            value: detail ? `${detail.stats.ownedCount}冊` : "--",
          },
          {
            label: "シリーズ",
            value: detail ? `${detail.stats.seriesCount}件` : "--",
          },
          {
            label: "フォロワー",
            value: detail ? `${detail.stats.followers}人` : "--",
          },
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
          {(detail?.recommendations ?? []).length === 0 ? (
            <div className="rounded-2xl border border-[#e4d8c7] bg-white p-4 text-sm text-[#5c5d63]">
              まだおすすめがありません。
            </div>
          ) : null}
          {(detail?.recommendations ?? []).map((title) => (
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
