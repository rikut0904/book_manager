"use client";

import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useState } from "react";

import { fetchJSON } from "@/lib/api";
import { setAuthState } from "@/lib/auth";

const adminSignupErrorMessages: Record<string, string> = {
  token_required: "招待トークンが必要です。",
  invalid_token: "無効な招待トークンです。",
  token_expired: "招待トークンの有効期限が切れています。",
  token_already_used: "この招待トークンは既に使用されています。",
  email_mismatch: "メールアドレスが招待と一致しません。",
  email_required: "メールアドレスを入力してください。",
  invalid_email: "有効なメールアドレスを入力してください。",
  password_required: "パスワードを入力してください。",
  password_too_short: "パスワードは8文字以上で入力してください。",
  display_name_too_long: "表示名は50文字以内で入力してください。",
  email_exists: "このメールアドレスは既に登録されています。",
};

function AdminSignupForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const tokenFromUrl = searchParams.get("token") || "";
  const emailFromUrl = searchParams.get("email") || "";

  const [form, setForm] = useState({
    token: tokenFromUrl,
    email: emailFromUrl,
    password: "",
    displayName: "",
  });
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setIsSubmitting(true);

    try {
      const data = await fetchJSON<{
        accessToken: string;
        refreshToken: string;
        user: { id: string; userId: string };
        emailVerified: boolean;
        isAdmin: boolean;
      }>("/auth/signup/admin", {
        method: "POST",
        body: JSON.stringify({
          token: form.token.trim(),
          email: form.email.trim(),
          password: form.password,
          displayName: form.displayName.trim() || undefined,
        }),
      });

      setAuthState({
        accessToken: data.accessToken,
        refreshToken: data.refreshToken,
        userId: data.user.id,
      });

      if (data.emailVerified) {
        router.push("/books");
      } else {
        router.push("/verify-email");
      }
    } catch (err) {
      const message = err instanceof Error ? err.message.trim() : "";
      setError(
        adminSignupErrorMessages[message] ||
          "登録に失敗しました。招待トークンを確認してください。"
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="mx-auto flex min-h-screen w-full max-w-4xl items-center justify-center px-6 py-16">
      <div className="grid w-full gap-8 rounded-[32px] border border-[#e4d8c7] bg-white/80 p-8 shadow-lg md:grid-cols-[1.1fr_0.9fr]">
        <div className="flex flex-col justify-between gap-6">
          <div>
            <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
              Admin Registration
            </p>
            <h1 className="mt-3 font-[var(--font-display)] text-3xl">
              管理者アカウント登録
            </h1>
            <p className="mt-3 text-sm leading-6 text-[#5c5d63]">
              招待を受けた管理者ユーザーとして登録します。招待トークンと登録情報を入力してください。
            </p>
          </div>
          <div className="rounded-2xl border border-[#c86b3c]/30 bg-[#fef7f0] p-4 text-xs text-[#5c5d63]">
            <p className="font-medium text-[#c86b3c]">管理者招待について</p>
            <p className="mt-2">
              このページは招待を受けた方専用です。招待トークンがない場合は、管理者にお問い合わせください。
            </p>
          </div>
        </div>

        <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
          <label className="text-sm text-[#1b1c1f]">
            招待トークン
            <input
              className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              name="token"
              placeholder="招待トークンを入力"
              type="text"
              value={form.token}
              onChange={(event) =>
                setForm((prev) => ({ ...prev, token: event.target.value }))
              }
            />
            <span className="mt-1 block text-xs text-[#5c5d63]">
              管理者から受け取った招待トークン
            </span>
          </label>
          <label className="text-sm text-[#1b1c1f]">
            メールアドレス
            <input
              className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              name="email"
              placeholder="user@example.com"
              type="email"
              value={form.email}
              onChange={(event) =>
                setForm((prev) => ({ ...prev, email: event.target.value }))
              }
            />
          </label>
          <label className="text-sm text-[#1b1c1f]">
            パスワード
            <input
              className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              name="password"
              placeholder="********"
              type="password"
              value={form.password}
              onChange={(event) =>
                setForm((prev) => ({ ...prev, password: event.target.value }))
              }
            />
            <span className="mt-1 block text-xs text-[#5c5d63]">
              8文字以上
            </span>
          </label>
          <label className="text-sm text-[#1b1c1f]">
            表示名（任意）
            <input
              className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              name="displayName"
              placeholder="表示名"
              type="text"
              value={form.displayName}
              onChange={(event) =>
                setForm((prev) => ({ ...prev, displayName: event.target.value }))
              }
            />
            <span className="mt-1 block text-xs text-[#5c5d63]">
              プロフィールに表示される名前（未入力の場合はメールアドレスが使用されます）
            </span>
          </label>
          {error ? <p className="text-xs text-red-600">{error}</p> : null}
          <button
            className="mt-2 rounded-full bg-[#c86b3c] px-5 py-3 text-sm font-medium text-white shadow-lg shadow-[#c86b3c]/20 transition hover:bg-[#8f3d1f]"
            type="submit"
            disabled={isSubmitting}
          >
            {isSubmitting ? "登録中..." : "管理者として登録"}
          </button>
          <div className="flex items-center justify-between text-xs text-[#5c5d63]">
            <span>
              登録すると
              <Link href="/terms" className="underline hover:text-[#1b1c1f]">
                利用規約
              </Link>
              に同意したことになります。
            </span>
            <Link className="hover:text-[#1b1c1f]" href="/login">
              ログインへ
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
}

export default function AdminSignupPage() {
  return (
    <Suspense
      fallback={
        <div className="flex min-h-screen items-center justify-center">
          <p className="text-sm text-[#5c5d63]">読み込み中...</p>
        </div>
      }
    >
      <AdminSignupForm />
    </Suspense>
  );
}
