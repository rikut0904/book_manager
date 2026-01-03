"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { fetchJSON } from "@/lib/api";

type Series = {
  id: string;
  name: string;
};

export default function SeriesSettingsPage() {
  const [items, setItems] = useState<Series[]>([]);
  const [seriesName, setSeriesName] = useState("");
  const [seriesError, setSeriesError] = useState<string | null>(null);

  useEffect(() => {
    let isMounted = true;
    fetchJSON<{ items: Series[] }>("/series", { auth: true })
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
        setSeriesError("シリーズ一覧を取得できませんでした。");
      });
    return () => {
      isMounted = false;
    };
  }, []);

  const handleCreateSeries = async () => {
    setSeriesError(null);
    if (!seriesName.trim()) {
      setSeriesError("シリーズ名を入力してください。");
      return;
    }
    try {
      const item = await fetchJSON<Series>("/series", {
        method: "POST",
        auth: true,
        body: JSON.stringify({ name: seriesName.trim() }),
      });
      setItems((prev) => [...prev, item]);
      setSeriesName("");
    } catch {
      setSeriesError("シリーズ作成に失敗しました。");
    }
  };

  const handleDeleteSeries = async (seriesId: string) => {
    setSeriesError(null);
    try {
      await fetchJSON(`/series/${seriesId}`, { method: "DELETE", auth: true });
      setItems((prev) => prev.filter((item) => item.id !== seriesId));
    } catch {
      setSeriesError("シリーズ削除に失敗しました。");
    }
  };

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <Link className="text-xs text-[#5c5d63]" href="/settings">
          ← 設定に戻る
        </Link>
        <h1 className="mt-3 font-[var(--font-display)] text-3xl">
          シリーズ管理
        </h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          シリーズの追加・確認を行います。
        </p>
      </section>

      <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
        <div className="flex flex-wrap gap-2 text-xs">
          {items.map((item) => (
            <span
              key={item.id}
              className="flex items-center gap-2 rounded-full border border-[#e4d8c7] bg-white px-3 py-2 text-[#5c5d63]"
            >
              {item.name}
              <button
                className="text-[10px] text-[#c86b3c] hover:text-[#8f3d1f]"
                type="button"
                onClick={() => handleDeleteSeries(item.id)}
              >
                削除
              </button>
            </span>
          ))}
        </div>
        <div className="mt-4 flex flex-wrap gap-2">
          <input
            className="min-w-[200px] rounded-full border border-[#e4d8c7] bg-white px-4 py-2 text-xs outline-none transition focus:border-[#c86b3c]"
            placeholder="新規シリーズ名"
            value={seriesName}
            onChange={(event) => setSeriesName(event.target.value)}
          />
          <button
            className="rounded-full bg-[#efe5d4] px-3 py-2 text-xs text-[#1b1c1f]"
            type="button"
            onClick={handleCreateSeries}
          >
            追加
          </button>
        </div>
        {seriesError ? (
          <p className="mt-3 text-xs text-red-600">{seriesError}</p>
        ) : null}
      </section>
    </div>
  );
}
