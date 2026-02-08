"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { fetchJSON } from "@/lib/api";
import { getAuthState } from "@/lib/auth";
import { commonErrorMessages } from "@/lib/errorMessages";
import { useUserProfile } from "@/lib/userProfile";

type UserResponse = {
  user: { id: string; email: string; userId: string; displayName: string };
  isAdmin?: boolean;
};

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

type Recommendation = {
  id: string;
  userId: string;
  bookId: string;
  comment: string;
  createdAt: string;
};

type DashboardResponse = {
  favorites: Favorite[];
  recommendations: Recommendation[];
  books: Book[];
  series: Series[];
  userBooks: UserBook[];
};

const errorMessages = commonErrorMessages;

export default function UserPage() {
  const { profile, error: profileError, refresh } = useUserProfile();
  const [bookmarks, setBookmarks] = useState<Favorite[]>([]);
  const [bookmarkError, setBookmarkError] = useState<string | null>(null);
  const [recs, setRecs] = useState<Recommendation[]>([]);
  const [recError, setRecError] = useState<string | null>(null);
  const [books, setBooks] = useState<Book[]>([]);
  const [series, setSeries] = useState<Series[]>([]);
  const [userBooks, setUserBooks] = useState<UserBook[]>([]);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [editForm, setEditForm] = useState({ displayName: "", email: "" });
  const [editError, setEditError] = useState<string | null>(null);
  const [editMessage, setEditMessage] = useState<string | null>(null);
  const initialUserId = getAuthState()?.userId ?? "";

  useEffect(() => {
    if (!profile) {
      return;
    }
    setEditForm({
      displayName: profile.displayName ?? "",
      email: profile.email ?? "",
    });
  }, [profile]);

  useEffect(() => {
    const auth = getAuthState();
    if (!auth?.userId) {
      return;
    }
    let isMounted = true;
    fetchJSON<DashboardResponse>("/user/dashboard", { auth: true })
      .then((data) => {
        if (!isMounted) {
          return;
        }
        setBookmarks(data.favorites ?? []);
        setRecs(data.recommendations ?? []);
        setBooks(data.books ?? []);
        setSeries(data.series ?? []);
        setUserBooks(data.userBooks ?? []);
      })
      .catch(() => {
        if (!isMounted) {
          return;
        }
        setBookmarkError("ブックマークを取得できませんでした。");
        setRecError("おすすめした本を取得できませんでした。");
        setBooks([]);
        setSeries([]);
        setUserBooks([]);
      });
    return () => {
      isMounted = false;
    };
  }, []);

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

  const handleOpenEdit = () => {
    if (profile) {
      setEditForm({
        displayName: profile.displayName ?? "",
        email: profile.email ?? "",
      });
    }
    setEditError(null);
    setEditMessage(null);
    setIsEditOpen(true);
  };

  const handleSaveProfile = async () => {
    setEditError(null);
    setEditMessage(null);
    try {
      const data = await fetchJSON<{ user: UserResponse["user"] }>("/users/me", {
        method: "PATCH",
        auth: true,
        body: JSON.stringify({
          displayName: editForm.displayName,
          email: editForm.email,
        }),
      });
      setEditForm({
        displayName: data.user?.displayName ?? "",
        email: data.user?.email ?? "",
      });
      await refresh();
      setEditMessage("更新しました。");
    } catch (err) {
      const message = err instanceof Error ? err.message.trim() : "";
      setEditError(errorMessages[message] || "更新に失敗しました。");
    }
  };

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <div className="flex items-start justify-between gap-4">
          <div>
            <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
              User
            </p>
            <h1 className="mt-2 font-[var(--font-display)] text-3xl">
              @{profile?.userId || initialUserId || "読み込み中..."}
            </h1>
          </div>
          <Link
            className="flex h-10 w-10 items-center justify-center rounded-full border border-[#e4d8c7] text-[#5c5d63] hover:bg-white"
            href="/settings"
            aria-label="設定"
          >
            <svg
              aria-hidden="true"
              className="h-5 w-5"
              viewBox="0 0 24 24"
              fill="currentColor"
            >
              <path d="M19.14 12.94c.04-.31.06-.63.06-.94s-.02-.63-.06-.94l2.03-1.58a.5.5 0 0 0 .12-.64l-1.92-3.32a.5.5 0 0 0-.6-.22l-2.39.96a7.03 7.03 0 0 0-1.63-.94l-.36-2.54A.5.5 0 0 0 13.9 1h-3.8a.5.5 0 0 0-.49.41l-.36 2.54c-.59.23-1.13.53-1.63.94l-2.39-.96a.5.5 0 0 0-.6.22L2.7 7.96a.5.5 0 0 0 .12.64l2.03 1.58c-.04.31-.06.63-.06.94s.02.63.06.94l-2.03 1.58a.5.5 0 0 0-.12.64l1.92 3.32c.13.23.4.32.64.22l2.39-.96c.5.41 1.04.75 1.63.98l.36 2.54c.05.24.25.41.49.41h3.8c.24 0 .45-.17.49-.41l.36-2.54c.59-.23 1.13-.56 1.63-.98l2.39.96c.24.1.51.01.64-.22l1.92-3.32a.5.5 0 0 0-.12-.64l-2.03-1.58zM12 15.5A3.5 3.5 0 1 1 12 8a3.5 3.5 0 0 1 0 7.5z" />
            </svg>
          </Link>
        </div>
      </section>

      <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
        <div className="flex items-center justify-between gap-4">
          <h2 className="font-[var(--font-display)] text-2xl">基本情報</h2>
          <button
            className="inline-flex items-center gap-2 text-sm text-[#c86b3c] hover:text-[#8f3d1f]"
            type="button"
            onClick={handleOpenEdit}
          >
            <svg
              aria-hidden="true"
              className="h-4 w-4"
              viewBox="0 0 24 24"
              fill="currentColor"
            >
              <path d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zm14.71-9.04a1 1 0 0 0 0-1.41l-2.5-2.5a1 1 0 0 0-1.41 0l-1.88 1.88 3.75 3.75 2.04-2.72z" />
            </svg>
            編集
          </button>
        </div>
        {profileError ? (
          <p className="mt-3 text-sm text-red-600">{profileError}</p>
        ) : null}
        {!profileError && !profile ? (
          <p className="mt-3 text-sm text-[#5c5d63]">
            情報を取得しています...
          </p>
        ) : null}
        {profile ? (
          <div className="mt-4 grid gap-3 text-sm md:grid-cols-3">
            <div className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-4">
              <p className="text-xs text-[#5c5d63]">ユーザー名</p>
              <p className="mt-1 text-[#1b1c1f]">
                {profile.displayName || "未設定"}
              </p>
            </div>
            <div className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-4">
              <p className="text-xs text-[#5c5d63]">ユーザーID</p>
              <p className="mt-1 text-[#1b1c1f]">@{profile.userId}</p>
            </div>
            <div className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-4">
              <p className="text-xs text-[#5c5d63]">メールアドレス</p>
              <p className="mt-1 text-[#1b1c1f]">{profile.email}</p>
            </div>
          </div>
        ) : null}
      </section>

      <section className="grid gap-4 lg:grid-cols-2">
        <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
          <div className="flex items-start justify-between gap-3">
            <h2 className="font-[var(--font-display)] text-2xl">
              ブックマーク
            </h2>
            <Link
              className="text-xs text-[#c86b3c] hover:text-[#8f3d1f]"
              href="/user/bookmark"
            >
              すべて見る
            </Link>
          </div>
          {bookmarkError ? (
            <p className="mt-3 text-sm text-red-600">{bookmarkError}</p>
          ) : null}
          {!bookmarkError && bookmarks.length === 0 ? (
            <p className="mt-3 text-sm text-[#5c5d63]">
              まだブックマークがありません。
            </p>
          ) : null}
          <div className="mt-4 grid gap-3">
            {bookmarks.slice(0, 2).map((item) => (
              <div
                key={item.id}
                className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm"
              >
                <p className="text-xs text-[#5c5d63]">
                  {item.type === "series" ? "シリーズ" : "単巻"}
                </p>
                <p className="mt-1 text-[#1b1c1f]">
                  {item.type === "series"
                    ? getSeriesName(item.seriesId)
                    : getBookTitleWithVolume(item.bookId)}
                </p>
              </div>
            ))}
          </div>
        </div>
        <div className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
          <div className="flex items-start justify-between gap-3">
            <h2 className="font-[var(--font-display)] text-2xl">おすすめした本</h2>
            <Link
              className="text-xs text-[#c86b3c] hover:text-[#8f3d1f]"
              href="/user/suggest"
            >
              すべて見る
            </Link>
          </div>
          {recError ? (
            <p className="mt-3 text-sm text-red-600">{recError}</p>
          ) : null}
          {!recError && recs.length === 0 ? (
            <p className="mt-3 text-sm text-[#5c5d63]">
              まだおすすめした本がありません。
            </p>
          ) : null}
          <div className="mt-4 grid gap-3">
            {recs.slice(0, 2).map((item) => (
              <div
                key={item.id}
                className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm"
              >
                <p className="text-xs text-[#5c5d63]">おすすめした本</p>
                <p className="mt-1 text-[#1b1c1f]">
                  {getBookTitleWithVolume(item.bookId)}
                </p>
                <p className="mt-2 text-xs text-[#5c5d63]">コメント: {item.comment}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {isEditOpen ? (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4">
          <div className="w-full max-w-lg rounded-3xl border border-[#e4d8c7] bg-white p-6 shadow-xl">
            <div className="flex items-start justify-between gap-4">
              <div>
                <h2 className="font-[var(--font-display)] text-2xl">
                  基本情報の編集
                </h2>
                <p className="mt-1 text-sm text-[#5c5d63]">
                  ユーザー名とメールアドレスを更新できます。
                </p>
              </div>
              <button
                className="text-sm text-[#5c5d63] hover:text-[#1b1c1f]"
                type="button"
                onClick={() => setIsEditOpen(false)}
              >
                閉じる
              </button>
            </div>
            <div className="mt-4 grid gap-3 text-sm">
              <label className="flex flex-col gap-2">
                ユーザーID
                <input
                  className="rounded-2xl border border-[#e4d8c7] bg-[#f6f1e7] px-4 py-3 text-sm text-[#5c5d63]"
                  value={profile?.userId ?? ""}
                  disabled
                />
              </label>
              <label className="flex flex-col gap-2">
                ユーザー名
                <input
                  className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
                  value={editForm.displayName}
                  onChange={(event) =>
                    setEditForm((prev) => ({
                      ...prev,
                      displayName: event.target.value,
                    }))
                  }
                />
              </label>
              <label className="flex flex-col gap-2">
                メールアドレス
                <input
                  className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
                  type="email"
                  value={editForm.email}
                  onChange={(event) =>
                    setEditForm((prev) => ({
                      ...prev,
                      email: event.target.value,
                    }))
                  }
                />
              </label>
            </div>
            {editError ? (
              <p className="mt-3 text-xs text-red-600">{editError}</p>
            ) : null}
            {editMessage ? (
              <p className="mt-3 text-xs text-[#5c5d63]">{editMessage}</p>
            ) : null}
            <div className="mt-4 flex justify-end gap-2 text-sm">
              <button
                className="rounded-full border border-[#e4d8c7] px-4 py-2 text-[#5c5d63] hover:bg-[#f6f1e7]"
                type="button"
                onClick={() => setIsEditOpen(false)}
              >
                キャンセル
              </button>
              <button
                className="rounded-full bg-[#1b1c1f] px-5 py-2 text-white"
                type="button"
                onClick={handleSaveProfile}
              >
                保存する
              </button>
            </div>
          </div>
        </div>
      ) : null}
    </div>
  );
}
