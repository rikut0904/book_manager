"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

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
  seriesSource: string;
};

type Series = {
  id: string;
  name: string;
};

export default function SeriesBookEditPage() {
  const params = useParams<{ seriesId: string; bookId: string }>();
  const bookId = params?.bookId;
  const [book, setBook] = useState<Book | null>(null);
  const [userBook, setUserBook] = useState<UserBook | null>(null);
  const [seriesList, setSeriesList] = useState<Series[]>([]);
  const [isSeries, setIsSeries] = useState(true);
  const [seriesId, setSeriesId] = useState(params?.seriesId ?? "");
  const [volumeNumber, setVolumeNumber] = useState("");
  const [note, setNote] = useState("");
  const [seriesMessage, setSeriesMessage] = useState<string | null>(null);
  const [noteMessage, setNoteMessage] = useState<string | null>(null);
  const [reportMessage, setReportMessage] = useState<string | null>(null);
  const [suggestion, setSuggestion] = useState("");
  const [reportNote, setReportNote] = useState("");

  const loadUserBook = useCallback(async () => {
    if (!bookId) {
      return;
    }
    const data = await fetchJSON<{ items: UserBook[] }>(
      `/user-books?bookId=${encodeURIComponent(bookId)}`,
      { auth: true }
    );
    setUserBook(data.items?.[0] ?? null);
  }, [bookId]);

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
        setBook(null);
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

  useEffect(() => {
    loadUserBook().catch(() => {
      setUserBook(null);
    });
  }, [loadUserBook]);

  useEffect(() => {
    if (!userBook) {
      setNote("");
      return;
    }
    setIsSeries(Boolean(userBook.seriesId));
    setSeriesId(userBook.seriesId || params?.seriesId || "");
    setVolumeNumber(userBook.volumeNumber ? String(userBook.volumeNumber) : "");
    setNote(userBook.note || "");
  }, [params?.seriesId, userBook]);

  const handleSeriesSave = async () => {
    if (!bookId) {
      return;
    }
    setSeriesMessage(null);
    try {
      if (!isSeries) {
        await fetchJSON("/user-series/override", {
          method: "PATCH",
          auth: true,
          body: JSON.stringify({
            bookId,
            isSeries: false,
          }),
        });
        setSeriesMessage("単行本として保存しました。");
        await loadUserBook();
        return;
      }
      if (!seriesId.trim() || !volumeNumber.trim()) {
        setSeriesMessage("seriesId と巻数を入力してください。");
        return;
      }
      await fetchJSON("/user-series/override", {
        method: "PATCH",
        auth: true,
        body: JSON.stringify({
          bookId,
          seriesId: seriesId.trim(),
          volumeNumber: Number(volumeNumber),
          isSeries: true,
        }),
      });
      setSeriesMessage("シリーズを保存しました。");
      await loadUserBook();
    } catch {
      setSeriesMessage("シリーズ設定に失敗しました。");
    }
  };

  const handleNoteSave = async () => {
    if (!bookId) {
      return;
    }
    setNoteMessage(null);
    try {
      if (userBook?.id) {
        await fetchJSON(`/user-books/${userBook.id}`, {
          method: "PATCH",
          auth: true,
          body: JSON.stringify({
            note,
          }),
        });
      } else {
        await fetchJSON("/user-books", {
          method: "POST",
          auth: true,
          body: JSON.stringify({
            bookId,
            note,
          }),
        });
      }
      setNoteMessage("所蔵メモを保存しました。");
      await loadUserBook();
    } catch {
      setNoteMessage("所蔵メモの保存に失敗しました。");
    }
  };

  const handleReportSubmit = async () => {
    if (!bookId) {
      return;
    }
    setReportMessage(null);
    if (!suggestion.trim()) {
      setReportMessage("修正内容を入力してください。");
      return;
    }
    try {
      await fetchJSON("/book-reports", {
        method: "POST",
        auth: true,
        body: JSON.stringify({
          bookId,
          suggestion: suggestion.trim(),
          note: reportNote.trim(),
        }),
      });
      setReportMessage("修正報告を送信しました。");
      setSuggestion("");
      setReportNote("");
    } catch {
      setReportMessage("修正報告の送信に失敗しました。");
    }
  };

  const seriesNameLabel =
    seriesList.find((item) => item.id === userBook?.seriesId)?.name ||
    book?.seriesName ||
    "未判定";
  const seriesSourceLabel =
    userBook?.seriesSource === "manual"
      ? "手動入力"
      : userBook?.seriesSource === "auto"
      ? "自動判定"
      : "未判定";

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

      <section className="grid gap-6 lg:grid-cols-[1.2fr_0.8fr]">
        <div className="flex flex-col gap-6">
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
            <h2 className="font-[var(--font-display)] text-2xl">
              現在の情報
            </h2>
            <div className="mt-4 grid gap-3 text-sm text-[#5c5d63]">
              <p>シリーズ: {seriesNameLabel}</p>
              <p>判定種別: {seriesSourceLabel}</p>
              <p>
                巻数:{" "}
                {userBook?.volumeNumber
                  ? `Vol.${userBook.volumeNumber}`
                  : "未判定"}
              </p>
              <p>所蔵メモ: {userBook?.note || "未登録"}</p>
            </div>
          </div>

          <details className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
            <summary className="flex cursor-pointer list-none items-center justify-between gap-3 font-[var(--font-display)] text-2xl">
              <span>シリーズ判定</span>
              <span className="rounded-full border border-[#e4d8c7] px-3 py-1 text-xs text-[#5c5d63]">
                編集する ▼
              </span>
            </summary>
            <p className="mt-2 text-sm text-[#5c5d63]">
              シリーズか単行本かを選択して保存します。
            </p>
            <label className="mt-4 flex items-center justify-between rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm">
              <span>シリーズ作品として登録</span>
              <input
                type="checkbox"
                checked={isSeries}
                onChange={(event) => {
                  const checked = event.target.checked;
                  setIsSeries(checked);
                  if (!checked) {
                    setSeriesId("");
                    setVolumeNumber("");
                  }
                }}
              />
            </label>
            {isSeries ? (
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
                  巻数
                  <input
                    className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
                    value={volumeNumber}
                    onChange={(event) => setVolumeNumber(event.target.value)}
                    placeholder="3"
                  />
                </label>
              </div>
            ) : null}
            {seriesMessage ? (
              <p className="mt-3 text-xs text-[#5c5d63]">{seriesMessage}</p>
            ) : null}
            <button
              className="mt-4 rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white"
              type="button"
              onClick={handleSeriesSave}
            >
              シリーズ判定を保存
            </button>
          </details>

          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
            <h2 className="font-[var(--font-display)] text-2xl">所蔵メモ</h2>
            <textarea
              className="mt-4 h-28 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              value={note}
              onChange={(event) => setNote(event.target.value)}
              placeholder="所蔵メモを入力"
            />
            {noteMessage ? (
              <p className="mt-3 text-xs text-[#5c5d63]">{noteMessage}</p>
            ) : null}
            <button
              className="mt-4 rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white"
              type="button"
              onClick={handleNoteSave}
            >
              所蔵メモを保存
            </button>
          </div>
        </div>

        <div className="flex flex-col gap-6">
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
            <h2 className="font-[var(--font-display)] text-2xl">
              書誌修正の連絡
            </h2>
            <p className="mt-2 text-sm text-[#5c5d63]">
              書誌情報に誤りがある場合はこちらから報告してください。
            </p>
            <label className="mt-4 block text-sm text-[#1b1c1f]">
              修正内容
              <textarea
                className="mt-2 h-24 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
                value={suggestion}
                onChange={(event) => setSuggestion(event.target.value)}
                placeholder="正しいタイトル・著者など"
              />
            </label>
            <label className="mt-4 block text-sm text-[#1b1c1f]">
              備考（任意）
              <textarea
                className="mt-2 h-20 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
                value={reportNote}
                onChange={(event) => setReportNote(event.target.value)}
                placeholder="補足情報"
              />
            </label>
            {reportMessage ? (
              <p className="mt-3 text-xs text-[#5c5d63]">{reportMessage}</p>
            ) : null}
            <button
              className="mt-4 rounded-full bg-[#c86b3c] px-5 py-3 text-sm font-medium text-white"
              type="button"
              onClick={handleReportSubmit}
            >
              報告を送信
            </button>
          </div>
        </div>
      </section>
    </div>
  );
}
