"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { fetchJSON } from "@/lib/api";

type Tag = {
  id: string;
  name: string;
};

export default function TagSettingsPage() {
  const [tags, setTags] = useState<Tag[]>([]);
  const [tagName, setTagName] = useState("");
  const [tagError, setTagError] = useState<string | null>(null);

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

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <Link className="text-xs text-[#5c5d63]" href="/settings">
          ← 設定に戻る
        </Link>
        <h1 className="mt-3 font-[var(--font-display)] text-3xl">タグ管理</h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          タグの追加・削除を行います。
        </p>
      </section>

      <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
        <div className="flex flex-wrap gap-2 text-xs">
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
      </section>
    </div>
  );
}
