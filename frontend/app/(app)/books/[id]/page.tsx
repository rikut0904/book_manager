"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
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
};

type Tag = {
  id: string;
  name: string;
};

export default function BookDetailPage() {
  const params = useParams<{ id: string }>();
  const [book, setBook] = useState<Book | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [tags, setTags] = useState<Tag[]>([]);
  const [tagId, setTagId] = useState("");
  const [tagMessage, setTagMessage] = useState<string | null>(null);
  const [seriesId, setSeriesId] = useState("");
  const [volumeNumber, setVolumeNumber] = useState("");
  const [seriesMessage, setSeriesMessage] = useState<string | null>(null);

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
  }, [params?.id]);

  useEffect(() => {
    let isMounted = true;
    fetchJSON<{ items: Tag[] }>("/tags", { auth: true })
      .then((data) => {
        if (!isMounted) {
          return;
        }
        setTags(data.items ?? []);
      })
      .catch(() => {
        if (!isMounted) {
          return;
        }
        setTags([]);
      });
    return () => {
      isMounted = false;
    };
  }, []);

  const handleSeriesOverride = async () => {
    setSeriesMessage(null);
    if (!params?.id) {
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
          bookId: params.id,
          seriesId: seriesId.trim(),
          volumeNumber: Number(volumeNumber),
        }),
      });
      setSeriesMessage("シリーズを上書きしました。");
    } catch {
      setSeriesMessage("シリーズ上書きに失敗しました。");
    }
  };

  const handleAddTag = async () => {
    setTagMessage(null);
    if (!params?.id) {
      return;
    }
    if (!tagId.trim()) {
      setTagMessage("tagId を入力してください。");
      return;
    }
    try {
      await fetchJSON("/book-tags", {
        method: "POST",
        auth: true,
        body: JSON.stringify({
          bookId: params.id,
          tagId: tagId.trim(),
        }),
      });
      setTagMessage("タグを付与しました。");
    } catch {
      setTagMessage("タグ付与に失敗しました。");
    }
  };

  const handleRemoveTag = async () => {
    setTagMessage(null);
    if (!params?.id) {
      return;
    }
    if (!tagId.trim()) {
      setTagMessage("tagId を入力してください。");
      return;
    }
    try {
      await fetchJSON("/book-tags", {
        method: "DELETE",
        auth: true,
        body: JSON.stringify({
          bookId: params.id,
          tagId: tagId.trim(),
        }),
      });
      setTagMessage("タグを解除しました。");
    } catch {
      setTagMessage("タグ解除に失敗しました。");
    }
  };

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
              Book Detail
            </p>
            <h1 className="mt-2 font-[var(--font-display)] text-3xl">
              {book?.title || "書籍"}
            </h1>
            <p className="mt-1 text-sm text-[#5c5d63]">
              {book?.authors?.join(" / ") || "著者未登録"}
            </p>
          </div>
          <Link
            className="rounded-full border border-[#e4d8c7] px-4 py-2 text-xs text-[#5c5d63]"
            href="/books"
          >
            一覧へ戻る
          </Link>
        </div>
        {error ? (
          <p className="mt-4 text-sm text-red-600">{error}</p>
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
      </section>

      <section className="grid gap-6 lg:grid-cols-[1.2fr_0.8fr]">
        <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
          <h2 className="font-[var(--font-display)] text-2xl">
            シリーズ判定
          </h2>
          <p className="mt-2 text-sm text-[#5c5d63]">
            自動判定されたシリーズ情報を確認し、必要なら上書きしてください。
          </p>
          <div className="mt-4 rounded-2xl border border-[#e4d8c7] bg-white p-4 text-sm text-[#5c5d63]">
            <div className="flex items-center justify-between">
              <span className="text-[#1b1c1f]">星街メトロ</span>
              <span>Vol.4</span>
            </div>
            <p className="mt-2 text-xs text-[#c86b3c]">自動判定: 信頼度 86%</p>
          </div>
          <div className="mt-4 grid gap-3 text-sm">
            <label className="text-[#1b1c1f]">
              シリーズ名
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

        <div className="flex flex-col gap-6">
          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
            <h2 className="font-[var(--font-display)] text-2xl">タグ</h2>
            <div className="mt-4 flex flex-wrap gap-2 text-xs">
              {tags.length === 0 ? (
                <span className="rounded-full border border-[#e4d8c7] bg-white px-3 py-2 text-[#5c5d63]">
                  タグがありません
                </span>
              ) : null}
              {tags.map((tag) => (
                <span
                  key={tag.id}
                  className="rounded-full border border-[#e4d8c7] bg-white px-3 py-2 text-[#5c5d63]"
                >
                  {tag.name}
                </span>
              ))}
            </div>
            <div className="mt-4 flex flex-wrap gap-2">
              <select
                className="min-w-[200px] rounded-full border border-[#e4d8c7] bg-white px-4 py-2 text-xs outline-none transition focus:border-[#c86b3c]"
                value={tagId}
                onChange={(event) => setTagId(event.target.value)}
              >
                <option value="">タグを選択</option>
                {tags.map((tag) => (
                  <option key={tag.id} value={tag.id}>
                    {tag.name}
                  </option>
                ))}
              </select>
              <input
                className="min-w-[200px] rounded-full border border-[#e4d8c7] bg-white px-4 py-2 text-xs outline-none transition focus:border-[#c86b3c]"
                placeholder="tagId を直接入力"
                value={tagId}
                onChange={(event) => setTagId(event.target.value)}
              />
              <button
                className="rounded-full bg-[#efe5d4] px-3 py-2 text-xs text-[#1b1c1f]"
                type="button"
                onClick={handleAddTag}
              >
                付与
              </button>
              <button
                className="rounded-full border border-[#e4d8c7] px-3 py-2 text-xs text-[#5c5d63]"
                type="button"
                onClick={handleRemoveTag}
              >
                解除
              </button>
            </div>
            {tagMessage ? (
              <p className="mt-3 text-xs text-[#5c5d63]">{tagMessage}</p>
            ) : null}
          </div>

          <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
            <h2 className="font-[var(--font-display)] text-2xl">お気に入り</h2>
            <p className="mt-2 text-sm text-[#5c5d63]">
              この巻またはシリーズをお気に入りに登録できます。
            </p>
            <div className="mt-4 flex gap-3">
              <button className="rounded-full border border-[#e4d8c7] px-4 py-2 text-xs text-[#5c5d63]">
                単巻で登録
              </button>
              <button className="rounded-full bg-[#c86b3c] px-4 py-2 text-xs text-white">
                シリーズで登録
              </button>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}
