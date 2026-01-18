"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";

import { fetchJSON } from "@/lib/api";
import { setAuthState } from "@/lib/auth";

const errorMessages: Record<string, string> = {
  email_required: "メールアドレスを入力してください。",
  password_required: "パスワードを入力してください。",
  invalid_email: "メールアドレスの形式が正しくありません。",
  password_too_short: "パスワードは8文字以上で入力してください。",
  user_id_required: "ユーザーIDを入力してください。",
  user_id_too_short: "ユーザーIDは2文字以上で入力してください。",
  user_id_too_long: "ユーザーIDは20文字以内で入力してください。",
  display_name_too_long: "表示名は50文字以内で入力してください。",
  email_exists: "このメールアドレスは既に登録されています。",
  user_id_exists: "このユーザーIDは既に使用されています。",
};

export default function SignupPage() {
  const router = useRouter();
  const [form, setForm] = useState({
    userId: "",
    displayName: "",
    email: "",
    password: "",
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
        user: { id: string };
      }>("/auth/signup", {
        method: "POST",
        body: JSON.stringify(form),
      });
      setAuthState({
        accessToken: data.accessToken,
        refreshToken: data.refreshToken,
        userId: data.user.id,
      });
      router.push("/books");
    } catch (err) {
      const message = err instanceof Error ? err.message.trim() : "";
      setError(errorMessages[message] || "登録に失敗しました。");
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
              Create account
            </p>
            <h1 className="mt-3 font-[var(--font-display)] text-3xl">
              新しい本棚をはじめる
            </h1>
            <p className="mt-3 text-sm leading-6 text-[#5c5d63]">
              ISBN検索、シリーズ自動判定、次に買う本まで。最初の登録から一緒に準備しましょう。
            </p>
          </div>
          <div className="rounded-2xl border border-[#e4d8c7] bg-[#f6f1e7] p-4 text-xs text-[#5c5d63]">
            <p className="font-medium text-[#1b1c1f]">無料ではじめる</p>
            <p className="mt-2">メールアドレスとユーザーIDだけで開始できます。</p>
          </div>
        </div>

        <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
          <label className="text-sm text-[#1b1c1f]">
            ユーザーID
            <input
              className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              name="userId"
              placeholder="book"
              type="text"
              value={form.userId}
              onChange={(event) =>
                setForm((prev) => ({ ...prev, userId: event.target.value }))
              }
            />
            <span className="mt-1 block text-xs text-[#5c5d63]">
              2〜20文字、ログインや識別に使用されます。<br />※変更不可です。
            </span>
          </label>
          <label className="text-sm text-[#1b1c1f]">
            表示名
            <input
              className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              name="displayName"
              placeholder="book-name"
              type="text"
              value={form.displayName}
              onChange={(event) =>
                setForm((prev) => ({ ...prev, displayName: event.target.value }))
              }
            />
            <span className="mt-1 block text-xs text-[#5c5d63]">
              プロフィールに表示される名前
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
          {error ? <p className="text-xs text-red-600">{error}</p> : null}
          <button
            className="mt-2 rounded-full bg-[#c86b3c] px-5 py-3 text-sm font-medium text-white shadow-lg shadow-[#c86b3c]/20 transition hover:bg-[#8f3d1f]"
            type="submit"
            disabled={isSubmitting}
          >
            {isSubmitting ? "作成中..." : "アカウント作成"}
          </button>
          <div className="flex items-center justify-between text-xs text-[#5c5d63]">
            <span>登録すると<Link href="/terms" className="hover:text-[#1b1c1f] underline">利用規約</Link>に同意したことになります。</span>
            <Link className="hover:text-[#1b1c1f]" href="/login">
              ログインへ
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
}
