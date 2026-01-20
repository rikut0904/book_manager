"use client";

import Link from "next/link";
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

type UserBook = {
  id: string;
  bookId: string;
  seriesId: string;
  volumeNumber: number;
};

export default function BookmarkPage() {
  const [items, setItems] = useState<Favorite[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [books, setBooks] = useState<Book[]>([]);
  const [series, setSeries] = useState<{ id: string; name: string }[]>([]);
  const [userBooks, setUserBooks] = useState<UserBook[]>([]);

  const loadBookmarks = () => {
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
        setError("ブックマークを取得できませんでした。");
      });
    return () => {
      isMounted = false;
    };
  };

  useEffect(() => {
    const cleanup = loadBookmarks();
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

  useEffect(() => {
    let isMounted = true;
    fetchJSON<{ items: UserBook[] }>("/user-books", { auth: true })
      .then((data) => {
        if (!isMounted) {
          return;
        }
        setUserBooks(data.items ?? []);
      })
      .catch(() => {
        if (!isMounted) {
          return;
        }
        setUserBooks([]);
      });
    return () => {
      isMounted = false;
    };
  }, []);

  const handleDelete = async (id: string) => {
    try {
      await fetchJSON(`/favorites/${id}`, { method: "DELETE", auth: true });
      loadBookmarks();
    } catch {
      setError("ブックマークの削除に失敗しました。");
    }
  };

  const getBookTitle = (id: string) =>
    books.find((item) => item.id === id)?.title || id;
  const getSeriesName = (id: string) =>
    series.find((item) => item.id === id)?.name || id;
  const getVolumeLabel = (id: string) => {
    const volume = userBooks.find((item) => item.bookId === id)?.volumeNumber;
    if (!volume) {
      return "";
    }
    return `Vol.${volume}`;
  };
  const getBookTitleWithVolume = (id: string) => {
    const title = getBookTitle(id);
    const volume = getVolumeLabel(id);
    return volume ? `${title} ${volume}` : title;
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
          Bookmarks
        </p>
        <h1 className="mt-2 font-[var(--font-display)] text-3xl">
          ブックマーク
        </h1>
      </section>
      <section className="grid gap-4 md:grid-cols-2">
        {error ? (
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 text-sm text-red-600">
            {error}
          </div>
        ) : null}
        {!error && items.length === 0 ? (
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 text-sm text-[#5c5d63]">
            まだブックマークがありません。
          </div>
        ) : null}
        {items.map((item) => (
          <div
            key={item.id}
            className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 shadow-sm"
          >
            <div className="flex items-start justify-between gap-3">
              <div>
                <p className="text-xs text-[#c86b3c]">
                  {item.type === "series" ? "シリーズ" : "単巻"}
                </p>
                <p className="mt-2 font-[var(--font-display)] text-xl">
                  {item.type === "series"
                    ? getSeriesName(item.seriesId)
                    : getBookTitleWithVolume(item.bookId)}
                </p>
              </div>
              <button
                aria-label="ブックマーク解除"
                className="flex h-9 w-9 items-center justify-center rounded-full border border-[#c86b3c] bg-[#c86b3c] text-sm text-white transition hover:opacity-80"
                type="button"
                onClick={() => handleDelete(item.id)}
              >
                ★
              </button>
            </div>
          </div>
        ))}
      </section>
    </div>
  );
}
