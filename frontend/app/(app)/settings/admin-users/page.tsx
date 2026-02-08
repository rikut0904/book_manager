"use client";

import Link from "next/link";
import { useCallback, useEffect, useMemo, useState } from "react";

import { fetchJSON } from "@/lib/api";
import { getAuthState } from "@/lib/auth";

type SettingsResponse = {
  isAdmin?: boolean;
};

type AdminUserInfo = {
  userId: string;
  source: "env" | "db";
  createdBy?: string;
  createdAt?: string;
  inDb?: boolean;
};

type AdminUsersResponse = {
  items: AdminUserInfo[];
  total: number;
};

const pageSizeOptions = [20, 50, 100, 200] as const;

export default function AdminUsersSettingsPage() {
  const [isAdmin, setIsAdmin] = useState(false);
  const [admins, setAdmins] = useState<AdminUserInfo[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState<(typeof pageSizeOptions)[number]>(50);
  const [formUserId, setFormUserId] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [message, setMessage] = useState<string | null>(null);

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

  const loadAdmins = useCallback(
    async (targetPage = page, targetPageSize = pageSize) => {
    setLoading(true);
    setError(null);
    try {
      const data = await fetchJSON<AdminUsersResponse>(
        `/admin/users?page=${targetPage}&pageSize=${targetPageSize}`,
        { auth: true }
      );
      setAdmins(data.items ?? []);
      setTotal(data.total ?? 0);
    } catch {
      setAdmins([]);
      setTotal(0);
      setError("管理者一覧の取得に失敗しました。");
    } finally {
      setLoading(false);
    }
  },
    [page, pageSize]
  );

  useEffect(() => {
    if (!isAdmin) {
      return;
    }
    loadAdmins();
  }, [isAdmin, loadAdmins]);

  const totalPages = useMemo(() => {
    if (total === 0) return 1;
    return Math.max(1, Math.ceil(total / pageSize));
  }, [total, pageSize]);

  const formatDate = (dateStr?: string) => {
    if (!dateStr) return "";
    const date = new Date(dateStr);
    return date.toLocaleString("ja-JP", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const handleGrant = async () => {
    setError(null);
    setMessage(null);
    const userId = formUserId.trim();
    if (!userId) {
      setError("ユーザーIDを入力してください。");
      return;
    }
    if (!window.confirm(`ユーザーID「${userId}」に管理者権限を付与します。よろしいですか？`)) {
      return;
    }
    try {
      await fetchJSON("/admin/users", {
        method: "POST",
        auth: true,
        body: JSON.stringify({ userId }),
      });
      setFormUserId("");
      setMessage("管理者権限を付与しました。");
      await loadAdmins(1, pageSize);
      setPage(1);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "管理者付与に失敗しました。";
      if (errorMessage.includes("already_admin")) {
        setError("既に管理者権限が付与されています。");
      } else if (errorMessage.includes("not found")) {
        setError("指定したユーザーが見つかりません。");
      } else {
        setError(errorMessage);
      }
    }
  };

  const handleRemove = async (userId: string) => {
    setError(null);
    setMessage(null);
    if (!window.confirm(`ユーザーID「${userId}」の管理者権限を削除します。よろしいですか？`)) {
      return;
    }
    try {
      await fetchJSON(`/admin/users/${userId}?confirm=true`, {
        method: "DELETE",
        auth: true,
      });
      setMessage("管理者権限を削除しました。");
      await loadAdmins();
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "管理者削除に失敗しました。";
      if (errorMessage.includes("cannot_remove_env_admin")) {
        setError("ENVで管理されている管理者は削除できません。");
      } else if (errorMessage.includes("not found")) {
        setError("指定した管理者が見つかりません。");
      } else {
        setError(errorMessage);
      }
    }
  };

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <Link className="text-xs text-[#5c5d63]" href="/settings">
          ← 設定に戻る
        </Link>
        <h1 className="mt-3 font-[var(--font-display)] text-3xl">
          管理者管理
        </h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          管理者権限の付与・削除と一覧の確認を行います。
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
            <h2 className="font-[var(--font-display)] text-2xl">権限付与</h2>
            <p className="mt-2 text-sm text-[#5c5d63]">
              既存ユーザーのユーザーIDを指定して管理者権限を付与します。
            </p>
            <div className="mt-4 grid gap-3 text-sm md:grid-cols-[1fr_auto]">
              <input
                className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
                placeholder="ユーザーID（必須）"
                value={formUserId}
                onChange={(e) => setFormUserId(e.target.value)}
              />
              <button
                className="rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white"
                type="button"
                onClick={handleGrant}
              >
                付与する
              </button>
            </div>
            {error ? (
              <p className="mt-3 text-xs text-red-600">{error}</p>
            ) : null}
            {message ? (
              <p className="mt-3 text-xs text-green-600">{message}</p>
            ) : null}
          </section>

          <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
            <div className="flex flex-wrap items-center justify-between gap-3">
              <div>
                <h2 className="font-[var(--font-display)] text-2xl">管理者一覧</h2>
                <p className="mt-2 text-sm text-[#5c5d63]">
                  付与済みの管理者のみ表示されます。
                </p>
              </div>
              <div className="flex items-center gap-2 text-xs text-[#5c5d63]">
                <span>表示数</span>
                <select
                  className="rounded-full border border-[#e4d8c7] bg-white px-3 py-2"
                  value={pageSize}
                  onChange={(e) => {
                    const value = Number(e.target.value) as (typeof pageSizeOptions)[number];
                    setPageSize(value);
                    setPage(1);
                  }}
                >
                  {pageSizeOptions.map((size) => (
                    <option key={size} value={size}>
                      {size}件
                    </option>
                  ))}
                </select>
              </div>
            </div>

            {loading ? (
              <p className="mt-4 text-sm text-[#5c5d63]">読み込み中...</p>
            ) : admins.length === 0 ? (
              <p className="mt-4 text-sm text-[#5c5d63]">
                管理者がいません。
              </p>
            ) : (
              <div className="mt-4 flex flex-col gap-3">
                {admins.map((item) => {
                  const canRemove = !(item.source === "env" && !item.inDb);
                  return (
                    <div
                      key={item.userId}
                      className="rounded-2xl border border-[#e4d8c7] bg-white p-4"
                    >
                      <div className="flex flex-wrap items-start justify-between gap-3">
                        <div>
                          <p className="text-sm font-medium text-[#1b1c1f]">
                            ユーザーID: {item.userId}
                          </p>
                          <div className="mt-2 flex flex-wrap items-center gap-2 text-xs">
                            <span className="rounded-full border border-[#e4d8c7] bg-[#f7f2ea] px-2 py-1 text-[#5c5d63]">
                              {item.source === "env" ? "ENV" : "DB"}
                            </span>
                            {item.source === "env" && !item.inDb ? (
                              <span className="rounded-full border border-[#f0c6a6] bg-[#fef7f0] px-2 py-1 text-[#c86b3c]">
                                ENV固定
                              </span>
                            ) : null}
                          </div>
                          {item.createdAt ? (
                            <p className="mt-2 text-xs text-[#5c5d63]">
                              付与日: {formatDate(item.createdAt)}
                            </p>
                          ) : null}
                          {item.createdBy ? (
                            <p className="mt-1 text-xs text-[#5c5d63]">
                              付与者: {item.createdBy}
                            </p>
                          ) : null}
                        </div>
                        <button
                          className={`text-xs ${
                            canRemove
                              ? "text-[#c86b3c] hover:text-[#8f3d1f]"
                              : "cursor-not-allowed text-[#b0a89b]"
                          }`}
                          type="button"
                          onClick={() => {
                            if (canRemove) {
                              void handleRemove(item.userId);
                            }
                          }}
                          disabled={!canRemove}
                        >
                          削除
                        </button>
                      </div>
                    </div>
                  );
                })}
              </div>
            )}

            <div className="mt-4 flex flex-wrap items-center justify-between gap-3 text-xs text-[#5c5d63]">
              <p>
                {total === 0
                  ? "全0件"
                  : `全${total}件 / ${page}ページ目`}
              </p>
              <div className="flex items-center gap-2">
                <button
                  className="rounded-full border border-[#e4d8c7] bg-white px-3 py-2 disabled:cursor-not-allowed disabled:opacity-60"
                  type="button"
                  disabled={page <= 1}
                  onClick={() => setPage((prev) => Math.max(1, prev - 1))}
                >
                  前へ
                </button>
                <button
                  className="rounded-full border border-[#e4d8c7] bg-white px-3 py-2 disabled:cursor-not-allowed disabled:opacity-60"
                  type="button"
                  disabled={page >= totalPages}
                  onClick={() => setPage((prev) => Math.min(totalPages, prev + 1))}
                >
                  次へ
                </button>
              </div>
            </div>
          </section>
        </>
      )}
    </div>
  );
}
