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
  seriesId?: string | null;
  volumeNumber?: number | null;
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

type BooksPageData = {
  books: Book[];
  userBooks: UserBook[];
  series: Series[];
  favorites: Favorite[];
};

let cachedBooksPageData: BooksPageData | null = null;
let cachedBooksPagePromise: Promise<BooksPageData> | null = null;

function resetBooksCache() {
  cachedBooksPageData = null;
  cachedBooksPagePromise = null;
}

export default function BooksPage() {
  const [items, setItems] = useState<Book[]>([]);
  const [userBooks, setUserBooks] = useState<UserBook[]>([]);
  const [seriesList, setSeriesList] = useState<Series[]>([]);
  const [favorites, setFavorites] = useState<Favorite[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  useEffect(() => {
    const handleAuthChanged = () => {
      resetBooksCache();
    };
    window.addEventListener("auth-changed", handleAuthChanged);
    let isMounted = true;
    const load = async () => {
      try {
        let data = cachedBooksPageData;
        if (!data) {
          if (!cachedBooksPagePromise) {
            cachedBooksPagePromise = fetchJSON<BooksPageData>("/books/overview", {
              auth: true,
            });
          }
          data = await cachedBooksPagePromise;
          cachedBooksPageData = data;
        }
        cachedBooksPageData = data;
        if (!isMounted) {
          return;
        }
        setItems(data.books);
        setUserBooks(data.userBooks);
        setSeriesList(data.series);
        setFavorites(data.favorites);
      } catch {
        cachedBooksPagePromise = null;
        cachedBooksPageData = null;
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
      window.removeEventListener("auth-changed", handleAuthChanged);
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
      const normalizedSeriesId =
        item.seriesId && item.seriesId !== "null" && item.seriesId !== "undefined"
          ? item.seriesId
          : null;
      if (
        !normalizedSeriesId ||
        !seriesList.some((series) => series.id === normalizedSeriesId) ||
        !item.volumeNumber
      ) {
        return;
      }
      const list = map.get(normalizedSeriesId) ?? [];
      list.push(item.volumeNumber);
      map.set(normalizedSeriesId, list);
    });
    map.forEach((list, key) => {
      const unique = Array.from(new Set(list)).sort((a, b) => a - b);
      map.set(key, unique);
    });
    return map;
  }, [userBooks, seriesList]);

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

    // シリーズに属する本を収集
    userBooks.forEach((item) => {
      const normalizedSeriesId =
        item.seriesId && item.seriesId !== "null" && item.seriesId !== "undefined"
          ? item.seriesId
          : null;
      const book = booksById.get(item.bookId);
      if (!book) {
        return;
      }
      const isSeriesBook =
        !!normalizedSeriesId && seriesById.has(normalizedSeriesId);
      if (!isSeriesBook) {
        return;
      }
      seriesBookIDs.add(item.bookId); // book.id ではなく item.bookId を使用
      const existing = grouped.get(normalizedSeriesId);
      if (existing) {
        existing.books.push(book);
      } else {
        grouped.set(normalizedSeriesId, {
          seriesId: normalizedSeriesId,
          name:
            seriesById.get(normalizedSeriesId)?.name ||
            book.seriesName ||
            "未判定",
          books: [book],
        });
      }
    });

    const seriesCards = Array.from(grouped.values()).sort((a, b) =>
      a.name.localeCompare(b.name, "ja")
    );

    // 単行本: ユーザーの蔵書でシリーズに属さない本
    const singleBooks = items.filter((book) => {
      if (!userBookIds.has(book.id)) return false;
      if (seriesBookIDs.has(book.id)) return false;
      return true;
    });

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
      {!isLoading && !error && seriesCards.length === 0 && singleBooks.length === 0 ? (
        <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 text-sm text-[#5c5d63]">
          まだ登録された書籍がありません。
        </div>
      ) : null}

      {/* シリーズ本の本棚 */}
      {!isLoading && seriesCards.length > 0 ? (
        <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
          <div className="mb-4 flex items-center gap-3">
            <span className="text-2xl">📚</span>
            <div>
              <h2 className="font-[var(--font-display)] text-xl text-[#1b1c1f]">
                シリーズ本棚
              </h2>
              <p className="text-xs text-[#5c5d63]">
                {seriesCards.length} シリーズ
              </p>
            </div>
          </div>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
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
                  className="group rounded-2xl border border-[#e4d8c7] bg-[#faf8f5] p-4 shadow-sm transition hover:-translate-y-1 hover:shadow-md"
                  href={`/books/series/${series.seriesId}`}
                >
                  <div className="flex items-start justify-between gap-3 text-xs text-[#5c5d63]">
                    <span className="rounded-full bg-[#c86b3c]/10 px-2 py-1 text-[#c86b3c]">
                      シリーズ
                    </span>
                    <button
                      aria-label={
                        favorite ? "シリーズのブックマーク解除" : "シリーズをブックマーク"
                      }
                      className={`flex h-7 w-7 items-center justify-center rounded-full border text-sm transition ${
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
                  <h3 className="mt-3 font-[var(--font-display)] text-lg text-[#1b1c1f] line-clamp-2">
                    {series.name}
                  </h3>
                  <p className="mt-1 text-sm text-[#5c5d63] line-clamp-1">
                    {authors.length > 0 ? authors.join(" / ") : "著者未登録"}
                  </p>
                  <div className="mt-3 flex flex-wrap items-center gap-2 text-xs text-[#5c5d63]">
                    <span className="rounded-full bg-[#f6f1e7] px-2 py-1">
                      {series.books.length} 冊
                    </span>
                    {volumeSummary ? (
                      <span className="rounded-full bg-[#f6f1e7] px-2 py-1">
                        {volumeSummary}
                      </span>
                    ) : null}
                  </div>
                  <p className="mt-3 text-xs text-[#c86b3c]">
                    詳細を見る →
                  </p>
                </Link>
              );
            })}
          </div>
        </section>
      ) : null}

      {/* 単行本の本棚 */}
      {!isLoading && singleBooks.length > 0 ? (
        <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
          <div className="mb-4 flex items-center gap-3">
            <span className="text-2xl">📖</span>
            <div>
              <h2 className="font-[var(--font-display)] text-xl text-[#1b1c1f]">
                単行本棚
              </h2>
              <p className="text-xs text-[#5c5d63]">
                {singleBooks.length} 冊
              </p>
            </div>
          </div>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {singleBooks.map((book, index) => (
              <Link
                key={book.id || book.isbn13 || `book-${index}`}
                className="group rounded-2xl border border-[#e4d8c7] bg-[#faf8f5] p-4 shadow-sm transition hover:-translate-y-1 hover:shadow-md"
                href={`/books/${book.id}`}
              >
                <div className="flex items-center justify-between text-xs text-[#5c5d63]">
                  <span className="rounded-full bg-[#5c5d63]/10 px-2 py-1">
                    単行本
                  </span>
                  <span>{book.publishedDate || ""}</span>
                </div>
                <h3 className="mt-3 font-[var(--font-display)] text-lg text-[#1b1c1f] line-clamp-2">
                  {book.title}
                </h3>
                <p className="mt-1 text-sm text-[#5c5d63] line-clamp-1">
                  {book.authors?.join(" / ") || "著者未登録"}
                </p>
                <div className="mt-3 flex flex-wrap items-center gap-2 text-xs text-[#5c5d63]">
                  <span className="rounded-full bg-[#f6f1e7] px-2 py-1">
                    {book.publisher || "出版社未登録"}
                  </span>
                  {book.isbn13 ? (
                    <span className="rounded-full bg-[#f6f1e7] px-2 py-1">
                      ISBN: {book.isbn13}
                    </span>
                  ) : null}
                </div>
                <p className="mt-3 text-xs text-[#c86b3c]">詳細を見る →</p>
              </Link>
            ))}
          </div>
        </section>
      ) : null}
    </div>
  );
}
