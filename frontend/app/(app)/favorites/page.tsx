"use client";

import { useEffect, useState } from "react";

import { fetchJSON } from "@/lib/api";

type Favorite = {
  id: string;
  type: "book" | "series";
  bookId: string;
  seriesId: string;
};

type Book = {
  id: string;
  title: string;
};

export default function FavoritesPage() {
  const [items, setItems] = useState<Favorite[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [form, setForm] = useState({
    type: "book",
    bookId: "",
    seriesId: "",
  });
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [books, setBooks] = useState<Book[]>([]);
  const [series, setSeries] = useState<{ id: string; name: string }[]>([]);

  const loadFavorites = () => {
    let isMounted = true;
    fetchJSON<{ items: Favorite[] }>("/favorites", { auth: true })
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
        setError("お気に入りを取得できませんでした。");
      });
    return () => {
      isMounted = false;
    };
  };

  useEffect(() => {
    const cleanup = loadFavorites();
    return () => cleanup();
  }, []);

  useEffect(() => {
    let isMounted = true;
    fetchJSON<{ items: { id: string; name: string }[] }>("/series", {
      auth: true,
    })
      .then((data) => {
        if (!isMounted) {
          return;
        }
        setSeries(data.items ?? []);
      })
      .catch(() => {
        if (!isMounted) {
          return;
        }
        setSeries([]);
      });
    return () => {
      isMounted = false;
    };
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
    if (form.type === "book" && !form.bookId.trim()) {
      setSubmitError("bookIdを入力してください。");
      return;
    }
    if (form.type === "series" && !form.seriesId.trim()) {
      setSubmitError("seriesIdを入力してください。");
      return;
    }
    try {
      await fetchJSON<Favorite>("/favorites", {
        method: "POST",
        auth: true,
        body: JSON.stringify({
          type: form.type,
          bookId: form.type === "book" ? form.bookId.trim() : "",
          seriesId: form.type === "series" ? form.seriesId.trim() : "",
        }),
      });
      setForm({ type: "book", bookId: "", seriesId: "" });
      loadFavorites();
    } catch {
      setSubmitError("お気に入りの登録に失敗しました。");
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await fetchJSON(`/favorites/${id}`, { method: "DELETE", auth: true });
      loadFavorites();
    } catch {
      setError("お気に入りの削除に失敗しました。");
    }
  };

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
          Favorites
        </p>
        <h1 className="mt-2 font-[var(--font-display)] text-3xl">
          お気に入り
        </h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          単巻とシリーズのお気に入りをまとめて管理します。
        </p>
      </section>
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 shadow-sm">
        <h2 className="font-[var(--font-display)] text-xl">お気に入り登録</h2>
        <div className="mt-4 grid gap-3 text-sm md:grid-cols-[120px_1fr]">
          <label className="text-[#1b1c1f]">
            種別
            <select
              className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-3 py-2 text-sm"
              value={form.type}
              onChange={(event) =>
                setForm((prev) => ({
                  ...prev,
                  type: event.target.value as "book" | "series",
                }))
              }
            >
              <option value="book">単巻</option>
              <option value="series">シリーズ</option>
            </select>
          </label>
          <label className="text-[#1b1c1f]">
            {form.type === "book" ? "bookId" : "seriesId"}
            {form.type === "book" ? (
              <select
                className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
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
            ) : (
              <select
                className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
                value={form.seriesId}
                onChange={(event) =>
                  setForm((prev) => ({ ...prev, seriesId: event.target.value }))
                }
              >
                <option value="">シリーズを選択</option>
                {series.map((item) => (
                  <option key={item.id} value={item.id}>
                    {item.name}
                  </option>
                ))}
              </select>
            )}
            <input
              className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-2 text-xs outline-none transition focus:border-[#c86b3c]"
              placeholder={form.type === "book" ? "bookId を直接入力" : "seriesId を直接入力"}
              value={form.type === "book" ? form.bookId : form.seriesId}
              onChange={(event) =>
                setForm((prev) => ({
                  ...prev,
                  [form.type === "book" ? "bookId" : "seriesId"]:
                    event.target.value,
                }))
              }
            />
          </label>
        </div>
        {submitError ? (
          <p className="mt-3 text-xs text-red-600">{submitError}</p>
        ) : null}
        <button
          className="mt-4 rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white"
          type="button"
          onClick={handleCreate}
        >
          登録する
        </button>
      </section>
      <section className="grid gap-4 md:grid-cols-2">
        {error ? (
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 text-sm text-red-600">
            {error}
          </div>
        ) : null}
        {!error && items.length === 0 ? (
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 text-sm text-[#5c5d63]">
            まだお気に入りがありません。
          </div>
        ) : null}
        {items.map((item) => (
          <div
            key={item.id}
            className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 shadow-sm"
          >
            <p className="text-xs text-[#c86b3c]">
              {item.type === "series" ? "シリーズ" : "単巻"}
            </p>
            <p className="mt-2 font-[var(--font-display)] text-xl">
              {item.type === "series" ? item.seriesId : item.bookId}
            </p>
            <p className="mt-2 text-sm text-[#5c5d63]">
              IDベースで管理中
            </p>
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
    </div>
  );
}
