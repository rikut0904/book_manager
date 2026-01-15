"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { useEffect, useState } from "react";

import BookEditSections, {
  type Book,
  type Series,
  type UserBook,
} from "@/app/(app)/books/_components/BookEditSections";
import { fetchJSON } from "@/lib/api";

export default function SeriesBookEditPage() {
  const params = useParams<{ seriesId: string; bookId: string }>();
  const bookId = params?.bookId;
  const [book, setBook] = useState<Book | null>(null);
  const [userBook, setUserBook] = useState<UserBook | null>(null);
  const [seriesList, setSeriesList] = useState<Series[]>([]);
  const seriesIdParam = params?.seriesId ?? "";

  const refreshUserBook = async () => {
    if (!bookId) {
      return;
    }
    const data = await fetchJSON<{ items: UserBook[] }>(
      `/user-books?bookId=${encodeURIComponent(bookId)}`,
      { auth: true }
    );
    setUserBook(data.items?.[0] ?? null);
  };

  useEffect(() => {
    if (!bookId) {
      return;
    }
    let isMounted = true;
    const load = async () => {
      const [bookRes, seriesRes, userBooksRes] = await Promise.all([
        fetchJSON<Book>(`/books/${bookId}`, { auth: true }).catch(() => null),
        fetchJSON<{ items: Series[] }>("/series", { auth: true }).catch(() => ({
          items: [],
        })),
        fetchJSON<{ items: UserBook[] }>(
          `/user-books?bookId=${encodeURIComponent(bookId)}`,
          { auth: true }
        ).catch(() => ({ items: [] })),
      ]);
      if (!isMounted) {
        return;
      }
      setBook(bookRes);
      setSeriesList(seriesRes.items ?? []);
      setUserBook(userBooksRes.items?.[0] ?? null);
    };
    load();
    return () => {
      isMounted = false;
    };
  }, [bookId]);

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
              Book Edit
            </p>
            <h1 className="mt-2 font-[var(--font-display)] text-3xl">
              {book?.title || "書籍編集"}
            </h1>
            <p className="mt-1 text-sm text-[#5c5d63]">
              {book?.authors?.join(" / ") || "著者未登録"}
            </p>
          </div>
          <Link
            className="rounded-full border border-[#e4d8c7] px-4 py-2 text-xs text-[#5c5d63]"
            href={`/books/series/${params?.seriesId ?? ""}/${bookId ?? ""}`}
          >
            詳細へ戻る
          </Link>
        </div>
      </section>

      {bookId ? (
        <BookEditSections
          key={`${userBook?.id ?? "none"}-${seriesIdParam}`}
          bookId={bookId}
          book={book}
          userBook={userBook}
          seriesList={seriesList}
          onRefreshUserBook={refreshUserBook}
          initialSeriesId={seriesIdParam}
          seriesSectionTitle="シリーズ判定"
        />
      ) : null}
    </div>
  );
}
