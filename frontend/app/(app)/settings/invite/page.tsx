"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { fetchJSON, APIError } from "@/lib/api";
import { getAuthState } from "@/lib/auth";

type SettingsResponse = {
  isAdmin?: boolean;
};

type InvitationItem = {
  id: string;
  userId: string;
  email: string;
  token?: string;
  createdBy: string;
  expiresAt: string;
  usedAt: string | null;
  usedBy: string;
  createdAt: string;
};

type CreateInvitationResponse = {
  id: string;
  token: string;
  userId: string;
  email: string;
  expiresAt: string;
  createdAt: string;
};

export default function InvitationsSettingsPage() {
  const [isAdmin, setIsAdmin] = useState(false);
  const [invitations, setInvitations] = useState<InvitationItem[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [message, setMessage] = useState<string | null>(null);
  const [form, setForm] = useState({ userId: "", email: "" });
  const [createdToken, setCreatedToken] = useState<string | null>(null);
  const [copiedId, setCopiedId] = useState<string | null>(null);

  useEffect(() => {
    const auth = getAuthState();
    if (!auth?.userId) {
      return;
    }
    fetchJSON<SettingsResponse>(`/users/${auth.userId}`, { auth: true })
      .then((data) => {
        setIsAdmin(Boolean(data.isAdmin));
      })
      .catch(() => {
        setIsAdmin(false);
      });
  }, []);

  useEffect(() => {
    if (!isAdmin) {
      return;
    }
    let isMounted = true;
    fetchJSON<{ items: InvitationItem[] }>("/admin/invitations", { auth: true })
      .then((data) => {
        if (!isMounted) return;
        setInvitations(data.items ?? []);
      })
      .catch(() => {
        if (!isMounted) return;
        setInvitations([]);
      });
    return () => {
      isMounted = false;
    };
  }, [isAdmin]);

  const getStatus = (item: InvitationItem) => {
    if (item.usedAt) {
      return { label: "使用済み", color: "text-[#5c5d63]" };
    }
    if (new Date(item.expiresAt) < new Date()) {
      return { label: "期限切れ", color: "text-red-600" };
    }
    return { label: "有効", color: "text-green-600" };
  };

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleString("ja-JP", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const handleCreate = async () => {
    setError(null);
    setMessage(null);
    setCreatedToken(null);

    if (!form.userId.trim()) {
      setError("ユーザーIDを入力してください。");
      return;
    }

    if (!form.email.trim()) {
      setError("メールアドレスを入力してください。");
      return;
    }

    try {
      const result = await fetchJSON<CreateInvitationResponse>(
        "/admin/invitations",
        {
          method: "POST",
          auth: true,
          body: JSON.stringify({
            userId: form.userId.trim(),
            email: form.email.trim(),
          }),
        }
      );
      setInvitations((prev) => [
        {
          id: result.id,
          userId: result.userId,
          email: result.email,
          createdBy: "",
          expiresAt: result.expiresAt,
          usedAt: null,
          usedBy: "",
          createdAt: result.createdAt,
        },
        ...prev,
      ]);
      setCreatedToken(result.token);
      setMessage("招待を作成しました。トークンをコピーして招待者に送信してください。");
      setForm({ userId: "", email: "" });
    } catch (err: unknown) {
      if (err instanceof APIError && err.data.existingUserId) {
        setError(`このメールアドレスは既に登録されています。（ユーザーID: ${err.data.existingUserId}）`);
        return;
      }
      const errorMessage = err instanceof Error ? err.message : "招待の作成に失敗しました。";
      if (errorMessage.includes("user_id_not_admin")) {
        setError("このユーザーIDは管理者IDリストに含まれていません。");
      } else if (errorMessage.includes("user_id_exists")) {
        setError("このユーザーIDは既に使用されています。");
      } else if (errorMessage.includes("user_id_already_invited")) {
        setError("このユーザーIDには既に有効な招待があります。");
      } else if (errorMessage.includes("email_exists")) {
        setError("このメールアドレスは既に登録されています。");
      } else if (errorMessage.includes("email is required")) {
        setError("メールアドレスを入力してください。");
      } else {
        setError(errorMessage);
      }
    }
  };

  const handleDelete = async (id: string) => {
    setError(null);
    setMessage(null);
    try {
      await fetchJSON(`/admin/invitations/${id}`, {
        method: "DELETE",
        auth: true,
      });
      setInvitations((prev) => prev.filter((item) => item.id !== id));
      setMessage("招待を削除しました。");
    } catch {
      setError("招待の削除に失敗しました。");
    }
  };

  const handleCopyToken = async (token: string, id: string) => {
    try {
      await navigator.clipboard.writeText(token);
      setCopiedId(id);
      setTimeout(() => setCopiedId(null), 2000);
    } catch {
      setError("クリップボードへのコピーに失敗しました。");
    }
  };

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <Link className="text-xs text-[#5c5d63]" href="/settings">
          ← 設定に戻る
        </Link>
        <h1 className="mt-3 font-[var(--font-display)] text-3xl">
          管理者招待
        </h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          新しい管理者ユーザーを招待できます。招待されたユーザーは予約された管理者IDで登録できます。
        </p>
      </section>

      {!isAdmin ? (
        <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
          <p className="text-sm text-[#5c5d63]">
            この機能は管理者ユーザーのみ利用できます。
          </p>
        </section>
      ) : (
        <>
          <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
            <h2 className="font-[var(--font-display)] text-2xl">新規招待作成</h2>
            <p className="mt-2 text-sm text-[#5c5d63]">
              新規ユーザーを招待します。ユーザーIDとメールアドレスを入力してください。既存ユーザーには招待できません。
            </p>
            <div className="mt-4 rounded-2xl border border-dashed border-[#e4d8c7] bg-white/80 p-4 text-sm text-[#5c5d63]">
              <p className="font-medium text-[#1b1c1f]">招待手順</p>
              <ol className="mt-2 list-decimal pl-5 text-xs leading-6">
                <li>新規ユーザー用のユーザーIDを入力</li>
                <li>招待先メールアドレスを入力（未登録のメールアドレスのみ）</li>
                <li>招待作成ボタンを押す</li>
                <li>招待メールが自動送信されます</li>
              </ol>
            </div>

            <div className="mt-4 grid gap-3 text-sm md:grid-cols-[1fr_1fr_auto]">
              <input
                className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
                placeholder="ユーザーID（必須）"
                value={form.userId}
                onChange={(e) =>
                  setForm((prev) => ({ ...prev, userId: e.target.value }))
                }
              />
              <input
                className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
                placeholder="メールアドレス（必須）"
                type="email"
                value={form.email}
                onChange={(e) =>
                  setForm((prev) => ({ ...prev, email: e.target.value }))
                }
              />
              <button
                className="rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white"
                type="button"
                onClick={handleCreate}
              >
                招待作成
              </button>
            </div>

            {error ? (
              <p className="mt-3 text-xs text-red-600">{error}</p>
            ) : null}
            {message ? (
              <p className="mt-3 text-xs text-green-600">{message}</p>
            ) : null}

            {createdToken ? (
              <div className="mt-4 rounded-2xl border border-[#c86b3c] bg-[#fef7f0] p-4">
                <p className="text-xs font-medium text-[#c86b3c]">
                  招待トークン（このトークンは一度しか表示されません）
                </p>
                <div className="mt-2 flex items-center gap-2">
                  <code className="flex-1 break-all rounded-lg bg-white px-3 py-2 text-xs text-[#1b1c1f]">
                    {createdToken}
                  </code>
                  <button
                    className="rounded-full bg-[#c86b3c] px-4 py-2 text-xs font-medium text-white"
                    type="button"
                    onClick={() => handleCopyToken(createdToken, "new")}
                  >
                    {copiedId === "new" ? "コピー済み" : "コピー"}
                  </button>
                </div>
              </div>
            ) : null}
          </section>

          <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
            <h2 className="font-[var(--font-display)] text-2xl">招待一覧</h2>
            <p className="mt-2 text-sm text-[#5c5d63]">
              作成された招待の一覧です。有効期限は作成から72時間です。
            </p>

            {invitations.length === 0 ? (
              <p className="mt-4 text-sm text-[#5c5d63]">招待はありません。</p>
            ) : (
              <div className="mt-4 flex flex-col gap-3">
                {invitations.map((item) => {
                  const status = getStatus(item);
                  return (
                    <div
                      key={item.id}
                      className="rounded-2xl border border-[#e4d8c7] bg-white p-4"
                    >
                      <div className="flex flex-wrap items-start justify-between gap-2">
                        <div>
                          <p className="text-sm font-medium text-[#1b1c1f]">
                            ユーザーID: {item.userId}
                          </p>
                          {item.email ? (
                            <p className="mt-1 text-xs text-[#5c5d63]">
                              メール: {item.email}
                            </p>
                          ) : null}
                          <p className="mt-1 text-xs text-[#5c5d63]">
                            作成日: {formatDate(item.createdAt)}
                          </p>
                          <p className="mt-1 text-xs text-[#5c5d63]">
                            有効期限: {formatDate(item.expiresAt)}
                          </p>
                          {item.usedAt ? (
                            <p className="mt-1 text-xs text-[#5c5d63]">
                              使用日: {formatDate(item.usedAt)}
                            </p>
                          ) : null}
                        </div>
                        <div className="flex items-center gap-2">
                          <span className={`text-xs font-medium ${status.color}`}>
                            {status.label}
                          </span>
                          {!item.usedAt ? (
                            <button
                              className="text-xs text-[#c86b3c] hover:text-[#8f3d1f]"
                              type="button"
                              onClick={() => handleDelete(item.id)}
                            >
                              削除
                            </button>
                          ) : null}
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </section>
        </>
      )}
    </div>
  );
}
