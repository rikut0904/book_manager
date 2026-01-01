"use client";

import { useEffect, useState } from "react";

import { fetchJSON } from "@/lib/api";

type Recommendation = {
  id: string;
  userId: string;
  bookId: string;
  comment: string;
  createdAt: string;
};

type Book = {
  id: string;
  title: string;
};

export default function RecommendationsPage() {
  const [items, setItems] = useState<Recommendation[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [form, setForm] = useState({ bookId: "", comment: "" });
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [books, setBooks] = useState<Book[]>([]);

  const loadItems = () => {
    let isMounted = true;
    fetchJSON<{ items: Recommendation[] }>("/recommendations", { auth: true })
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
        setError("おすすめを取得できませんでした。");
      });
    return () => {
      isMounted = false;
    };
  };

  useEffect(() => {
    const cleanup = loadItems();
    return () => cleanup();
  }, []);

  useEffect(() => {
    let isMounted = true;
    fetchJSON<{ items: Book[] }>("/books", { auth: true })
      .then((data) => {
        if (!isMounted) {
          return;
        }
        setBooks(data.items ?? []);
      })
      .catch(() => {
        if (!isMounted) {
          return;
        }
        setBooks([]);
      });
    return () => {
      isMounted = false;
    };
  }, []);

  const handleCreate = async () => {
    setSubmitError(null);
    if (!form.bookId.trim()) {
      setSubmitError("bookIdを入力してください。");
      return;
    }
    if (!form.comment.trim()) {
      setSubmitError("コメントを入力してください。");
      return;
    }
    try {
      await fetchJSON<Recommendation>("/recommendations", {
        method: "POST",
        auth: true,
        body: JSON.stringify({
          bookId: form.bookId.trim(),
          comment: form.comment.trim(),
        }),
      });
      setForm({ bookId: "", comment: "" });
      loadItems();
    } catch {
      setSubmitError("おすすめの投稿に失敗しました。");
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await fetchJSON(`/recommendations/${id}`, { method: "DELETE", auth: true });
      loadItems();
    } catch {
      setError("削除に失敗しました。");
    }
  };

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
          Recommendations
        </p>
        <h1 className="mt-2 font-[var(--font-display)] text-3xl">
          おすすめ一覧
        </h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          みんなのおすすめがタイムライン形式で表示されます。
        </p>
      </section>

      <section className="grid gap-4">
        {error ? (
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 text-sm text-red-600">
            {error}
          </div>
        ) : null}
        {!error && items.length === 0 ? (
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 text-sm text-[#5c5d63]">
            まだおすすめが投稿されていません。
          </div>
        ) : null}
        {items.map((item) => (
          <div
            key={item.id}
            className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm"
          >
            <p className="text-xs text-[#c86b3c]">@{item.userId}</p>
            <p className="mt-2 font-[var(--font-display)] text-xl">
              {item.bookId}
            </p>
            <p className="mt-2 text-sm text-[#5c5d63]">{item.comment}</p>
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
        <h2 className="font-[var(--font-display)] text-2xl">おすすめ投稿</h2>
        <div className="mt-4 flex flex-col gap-3 text-sm">
          <select
            className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
            value={form.bookId}
            onChange={(event) =>
              setForm((prev) => ({ ...prev, bookId: event.target.value }))
            }
          >
            <option value="">書籍を選択</option>
            {books.map((book) => (
              <option key={book.id} value={book.id}>
                {book.title || book.id}
              </option>
            ))}
          </select>
          <input
            className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-2 text-xs outline-none transition focus:border-[#c86b3c]"
            placeholder="bookId を直接入力"
            value={form.bookId}
            onChange={(event) =>
              setForm((prev) => ({ ...prev, bookId: event.target.value }))
            }
          />
          <textarea
            className="min-h-[120px] rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
            placeholder="コメント"
            value={form.comment}
            onChange={(event) =>
              setForm((prev) => ({ ...prev, comment: event.target.value }))
            }
          />
          {submitError ? (
            <p className="text-xs text-red-600">{submitError}</p>
          ) : null}
          <button
            className="rounded-full bg-[#c86b3c] px-5 py-3 text-sm font-medium text-white"
            type="button"
            onClick={handleCreate}
          >
            投稿する
          </button>
        </div>
      </section>
    </div>
  );
}
