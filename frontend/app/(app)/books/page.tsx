"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

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
};

export default function BooksPage() {
  const [items, setItems] = useState<Book[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let isMounted = true;
    fetchJSON<{ items: Book[] }>("/books", { auth: true })
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
        setError("書籍一覧を取得できませんでした。");
      })
      .finally(() => {
        if (!isMounted) {
          return;
        }
        setIsLoading(false);
      });
    return () => {
      isMounted = false;
    };
  }, []);

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
              Library
            </p>
            <h1 className="mt-2 font-[var(--font-display)] text-3xl">
              所蔵一覧
            </h1>
          </div>
          <div className="flex flex-wrap gap-3">
            <button className="rounded-full border border-[#e4d8c7] px-4 py-2 text-xs text-[#5c5d63]">
              タグ
            </button>
            <button className="rounded-full border border-[#e4d8c7] px-4 py-2 text-xs text-[#5c5d63]">
              シリーズ
            </button>
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
            <option>お気に入り</option>
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
        {items.map((book, index) => (
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
