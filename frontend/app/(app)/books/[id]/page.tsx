"use client";

import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import { fetchJSON } from "@/lib/api";

type Book = {
  id: string;
  title: string;
  authors: string[];
  isbn13: string;
  publisher: string;
  publishedDate: string;
  thumbnailUrl: string;
  seriesName?: string;
};

type UserBook = {
  id: string;
  bookId: string;
  seriesId: string;
  volumeNumber: number;
  note: string;
};

type Favorite = {
  id: string;
  type: "book" | "series";
  bookId: string;
  seriesId: string;
};

export default function BookDetailPage() {
  const params = useParams<{ id: string }>();
  const [book, setBook] = useState<Book | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [userBook, setUserBook] = useState<UserBook | null>(null);
  const [deleteMessage, setDeleteMessage] = useState<string | null>(null);
  const [favorites, setFavorites] = useState<Favorite[]>([]);
  const router = useRouter();

  useEffect(() => {
    if (!params?.id) {
      return;
    }
    let isMounted = true;
    fetchJSON<Book>(`/books/${params.id}`, { auth: true })
      .then((data) => {
        if (!isMounted) {
          return;
        }
        setBook(data);
      })
      .catch(() => {
        if (!isMounted) {
          return;
        }
        setError("書籍詳細を取得できませんでした。");
      });
    return () => {
      isMounted = false;
    };
  }, [params?.id, router]);

  useEffect(() => {
    let isMounted = true;
    fetchJSON<{ items: Favorite[] }>("/favorites", { auth: true })
      .then((data) => {
        if (!isMounted) {
          return;
        }
        setFavorites(data.items ?? []);
      })
      .catch(() => {
        if (!isMounted) {
          return;
        }
        setFavorites([]);
      });
    return () => {
      isMounted = false;
    };
  }, []);

  useEffect(() => {
    if (!params?.id) {
      return;
    }
    let isMounted = true;
    fetchJSON<{ items: UserBook[] }>(
      `/user-books?bookId=${encodeURIComponent(params.id)}`,
      { auth: true }
    )
      .then((data) => {
        if (!isMounted) {
          return;
        }
        const item = data.items?.[0] ?? null;
        if (item?.seriesId) {
          router.replace(`/books/series/${item.seriesId}/${params.id}`);
          return;
        }
        setUserBook(item);
      })
      .catch(() => {
        if (!isMounted) {
          return;
        }
        setUserBook(null);
      });
    return () => {
      isMounted = false;
    };
  }, [params?.id, router]);

  const handleEdit = () => {
    if (!params?.id) {
      return;
    }
    router.push(`/books/${params.id}/edit`);
  };

  const handleDeleteBook = async () => {
    setDeleteMessage(null);
    if (!params?.id) {
      return;
    }
    if (!window.confirm("この書籍を削除しますか？")) {
      return;
    }
    try {
      await fetchJSON(`/books/${params.id}`, {
        method: "DELETE",
        auth: true,
      });
      router.push("/books");
    } catch {
      setDeleteMessage("書籍の削除に失敗しました。");
    }
  };

  const handleToggleFavorite = async () => {
    if (!params?.id) {
      return;
    }
    const existing = favorites.find(
      (item) => item.type === "book" && item.bookId === params.id
    );
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
            type: "book",
            bookId: params.id,
          }),
        });
        setFavorites((prev) => [...prev, created]);
      }
    } catch {
      setDeleteMessage("お気に入りの更新に失敗しました。");
    }
  };

  const displayVolume = userBook?.volumeNumber || 0;
  const favorite = favorites.find(
    (item) => item.type === "book" && item.bookId === params?.id
  );

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
              Book Detail
            </p>
            <div className="mt-2 flex items-baseline gap-3">
              <h1 className="font-[var(--font-display)] text-3xl">
                {book?.title || "書籍"}
              </h1>
              {displayVolume ? (
                <span className="text-sm text-[#5c5d63]">
                  Vol.{displayVolume}
                </span>
              ) : null}
            </div>
            <p className="mt-1 text-sm text-[#5c5d63]">
              {book?.authors?.join(" / ") || "著者未登録"}
            </p>
          </div>
          <div className="flex flex-wrap items-center gap-3">
            <button
              aria-label={favorite ? "お気に入り解除" : "お気に入り登録"}
              className={`flex h-10 w-10 items-center justify-center rounded-full border text-base transition ${
                favorite
                  ? "border-[#c86b3c] bg-[#c86b3c] text-white"
                  : "border-[#e4d8c7] text-[#5c5d63] hover:bg-white"
              }`}
              type="button"
              onClick={handleToggleFavorite}
            >
              {favorite ? "★" : "☆"}
            </button>
            <button
              className="rounded-full border border-[#e4d8c7] px-4 py-2 text-xs text-[#5c5d63] hover:bg-white"
              type="button"
              onClick={handleEdit}
            >
              編集
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
        {deleteMessage ? (
          <p className="mt-4 text-sm text-red-600">{deleteMessage}</p>
        ) : null}
        <div className="mt-6 grid gap-4 md:grid-cols-3">
          <div className="rounded-2xl border border-[#e4d8c7] bg-white/70 p-4 text-sm text-[#5c5d63]">
            <p className="text-xs text-[#c86b3c]">ISBN</p>
            <p className="mt-2 text-[#1b1c1f]">
              {book?.isbn13 || "未登録"}
            </p>
          </div>
          <div className="rounded-2xl border border-[#e4d8c7] bg-white/70 p-4 text-sm text-[#5c5d63]">
            <p className="text-xs text-[#c86b3c]">出版社</p>
            <p className="mt-2 text-[#1b1c1f]">
              {book?.publisher || "未登録"}
            </p>
          </div>
          <div className="rounded-2xl border border-[#e4d8c7] bg-white/70 p-4 text-sm text-[#5c5d63]">
            <p className="text-xs text-[#c86b3c]">所蔵メモ</p>
            <p className="mt-2 text-[#1b1c1f]">
              {userBook?.note || ""}
            </p>
          </div>
        </div>
        <div className="mt-6 flex justify-end">
          <button
            className="rounded-full border border-[#e4d8c7] px-4 py-2 text-xs text-[#c86b3c] hover:bg-white"
            type="button"
            onClick={handleDeleteBook}
          >
            書籍を削除
          </button>
        </div>
      </section>

    </div>
  );
}
