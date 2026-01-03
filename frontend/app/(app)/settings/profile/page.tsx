"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { fetchJSON } from "@/lib/api";
import { getAuthState } from "@/lib/auth";

export default function ProfileSettingsPage() {
  const [visibility, setVisibility] = useState("public");
  const [settingError, setSettingError] = useState<string | null>(null);

  useEffect(() => {
    const auth = getAuthState();
    if (!auth?.userId) {
      return;
    }
    fetchJSON<{ settings: { visibility: string } }>(`/users/${auth.userId}`, {
      auth: true,
    })
      .then((data) => {
        if (data.settings?.visibility) {
          setVisibility(data.settings.visibility);
        }
      })
      .catch(() => {
        setSettingError("公開設定を取得できませんでした。");
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
        <Link className="text-xs text-[#5c5d63]" href="/settings">
          ← 設定に戻る
        </Link>
        <h1 className="mt-3 font-[var(--font-display)] text-3xl">
          プロフィール公開
        </h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          公開範囲の変更と表示設定。
        </p>
      </section>

      <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
        <div className="flex flex-col gap-3 text-sm">
          <label className="flex items-center justify-between rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3">
            <span>公開</span>
            <input
              type="radio"
              name="visibility"
              value="public"
              checked={visibility === "public"}
              onChange={() => setVisibility("public")}
            />
          </label>
          <label className="flex items-center justify-between rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3">
            <span>フォロワー限定</span>
            <input
              type="radio"
              name="visibility"
              value="followers"
              checked={visibility === "followers"}
              onChange={() => setVisibility("followers")}
            />
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
    </div>
  );
}
