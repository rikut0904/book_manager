"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { fetchJSON } from "@/lib/api";
import { getAuthState } from "@/lib/auth";

type UserResponse = {
  isAdmin?: boolean;
  user?: { userId?: string; displayName?: string };
  settings?: { visibility?: string };
};

export default function SettingsPage() {
  const [isAdmin, setIsAdmin] = useState(false);
  const [username, setUsername] = useState("");
  const [visibility, setVisibility] = useState("public");
  const [settingError, setSettingError] = useState<string | null>(null);

  useEffect(() => {
    const auth = getAuthState();
    if (!auth?.userId) {
      return;
    }
    fetchJSON<UserResponse>(`/users/${auth.userId}`, { auth: true })
      .then((data) => {
        setIsAdmin(Boolean(data.isAdmin));
        const name = data.user?.displayName?.trim() || data.user?.userId?.trim() || "";
        setUsername(name);
        if (data.settings?.visibility) {
          setVisibility(data.settings.visibility);
        }
      })
      .catch(() => {
        setIsAdmin(false);
        setUsername("");
        setSettingError("設定情報を取得できませんでした。");
      });
  }, []);

  const handleSaveVisibility = async () => {
    setSettingError(null);
    try {
      await fetchJSON("/users/me/settings", {
        method: "PATCH",
        auth: true,
        body: JSON.stringify({ visibility }),
      });
    } catch {
      setSettingError("公開設定の保存に失敗しました。");
    }
  };
  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <Link className="inline-flex items-center gap-2 text-xs text-[#5c5d63]" href="/user">
          <svg
            aria-hidden="true"
            className="h-4 w-4"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <path d="M15 6l-6 6 6 6" />
          </svg>
          ユーザーに戻る
        </Link>
        <p className="mt-4 text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
          Settings
        </p>
        <h1 className="mt-2 font-[var(--font-display)] text-3xl">設定</h1>
      </section>

      <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
        <h2 className="font-[var(--font-display)] text-2xl">公開設定</h2>
        <p className="mt-2 text-sm text-[#5c5d63]">
          プロフィールの公開範囲を選択できます。
        </p>
        <div className="mt-4 flex items-center justify-between rounded-2xl border border-[#e4d8c7] bg-white px-4 py-4 text-sm">
          <div>
            <p className="text-xs text-[#5c5d63]">公開範囲</p>
            <p className="mt-1 text-[#1b1c1f]">
              {visibility === "followers" ? "フォロワー限定" : "公開"}
            </p>
          </div>
          <label className="relative inline-flex cursor-pointer items-center">
            <input
              className="peer sr-only"
              type="checkbox"
              checked={visibility === "followers"}
              onChange={(event) =>
                setVisibility(event.target.checked ? "followers" : "public")
              }
            />
            <span className="h-6 w-11 rounded-full bg-[#efe5d4] transition peer-checked:bg-[#1b1c1f]" />
            <span className="absolute left-1 top-1 h-4 w-4 rounded-full bg-white transition peer-checked:translate-x-5" />
          </label>
        </div>
        {settingError ? (
          <p className="mt-3 text-xs text-red-600">{settingError}</p>
        ) : null}
        <button
          className="mt-4 rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white"
          type="button"
          onClick={handleSaveVisibility}
        >
          変更を保存
        </button>
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
      </section>

      {isAdmin ? (
        <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
          <h2 className="font-[var(--font-display)] text-2xl">管理者メニュー</h2>
          <p className="mt-2 text-sm text-[#5c5d63]">
            管理者向けの設定へ移動します。
          </p>
          <div className="mt-4 grid gap-3 text-sm md:grid-cols-3">
            <Link
              className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-4 text-[#1b1c1f] transition hover:-translate-y-0.5 hover:shadow-sm"
              href="/settings/admin-users"
            >
              <p className="text-xs text-[#5c5d63]">管理者管理</p>
              <p className="mt-1 font-medium">権限の付与・削除</p>
            </Link>
            <Link
              className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-4 text-[#1b1c1f] transition hover:-translate-y-0.5 hover:shadow-sm"
              href="/settings/invite"
            >
              <p className="text-xs text-[#5c5d63]">管理者招待</p>
              <p className="mt-1 font-medium">招待トークンの発行</p>
            </Link>
            <Link
              className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-4 text-[#1b1c1f] transition hover:-translate-y-0.5 hover:shadow-sm"
              href="/settings/ai"
            >
              <p className="text-xs text-[#5c5d63]">AI設定</p>
              <p className="mt-1 font-medium">モデルの選択</p>
            </Link>
          </div>
        </section>
      ) : null}
    </div>
  );
}
