"use client";

import { useEffect, useState } from "react";

import { fetchJSON } from "@/lib/api";
import { getAuthState } from "@/lib/auth";

type Tag = {
  id: string;
  name: string;
};

export default function SettingsPage() {
  const [tags, setTags] = useState<Tag[]>([]);
  const [tagName, setTagName] = useState("");
  const [tagError, setTagError] = useState<string | null>(null);
  const [visibility, setVisibility] = useState("public");
  const [settingError, setSettingError] = useState<string | null>(null);

  useEffect(() => {
    let isMounted = true;
    fetchJSON<{ items: Tag[] }>("/tags", { auth: true })
      .then((data) => {
        if (!isMounted) {
          return;
        }
        setTags(data.items ?? []);
      })
      .catch(() => {
        if (!isMounted) {
          return;
        }
        setTagError("タグ一覧を取得できませんでした。");
      });
    return () => {
      isMounted = false;
    };
  }, []);

  useEffect(() => {
    const auth = getAuthState();
    if (!auth?.userId) {
      return;
    }
    fetchJSON<{
      settings: { visibility: string };
    }>(`/users/${auth.userId}`, { auth: true })
      .then((data) => {
        if (data.settings?.visibility) {
          setVisibility(data.settings.visibility);
        }
      })
      .catch(() => {
        setSettingError("公開設定を取得できませんでした。");
      });
  }, []);

  const handleCreateTag = async () => {
    setTagError(null);
    if (!tagName.trim()) {
      setTagError("タグ名を入力してください。");
      return;
    }
    try {
      const tag = await fetchJSON<Tag>("/tags", {
        method: "POST",
        auth: true,
        body: JSON.stringify({ name: tagName.trim() }),
      });
      setTags((prev) => [...prev, tag]);
      setTagName("");
    } catch {
      setTagError("タグ作成に失敗しました。");
    }
  };

  const handleDeleteTag = async (tagId: string) => {
    setTagError(null);
    try {
      await fetchJSON(`/tags/${tagId}`, { method: "DELETE", auth: true });
      setTags((prev) => prev.filter((tag) => tag.id !== tagId));
    } catch {
      setTagError("タグ削除に失敗しました。");
    }
  };

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
            {tags.map((tag) => (
              <span
                key={tag.id}
                className="flex items-center gap-2 rounded-full border border-[#e4d8c7] bg-white px-3 py-2 text-[#5c5d63]"
              >
                {tag.name}
                <button
                  className="text-[10px] text-[#c86b3c] hover:text-[#8f3d1f]"
                  type="button"
                  onClick={() => handleDeleteTag(tag.id)}
                >
                  削除
                </button>
              </span>
            ))}
          </div>
          <div className="mt-4 flex flex-wrap gap-2">
            <input
              className="min-w-[200px] rounded-full border border-[#e4d8c7] bg-white px-4 py-2 text-xs outline-none transition focus:border-[#c86b3c]"
              placeholder="新規タグ名"
              value={tagName}
              onChange={(event) => setTagName(event.target.value)}
            />
            <button
              className="rounded-full bg-[#efe5d4] px-3 py-2 text-xs text-[#1b1c1f]"
              type="button"
              onClick={handleCreateTag}
            >
              追加
            </button>
          </div>
          {tagError ? (
            <p className="mt-3 text-xs text-red-600">{tagError}</p>
          ) : null}
        </div>

        <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
          <h2 className="font-[var(--font-display)] text-2xl">
            プロフィール公開
          </h2>
          <div className="mt-4 flex flex-col gap-3 text-sm">
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
        </div>
      </section>
    </div>
  );
}
