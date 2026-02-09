"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";

import { fetchJSON } from "@/lib/api";
import { setAuthState, updateAuthState } from "@/lib/auth";
import { loginErrorMessages } from "@/lib/errorMessages";

export default function LoginPage() {
  const router = useRouter();
  const [form, setForm] = useState({ email: "", password: "" });
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
        user: { id: string };
        emailVerified: boolean;
      }>("/auth/login", {
        method: "POST",
        body: JSON.stringify(form),
      });
      setAuthState({
        accessToken: data.accessToken,
        refreshToken: data.refreshToken,
        userId: data.user.id,
      });
      try {
        const profile = await fetchJSON<{ user: { displayName: string } }>(
          "/users/profile",
          { auth: true }
        );
        if (profile.user?.displayName) {
          updateAuthState({ displayName: profile.user.displayName });
        }
      } catch {
        // ignore profile fetch errors
      }
      if (data.emailVerified) {
        router.push("/books");
      } else {
        router.push("/verify-email");
      }
    } catch (err) {
      const message = err instanceof Error ? err.message.trim() : "";
      setError(loginErrorMessages[message] || "ログインに失敗しました。");
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
              Sign in
            </p>
            <h1 className="mt-3 font-[var(--font-display)] text-3xl">
              本棚に戻る
            </h1>
            <p className="mt-3 text-sm leading-6 text-[#5c5d63]">
              ログインして所蔵の更新、シリーズの確認、次に買う本の管理を続けましょう。
            </p>
          </div>
        </div>

        <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
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
          </label>
          {error ? <p className="text-xs text-red-600">{error}</p> : null}
          <button
            className="mt-2 rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white shadow-lg shadow-[#1b1c1f]/10 transition hover:bg-black"
            type="submit"
            disabled={isSubmitting}
          >
            {isSubmitting ? "ログイン中..." : "ログイン"}
          </button>
          <div className="flex items-center justify-between text-xs text-[#5c5d63]">
            <Link className="hover:text-[#1b1c1f]" href="#">
              パスワードを忘れた
            </Link>
            <Link className="hover:text-[#1b1c1f]" href="/signup">
              新規登録へ
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
}
