"use client";

import Link from "next/link";
import type { MouseEvent } from "react";
import { useEffect, useMemo, useState } from "react";

import { fetchJSON } from "@/lib/api";

type Book = {
  id: string;
  title: string;
  authors: string[];
  publisher: string;
  publishedDate: string;
  thumbnailUrl: string;
  isbn13: string;
  source: string;
  seriesName?: string;
};

type UserBook = {
  id: string;
  bookId: string;
  seriesId: string;
  volumeNumber: number;
};

type Series = {
  id: string;
  name: string;
};

type Favorite = {
  id: string;
  type: "book" | "series";
  bookId: string;
  seriesId: string;
};

export default function BooksPage() {
  const [items, setItems] = useState<Book[]>([]);
  const [userBooks, setUserBooks] = useState<UserBook[]>([]);
  const [seriesList, setSeriesList] = useState<Series[]>([]);
  const [favorites, setFavorites] = useState<Favorite[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let isMounted = true;
    const load = async () => {
      try {
        const [booksRes, userBooksRes, seriesRes, favoritesRes] =
          await Promise.all([
            fetchJSON<{ items: Book[] }>("/books", { auth: true }),
            fetchJSON<{ items: UserBook[] }>("/user-books", { auth: true }).catch(
              () => ({ items: [] })
            ),
            fetchJSON<{ items: Series[] }>("/series", { auth: true }).catch(
              () => ({ items: [] })
            ),
            fetchJSON<{ items: Favorite[] }>("/favorites", { auth: true }).catch(
              () => ({ items: [] })
            ),
          ]);
        if (!isMounted) {
          return;
        }
        setItems(booksRes.items ?? []);
        setUserBooks(userBooksRes.items ?? []);
        setSeriesList(seriesRes.items ?? []);
        setFavorites(favoritesRes.items ?? []);
      } catch {
        if (!isMounted) {
          return;
        }
        setError("書籍一覧を取得できませんでした。");
      } finally {
        if (!isMounted) {
          return;
        }
        setIsLoading(false);
      }
    };
    load();
    return () => {
      isMounted = false;
    };
  }, []);

  const favoritesBySeriesId = useMemo(() => {
    const map = new Map<string, Favorite>();
    favorites
      .filter((item) => item.type === "series" && item.seriesId)
      .forEach((item) => {
        map.set(item.seriesId, item);
      });
    return map;
  }, [favorites]);

  const volumesBySeriesId = useMemo(() => {
    const map = new Map<string, number[]>();
    userBooks.forEach((item) => {
      if (!item.seriesId || !item.volumeNumber) {
        return;
      }
      const list = map.get(item.seriesId) ?? [];
      list.push(item.volumeNumber);
      map.set(item.seriesId, list);
    });
    map.forEach((list, key) => {
      const unique = Array.from(new Set(list)).sort((a, b) => a - b);
      map.set(key, unique);
    });
    return map;
  }, [userBooks]);

  const getVolumeSummary = (seriesId: string) => {
    const list = volumesBySeriesId.get(seriesId);
    if (!list || list.length === 0) {
      return "";
    }
    if (list.length === 1) {
      return `Vol.${list[0]}`;
    }
    return `Vol.${list[0]}〜${list[list.length - 1]}`;
  };

  const { seriesCards, singleBooks } = useMemo(() => {
    const booksById = new Map(items.map((book) => [book.id, book]));
    const seriesById = new Map(seriesList.map((series) => [series.id, series]));
    const userBookIds = new Set(userBooks.map((ub) => ub.bookId));
    const grouped = new Map<
      string,
      { seriesId: string; name: string; books: Book[] }
    >();
    const seriesBookIDs = new Set<string>();
    userBooks.forEach((item) => {
      if (!item.seriesId) {
        return;
      }
      const book = booksById.get(item.bookId);
      if (!book) {
        return;
      }
      seriesBookIDs.add(book.id);
      const existing = grouped.get(item.seriesId);
      if (existing) {
        existing.books.push(book);
      } else {
        grouped.set(item.seriesId, {
          seriesId: item.seriesId,
          name:
            seriesById.get(item.seriesId)?.name ||
            book.seriesName ||
            "未判定",
          books: [book],
        });
      }
    });
    const seriesCards = Array.from(grouped.values()).sort((a, b) =>
      a.name.localeCompare(b.name, "ja")
    );
    const singleBooks = items.filter(
      (book) => userBookIds.has(book.id) && !seriesBookIDs.has(book.id)
    );
    return { seriesCards, singleBooks };
  }, [items, seriesList, userBooks]);

  const handleToggleSeriesFavorite = async (
    seriesId: string,
    event: MouseEvent<HTMLButtonElement>
  ) => {
    event.preventDefault();
    event.stopPropagation();
    const existing = favoritesBySeriesId.get(seriesId);
    try {
      if (existing) {
        await fetchJSON(`/favorites/${existing.id}`, {
          method: "DELETE",
          auth: true,
        });
        setFavorites((prev) => prev.filter((item) => item.id !== existing.id));
      } else {
        const created = await fetchJSON<Favorite>("/favorites", {
          method: "POST",
          auth: true,
          body: JSON.stringify({
            type: "series",
            seriesId,
          }),
        });
        setFavorites((prev) => [...prev, created]);
      }
    } catch {
      setError("ブックマークの更新に失敗しました。");
    }
  };

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
              Library
            </p>
            <h1 className="mt-2 font-[var(--font-display)] text-3xl">
              ホーム
            </h1>
          </div>
          <div className="flex flex-wrap gap-3">
            <Link
              className="rounded-full bg-[#1b1c1f] px-4 py-2 text-xs font-medium text-white"
              href="/books/new"
            >
              書籍登録
            </Link>
          </div>
        </div>
        <div className="mt-6 flex flex-wrap gap-3">
          <input
            className="w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c] md:flex-1"
            placeholder="タイトル・著者・ISBNで検索"
          />
          <select className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm text-[#5c5d63]">
            <option>全て</option>
            <option>ブックマーク</option>
            <option>最近追加</option>
          </select>
        </div>
      </section>

      <section className="grid gap-4 lg:grid-cols-3">
        {isLoading ? (
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 text-sm text-[#5c5d63]">
            読み込み中...
          </div>
        ) : null}
        {error ? (
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 text-sm text-red-600">
            {error}
          </div>
        ) : null}
        {!isLoading && !error && items.length === 0 ? (
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 text-sm text-[#5c5d63]">
            まだ登録された書籍がありません。
          </div>
        ) : null}
        {seriesCards.length > 0 ? (
          <div className="lg:col-span-3">
            <div className="rounded-2xl border border-[#e4d8c7] bg-white/80 px-4 py-3 text-sm text-[#5c5d63]">
              シリーズ
            </div>
          </div>
        ) : null}
        {seriesCards.map((series) => {
          const favorite = favoritesBySeriesId.get(series.seriesId);
          const authors = Array.from(
            new Set(
              series.books.flatMap((book) => book.authors || []).filter(Boolean)
            )
          );
          const volumeSummary = getVolumeSummary(series.seriesId);
          return (
            <Link
              key={series.seriesId}
              className="group rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 shadow-sm transition hover:-translate-y-1 hover:shadow-md"
              href={`/books/series/${series.seriesId}`}
            >
              <div className="flex items-start justify-between gap-3 text-xs text-[#5c5d63]">
                <span>シリーズ</span>
                <button
                  aria-label={
                    favorite ? "シリーズのブックマーク解除" : "シリーズをブックマーク"
                  }
                  className={`flex h-8 w-8 items-center justify-center rounded-full border text-sm transition ${
                    favorite
                      ? "border-[#c86b3c] bg-[#c86b3c] text-white"
                      : "border-[#e4d8c7] text-[#5c5d63] hover:bg-white"
                  }`}
                  type="button"
                  onClick={(event) =>
                    handleToggleSeriesFavorite(series.seriesId, event)
                  }
                >
                  {favorite ? "★" : "☆"}
                </button>
              </div>
              <h2 className="mt-4 font-[var(--font-display)] text-xl text-[#1b1c1f]">
                {series.name}
              </h2>
              <p className="mt-2 text-sm text-[#5c5d63]">
                {authors.length > 0 ? authors.join(" / ") : "著者未登録"}
              </p>
              <div className="mt-4 flex flex-wrap items-center gap-2 text-xs text-[#5c5d63]">
                <span className="rounded-2xl bg-[#f6f1e7] px-3 py-2">
                  巻数合計: {series.books.length}
                </span>
                {volumeSummary ? (
                  <span className="rounded-2xl bg-[#f6f1e7] px-3 py-2">
                    {volumeSummary}
                  </span>
                ) : null}
              </div>
              <p className="mt-4 text-xs text-[#c86b3c]">
                シリーズ詳細を見る →
              </p>
            </Link>
          );
        })}
        {singleBooks.length > 0 ? (
          <div className="lg:col-span-3">
            <div className="rounded-2xl border border-[#e4d8c7] bg-white/80 px-4 py-3 text-sm text-[#5c5d63]">
              単巻
            </div>
          </div>
        ) : null}
        {singleBooks.map((book, index) => (
          <Link
            key={book.id || book.isbn13 || `book-${index}`}
            className="group rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 shadow-sm transition hover:-translate-y-1 hover:shadow-md"
            href={`/books/${book.id}`}
          >
            <div className="flex items-center justify-between text-xs text-[#5c5d63]">
              <span>{book.publisher || "出版社未登録"}</span>
              <span>{book.publishedDate || "発売日未登録"}</span>
            </div>
            <h2 className="mt-4 font-[var(--font-display)] text-xl text-[#1b1c1f]">
              {book.title}
            </h2>
            <p className="mt-2 text-sm text-[#5c5d63]">
              {book.authors?.join(" / ") || "著者未登録"}
            </p>
            <div className="mt-4 rounded-2xl bg-[#f6f1e7] px-3 py-2 text-xs text-[#5c5d63]">
              ISBN: {book.isbn13 || "未登録"}
            </div>
            <p className="mt-4 text-xs text-[#c86b3c]">詳細を見る →</p>
          </Link>
        ))}
      </section>
    </div>
  );
}
