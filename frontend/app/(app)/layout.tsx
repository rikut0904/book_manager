"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import { fetchJSON } from "@/lib/api";
import { clearAuthState, getAuthState } from "@/lib/auth";

const navItems = [
  { href: "/books", label: "ホーム" },
  { href: "/next-to-buy", label: "次に買う本" },
  { href: "/users", label: "ユーザー検索" },
  { href: "/user", label: "ユーザー" },
];

export default function AppLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [isReady, setIsReady] = useState(false);
  const [isAuthed, setIsAuthed] = useState(false);
  const [accountId, setAccountId] = useState("");
  const [userId, setUserId] = useState("");
  const [displayName, setDisplayName] = useState("");
  const router = useRouter();

  useEffect(() => {
    const auth = getAuthState();
    setIsAuthed(Boolean(auth?.accessToken));
    setAccountId(auth?.userId ?? "");
    setIsReady(true);
  }, []);

  useEffect(() => {
    if (!accountId) {
      return;
    }
    let isMounted = true;
    fetchJSON<{ user: { userId: string; displayName: string } }>(`/users/${accountId}`, {
      auth: true,
    })
      .then((data) => {
        if (!isMounted) {
          return;
        }
        setUserId(data.user?.userId ?? "");
        setDisplayName(data.user?.displayName ?? "");
      })
      .catch(() => {
        if (!isMounted) {
          return;
        }
        setUserId("");
        setDisplayName("");
      });
    return () => {
      isMounted = false;
    };
  }, [accountId]);

  const handleLogout = async () => {
    const auth = getAuthState();
    try {
      await fetchJSON("/auth/logout", {
        method: "POST",
        body: JSON.stringify({ refreshToken: auth?.refreshToken ?? "" }),
      });
    } catch {
      // ignore logout errors
    } finally {
      clearAuthState();
      router.push("/login");
    }
  };

  if (!isReady) {
    return <div className="min-h-screen" />;
  }

  if (!isAuthed) {
    return (
      <div className="min-h-screen">
        <header className="mx-auto flex w-full max-w-6xl items-center justify-between px-6 pt-8">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full border border-[#e4d8c7] bg-white/70 text-lg font-semibold">
              B
            </div>
            <div>
              <p className="text-sm tracking-[0.2em] text-[#5c5d63]">BOOK</p>
              <p className="font-[var(--font-display)] text-xl font-semibold">
                Manager
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3 text-sm">
            <Link className="text-[#5c5d63] hover:text-[#1b1c1f]" href="/login">
              ログイン
            </Link>
            <Link
              className="rounded-full bg-[#c86b3c] px-4 py-2 text-white"
              href="/signup"
            >
              新規登録
            </Link>
          </div>
        </header>
        <main className="mx-auto flex min-h-[70vh] max-w-3xl flex-col items-center justify-center gap-4 px-6 text-center">
          <h1 className="font-[var(--font-display)] text-3xl">
            ログインが必要です
          </h1>
          <p className="text-sm text-[#5c5d63]">
            所蔵管理やブックマーク機能を使うにはログインしてください。
          </p>
          <div className="flex gap-3">
            <Link
              className="rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white"
              href="/login"
            >
              ログインする
            </Link>
            <Link
              className="rounded-full border border-[#1b1c1f] px-5 py-3 text-sm font-medium text-[#1b1c1f]"
              href="/signup"
            >
              新規登録
            </Link>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <header className="sticky top-0 z-10 border-b border-[#e4d8c7] bg-[#f6f1e7]/90 backdrop-blur">
        <div className="mx-auto flex w-full max-w-6xl items-center justify-between px-6 py-4">
          <div className="flex items-center gap-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-full border border-[#e4d8c7] bg-white text-sm font-semibold">
              B
            </div>
            <div>
              <p className="text-[10px] tracking-[0.25em] text-[#5c5d63]">
                BOOK
              </p>
              <p className="font-[var(--font-display)] text-lg font-semibold">
                Manager
              </p>
            </div>
          </div>
          <div className="flex items-center gap-4">
            <div className="hidden items-center gap-2 rounded-full border border-[#e4d8c7] bg-white/80 px-4 py-2 text-xs text-[#5c5d63] md:flex">
              <span className="h-2 w-2 rounded-full bg-[#c86b3c]" />
              同期済み
            </div>
            <div className="flex items-center gap-2 text-sm">
              <div className="text-right">
                <span className="block text-[#1b1c1f]">
                  {displayName || userId || "user"}
                </span>
                <span className="block text-xs text-[#5c5d63]">
                  @{userId || accountId || "user"}
                </span>
              </div>
              <div className="h-9 w-9 rounded-full bg-[#1b1c1f] text-center text-sm leading-9 text-white">
                {(displayName || userId || "U").charAt(0).toUpperCase()}
              </div>
            </div>
            <button
              className="rounded-full border border-[#e4d8c7] px-3 py-2 text-xs text-[#5c5d63] hover:bg-white"
              type="button"
              onClick={handleLogout}
            >
              ログアウト
            </button>
          </div>
        </div>
      </header>
      <div className="mx-auto grid w-full max-w-6xl gap-6 px-6 py-8 lg:grid-cols-[220px_1fr]">
        <aside className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-4 shadow-sm">
          <nav className="flex flex-col gap-1 text-sm">
            {navItems.map((item) => (
              <Link
                key={item.href}
                className="rounded-2xl px-3 py-2 text-[#5c5d63] transition hover:bg-[#efe5d4] hover:text-[#1b1c1f]"
                href={item.href}
              >
                {item.label}
              </Link>
            ))}
          </nav>
          <div className="mt-6 rounded-2xl border border-[#e4d8c7] bg-[#f6f1e7] p-3 text-xs text-[#5c5d63]">
            <p className="font-medium text-[#1b1c1f]">書誌情報の修正</p>
            <p className="mt-2">
              誤り報告はフォームから送信。運営が確認後に反映します。
            </p>
          </div>
        </aside>
        <main className="flex min-h-[70vh] flex-col gap-6">
          {children}
        </main>
      </div>
    </div>
  );
}
