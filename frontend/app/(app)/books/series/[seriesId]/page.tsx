"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
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

export default function SeriesDetailPage() {
  const params = useParams<{ seriesId: string }>();
  const [books, setBooks] = useState<Book[]>([]);
  const [userBooks, setUserBooks] = useState<UserBook[]>([]);
  const [seriesList, setSeriesList] = useState<Series[]>([]);
  const [favorites, setFavorites] = useState<Favorite[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

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
        setBooks(booksRes.items ?? []);
        setUserBooks(userBooksRes.items ?? []);
        setSeriesList(seriesRes.items ?? []);
        setFavorites(favoritesRes.items ?? []);
      } catch {
        if (!isMounted) {
          return;
        }
        setError("シリーズ詳細を取得できませんでした。");
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

  const favoriteBySeriesId = useMemo(() => {
    const map = new Map<string, Favorite>();
    favorites
      .filter((item) => item.type === "series" && item.seriesId)
      .forEach((item) => {
        map.set(item.seriesId, item);
      });
    return map;
  }, [favorites]);

  const { seriesName, seriesBooks } = useMemo(() => {
    const seriesId = params?.seriesId || "";
    const seriesName =
      seriesList.find((item) => item.id === seriesId)?.name || "未判定";
    const booksById = new Map(books.map((book) => [book.id, book]));
    const seriesBooks = userBooks
      .filter((item) => item.seriesId === seriesId)
      .map((item) => ({
        ...item,
        book: booksById.get(item.bookId),
      }))
      .filter((item) => item.book)
      .sort((a, b) => {
        const aVolume = a.volumeNumber || Number.MAX_SAFE_INTEGER;
        const bVolume = b.volumeNumber || Number.MAX_SAFE_INTEGER;
        if (aVolume !== bVolume) {
          return aVolume - bVolume;
        }
        return (a.book?.title || "").localeCompare(b.book?.title || "", "ja");
      });
    return { seriesName, seriesBooks };
  }, [books, params?.seriesId, seriesList, userBooks]);

  const authors = useMemo(() => {
    const items = seriesBooks.flatMap((item) => item.book?.authors || []);
    return Array.from(new Set(items.filter(Boolean)));
  }, [seriesBooks]);

  const handleToggleFavorite = async (
    event: MouseEvent<HTMLButtonElement>
  ) => {
    event.preventDefault();
    event.stopPropagation();
    const seriesId = params?.seriesId;
    if (!seriesId) {
      return;
    }
    const existing = favoriteBySeriesId.get(seriesId);
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
      setError("お気に入りの更新に失敗しました。");
    }
  };

  const favorite = params?.seriesId
    ? favoriteBySeriesId.get(params.seriesId)
    : undefined;

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
              Series
            </p>
            <h1 className="mt-2 font-[var(--font-display)] text-3xl">
              {seriesName}
            </h1>
            <p className="mt-2 text-sm text-[#5c5d63]">
              {authors.length > 0 ? authors.join(" / ") : "著者未登録"}
            </p>
            <p className="mt-1 text-sm text-[#5c5d63]">
              巻数合計: {seriesBooks.length}
            </p>
          </div>
          <div className="flex flex-wrap gap-3">
            <button
              className="rounded-full border border-[#e4d8c7] px-4 py-2 text-xs text-[#5c5d63] hover:bg-white"
              type="button"
              onClick={handleToggleFavorite}
            >
              {favorite ? "お気に入り解除" : "お気に入り登録"}
            </button>
            <Link
              className="rounded-full border border-[#e4d8c7] px-4 py-2 text-xs text-[#5c5d63]"
              href="/books"
            >
              一覧へ戻る
            </Link>
          </div>
        </div>
        {error ? (
          <p className="mt-4 text-sm text-red-600">{error}</p>
        ) : null}
      </section>

      <section className="grid gap-4 lg:grid-cols-3">
        {isLoading ? (
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 text-sm text-[#5c5d63]">
            読み込み中...
          </div>
        ) : null}
        {!isLoading && seriesBooks.length === 0 ? (
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 text-sm text-[#5c5d63]">
            シリーズの書籍が見つかりませんでした。
          </div>
        ) : null}
        {seriesBooks.map((item) => {
          const book = item.book;
          if (!book) {
            return null;
          }
          return (
            <Link
              key={book.id}
              className="group rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 shadow-sm transition hover:-translate-y-1 hover:shadow-md"
              href={`/books/series/${item.seriesId}/${book.id}`}
            >
              <div className="flex items-center justify-between text-xs text-[#5c5d63]">
                <span>{book.publisher || "出版社未登録"}</span>
                <span>
                  {item.volumeNumber ? `Vol.${item.volumeNumber}` : "Vol.--"}
                </span>
              </div>
              <h2 className="mt-4 font-[var(--font-display)] text-xl text-[#1b1c1f]">
                {book.title}
              </h2>
              <p className="mt-2 text-sm text-[#5c5d63]">
                {book.authors?.join(" / ") || "著者未登録"}
              </p>
              <div className="mt-4 rounded-2xl bg-[#f6f1e7] px-3 py-2 text-xs text-[#5c5d63]">
                発売日: {book.publishedDate || "未登録"}
              </div>
              <p className="mt-4 text-xs text-[#c86b3c]">詳細を見る →</p>
            </Link>
          );
        })}
      </section>
    </div>
  );
}
