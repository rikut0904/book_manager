"use client";

import { useEffect, useState } from "react";

import { fetchJSON } from "@/lib/api";

type NextToBuyItem = {
  id: string;
  title: string;
  seriesName: string;
  volumeNumber: number;
  note: string;
};

export default function NextToBuyPage() {
  const [items, setItems] = useState<NextToBuyItem[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [form, setForm] = useState({
    title: "",
    seriesName: "",
    volumeNumber: "",
    note: "",
  });
  const [submitError, setSubmitError] = useState<string | null>(null);

  const loadItems = () => {
    let isMounted = true;
    fetchJSON<{ items: NextToBuyItem[] }>("/next-to-buy", { auth: true })
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
        setError("次に買う本を取得できませんでした。");
      });
    return () => {
      isMounted = false;
    };
  };

  useEffect(() => {
    const cleanup = loadItems();
    return () => cleanup();
  }, []);

  const handleCreate = async () => {
    setSubmitError(null);
    if (!form.title.trim()) {
      setSubmitError("タイトルを入力してください。");
      return;
    }
    const volume =
      form.volumeNumber.trim() === ""
        ? undefined
        : Number(form.volumeNumber);
    try {
      await fetchJSON<NextToBuyItem>("/next-to-buy/manual", {
        method: "POST",
        auth: true,
        body: JSON.stringify({
          title: form.title.trim(),
          seriesName: form.seriesName.trim(),
          volumeNumber: Number.isFinite(volume) ? volume : undefined,
          note: form.note.trim(),
        }),
      });
      setForm({ title: "", seriesName: "", volumeNumber: "", note: "" });
      loadItems();
    } catch {
      setSubmitError("手動追加に失敗しました。");
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await fetchJSON(`/next-to-buy/manual/${id}`, {
        method: "DELETE",
        auth: true,
      });
      loadItems();
    } catch {
      setError("削除に失敗しました。");
    }
  };

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
          Next To Buy
        </p>
        <h1 className="mt-2 font-[var(--font-display)] text-3xl">
          次に買う本
        </h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          お気に入りの次巻と手入力メモをまとめます。
        </p>
      </section>

      <section className="grid gap-4 lg:grid-cols-3">
        {error ? (
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 text-sm text-red-600">
            {error}
          </div>
        ) : null}
        {!error && items.length === 0 ? (
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 text-sm text-[#5c5d63]">
            まだ登録がありません。
          </div>
        ) : null}
        {items.map((item) => (
          <div
            key={item.id}
            className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 shadow-sm"
          >
            <p className="font-[var(--font-display)] text-xl">
              {item.title}
            </p>
            <p className="mt-2 text-sm text-[#5c5d63]">
              {item.seriesName ? `${item.seriesName} ` : ""}
              {item.volumeNumber ? `Vol.${item.volumeNumber}` : ""}
            </p>
            <p className="mt-2 text-sm text-[#5c5d63]">{item.note}</p>
            <button
              className="mt-4 rounded-full border border-[#e4d8c7] px-4 py-2 text-xs text-[#5c5d63] hover:bg-white"
              type="button"
              onClick={() => handleDelete(item.id)}
            >
              削除
            </button>
          </div>
        ))}
      </section>

      <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
        <h2 className="font-[var(--font-display)] text-2xl">手動追加</h2>
        <div className="mt-4 grid gap-3 text-sm md:grid-cols-[1fr_1fr]">
          <input
            className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
            placeholder="タイトル"
            value={form.title}
            onChange={(event) =>
              setForm((prev) => ({ ...prev, title: event.target.value }))
            }
          />
          <input
            className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
            placeholder="シリーズ名 (任意)"
            value={form.seriesName}
            onChange={(event) =>
              setForm((prev) => ({ ...prev, seriesName: event.target.value }))
            }
          />
          <input
            className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
            placeholder="巻数"
            value={form.volumeNumber}
            onChange={(event) =>
              setForm((prev) => ({
                ...prev,
                volumeNumber: event.target.value,
              }))
            }
          />
          <input
            className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
            placeholder="メモ"
            value={form.note}
            onChange={(event) =>
              setForm((prev) => ({ ...prev, note: event.target.value }))
            }
          />
        </div>
        {submitError ? (
          <p className="mt-3 text-xs text-red-600">{submitError}</p>
        ) : null}
        <button
          className="mt-4 rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white"
          type="button"
          onClick={handleCreate}
        >
          追加
        </button>
      </section>
    </div>
  );
}
