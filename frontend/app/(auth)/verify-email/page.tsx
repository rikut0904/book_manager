"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import { fetchJSON } from "@/lib/api";
import { getAuthState, setAuthState } from "@/lib/auth";

type StatusResponse = {
  user?: { email?: string; emailVerified?: boolean };
};

export default function VerifyEmailPage() {
  const router = useRouter();
  const [status, setStatus] = useState<"waiting" | "verified">("waiting");
  const [error, setError] = useState<string | null>(null);
  const [currentEmail, setCurrentEmail] = useState("");
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editEmail, setEditEmail] = useState("");
  const [editError, setEditError] = useState<string | null>(null);
  const [message, setMessage] = useState<string | null>(null);

  useEffect(() => {
    const auth = getAuthState();
    if (!auth?.accessToken) {
      router.push("/login");
      return;
    }

    const checkStatus = async () => {
      setError(null);
      try {
        const currentAuth = getAuthState();
        if (!currentAuth?.refreshToken) {
          throw new Error("missing refresh token");
        }
        const refreshed = await fetchJSON<{
          accessToken: string;
          refreshToken: string;
        }>("/auth/refresh", {
          method: "POST",
          body: JSON.stringify({ refreshToken: currentAuth.refreshToken }),
        });
        setAuthState({
          accessToken: refreshed.accessToken,
          refreshToken: refreshed.refreshToken,
          userId: currentAuth.userId,
        });
        const data = await fetchJSON<StatusResponse>("/auth/status", { auth: true });
        if (data.user?.email) {
          setCurrentEmail(data.user.email);
        }
        if (data.user?.emailVerified) {
          setStatus("verified");
          router.push("/books");
        } else {
          setStatus("waiting");
        }
      } catch {
        setError("認証状態を確認できませんでした。");
      }
    };

    checkStatus();
    const timer = window.setInterval(checkStatus, 5000);
    return () => window.clearInterval(timer);
  }, [router]);

  const handleResend = async () => {
    setMessage(null);
    setError(null);
    setEditError(null);
    const auth = getAuthState();
    if (!auth?.refreshToken) {
      setError("認証情報が見つかりませんでした。");
      return;
    }
    try {
      if (!editEmail.trim()) {
        setEditError("メールアドレスを入力してください。");
        return;
      }
      let currentRefreshToken = auth.refreshToken;
      // メールアドレスが変更された場合のみ、メールアドレス更新APIを呼び出す
      if (editEmail.trim() !== currentEmail) {
        const data = await fetchJSON<{
          accessToken: string;
          refreshToken: string;
        }>("/auth/email", {
          method: "PATCH",
          body: JSON.stringify({
            email: editEmail.trim(),
            refreshToken: auth.refreshToken,
          }),
        });
        setAuthState({
          accessToken: data.accessToken,
          refreshToken: data.refreshToken,
          userId: auth.userId,
        });
        currentRefreshToken = data.refreshToken;
        setCurrentEmail(editEmail.trim());
      }
      await fetchJSON("/auth/resend-verify", {
        method: "POST",
        body: JSON.stringify({ refreshToken: currentRefreshToken }),
      });
      setIsModalOpen(false);
      setMessage("確認メールを再送しました。");
    } catch (err) {
      const messageText = err instanceof Error ? err.message.trim() : "";
      if (messageText === "email_exists") {
        setEditError("このメールアドレスは既に登録されています。");
        return;
      }
      if (messageText === "invalid_email") {
        setEditError("メールアドレスの形式が正しくありません。");
        return;
      }
      setError("再送に失敗しました。");
    }
  };

  const handleOpenModal = () => {
    setEditEmail(currentEmail);
    setEditError(null);
    setIsModalOpen(true);
  };

  return (
    <div className="mx-auto flex min-h-screen w-full max-w-3xl items-center justify-center px-6 py-16">
      <div className="w-full rounded-[32px] border border-[#e4d8c7] bg-white/80 p-8 shadow-lg">
        <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
          Verify
        </p>
        <h1 className="mt-3 font-[var(--font-display)] text-3xl">
          メール認証を完了してください
        </h1>
        <p className="mt-3 text-sm leading-6 text-[#5c5d63]">
          登録したメールアドレスに確認メールを送信しました。認証が完了すると自動でホームへ移動します。
        </p>
        <div className="mt-6 flex flex-col gap-4">
          <div className="flex items-center gap-3 rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3">
            <span className="h-2 w-2 animate-pulse rounded-full bg-[#c86b3c]" />
            <span className="text-sm text-[#1b1c1f]">認証待ち...</span>
          </div>
          <div className="flex flex-wrap gap-3">
            <Link
              className="inline-flex w-fit items-center gap-2 rounded-full border border-[#e4d8c7] px-4 py-2 text-sm text-[#5c5d63] hover:bg-white"
              href="/login"
            >
              ログイン画面へ戻る
            </Link>
            <button
              className="rounded-full border border-[#e4d8c7] px-4 py-2 text-sm text-[#5c5d63] hover:bg-white"
              type="button"
              onClick={handleOpenModal}
            >
              もう一度送信する
            </button>
          </div>
        </div>
        {error ? (
          <p className="mt-3 text-xs text-red-600">{error}</p>
        ) : null}
        {message ? (
          <p className="mt-3 text-xs text-[#5c5d63]">{message}</p>
        ) : null}
      </div>
      {isModalOpen ? (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4">
          <div className="w-full max-w-md rounded-3xl border border-[#e4d8c7] bg-white p-6 shadow-xl">
            <div className="flex items-start justify-between gap-4">
              <div>
                <h2 className="font-[var(--font-display)] text-2xl">
                  メールアドレスの確認
                </h2>
              </div>
            </div>
            <div className="mt-4 flex flex-col gap-2 text-sm">
              <label className="flex flex-col gap-2">
                <input
                  className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
                  type="email"
                  value={editEmail}
                  onChange={(event) => setEditEmail(event.target.value)}
                />
              </label>
            </div>
            {error ? (
              <p className="mt-3 text-xs text-red-600">{error}</p>
            ) : null}
            <div className="mt-4 flex justify-end gap-2 text-sm">
              <button
                className="rounded-full border border-[#e4d8c7] px-4 py-2 text-[#5c5d63] hover:bg-[#f6f1e7]"
                type="button"
                onClick={() => setIsModalOpen(false)}
              >
                キャンセル
              </button>
              <button
                className="rounded-full bg-[#1b1c1f] px-5 py-2 text-white"
                type="button"
                onClick={handleResend}
              >
                確認して送信
              </button>
            </div>
          </div>
        </div>
      ) : null}
    </div>
  );
}
