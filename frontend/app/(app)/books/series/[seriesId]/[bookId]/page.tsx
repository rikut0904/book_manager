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

type Series = {
  id: string;
  name: string;
};

type UserBook = {
  id: string;
  bookId: string;
  seriesId: string;
  volumeNumber: number;
};

type Favorite = {
  id: string;
  type: "book" | "series";
  bookId: string;
  seriesId: string;
};

export default function SeriesBookDetailPage() {
  const params = useParams<{ seriesId: string; bookId: string }>();
  const bookId = params?.bookId;
  const [book, setBook] = useState<Book | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [seriesId, setSeriesId] = useState("");
  const [seriesList, setSeriesList] = useState<Series[]>([]);
  const [volumeNumber, setVolumeNumber] = useState("");
  const [seriesMessage, setSeriesMessage] = useState<string | null>(null);
  const [userBook, setUserBook] = useState<UserBook | null>(null);
  const [deleteMessage, setDeleteMessage] = useState<string | null>(null);
  const [favorites, setFavorites] = useState<Favorite[]>([]);
  const manualSeriesName =
    seriesList.find((item) => item.id === seriesId)?.name || seriesId;
  const storedSeriesName = userBook?.seriesId
    ? seriesList.find((item) => item.id === userBook.seriesId)?.name ||
      userBook.seriesId
    : book?.seriesName || "未判定";
  const router = useRouter();

  useEffect(() => {
    if (!bookId) {
      return;
    }
    let isMounted = true;
    fetchJSON<Book>(`/books/${bookId}`, { auth: true })
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
  }, [bookId]);

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
    if (!bookId) {
      return;
    }
    let isMounted = true;
    fetchJSON<{ items: UserBook[] }>(
      `/user-books?bookId=${encodeURIComponent(bookId)}`,
      { auth: true }
    )
      .then((data) => {
        if (!isMounted) {
          return;
        }
        setUserBook(data.items?.[0] ?? null);
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
  }, [bookId]);

  useEffect(() => {
    let isMounted = true;
    fetchJSON<{ items: Series[] }>("/series", { auth: true })
      .then((data) => {
        if (!isMounted) {
          return;
        }
        setSeriesList(data.items ?? []);
      })
      .catch(() => {
        if (!isMounted) {
          return;
        }
        setSeriesList([]);
      });
    return () => {
      isMounted = false;
    };
  }, []);

  const handleSeriesOverride = async () => {
    setSeriesMessage(null);
    if (!bookId) {
      return;
    }
    if (!seriesId.trim() || !volumeNumber.trim()) {
      setSeriesMessage("seriesId と巻数を入力してください。");
      return;
    }
    try {
      await fetchJSON("/user-series/override", {
        method: "PATCH",
        auth: true,
        body: JSON.stringify({
          bookId,
          seriesId: seriesId.trim(),
          volumeNumber: Number(volumeNumber),
        }),
      });
      setSeriesMessage("シリーズを上書きしました。");
    } catch {
      setSeriesMessage("シリーズ上書きに失敗しました。");
    }
  };

  const handleDeleteBook = async () => {
    setDeleteMessage(null);
    if (!bookId) {
      return;
    }
    if (!window.confirm("この書籍を削除しますか？")) {
      return;
    }
    try {
      await fetchJSON(`/books/${bookId}`, {
        method: "DELETE",
        auth: true,
      });
      router.push("/books");
    } catch {
      setDeleteMessage("書籍の削除に失敗しました。");
    }
  };

  const handleToggleFavorite = async () => {
    if (!bookId) {
      return;
    }
    const existing = favorites.find(
      (item) => item.type === "book" && item.bookId === bookId
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
            bookId,
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
    (item) => item.type === "book" && item.bookId === bookId
  );
  const backToSeries = params?.seriesId
    ? `/books/series/${params.seriesId}`
    : "/books";

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
            <Link
              className="rounded-full border border-[#e4d8c7] px-4 py-2 text-xs text-[#5c5d63]"
              href={backToSeries}
            >
              シリーズへ戻る
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
            <p className="mt-2 text-[#1b1c1f]">サイン本</p>
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

      <section className="grid gap-6">
        <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
          <h2 className="font-[var(--font-display)] text-2xl">
            シリーズ判定
          </h2>
          <p className="mt-2 text-sm text-[#5c5d63]">
            保存済みのシリーズ情報を確認し、必要なら上書きしてください。
          </p>
          <div className="mt-4 rounded-2xl border border-[#e4d8c7] bg-white p-4 text-sm text-[#5c5d63]">
            <div className="flex items-center justify-between">
              <span className="text-[#1b1c1f]">
                {seriesId ? manualSeriesName || "未判定" : storedSeriesName}
              </span>
              <span>
                {seriesId && volumeNumber
                  ? `Vol.${volumeNumber}`
                  : userBook?.volumeNumber
                  ? `Vol.${userBook.volumeNumber}`
                  : "--"}
              </span>
            </div>
            <p className="mt-2 text-xs text-[#c86b3c]">
              {seriesId
                ? "手動入力"
                : userBook?.seriesId
                ? "保存済み"
                : "未判定"}
            </p>
          </div>
          <div className="mt-4 grid gap-3 text-sm">
            <label className="text-[#1b1c1f]">
              シリーズ名
              <select
                className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
                value={seriesId}
                onChange={(event) => setSeriesId(event.target.value)}
              >
                <option value="">シリーズを選択</option>
                {seriesList.map((item) => (
                  <option key={item.id} value={item.id}>
                    {item.name}
                  </option>
                ))}
              </select>
            </label>
            <label className="text-[#1b1c1f]">
              seriesId（直接入力）
              <input
                className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
                value={seriesId}
                onChange={(event) => setSeriesId(event.target.value)}
                placeholder="series_123"
              />
            </label>
            <label className="text-[#1b1c1f]">
              巻数
              <input
                className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
                value={volumeNumber}
                onChange={(event) => setVolumeNumber(event.target.value)}
                placeholder="3"
              />
            </label>
            {seriesMessage ? (
              <p className="text-xs text-[#5c5d63]">{seriesMessage}</p>
            ) : null}
            <button
              className="rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white"
              type="button"
              onClick={handleSeriesOverride}
            >
              上書き保存
            </button>
          </div>
        </div>

      </section>
    </div>
  );
}
