"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { fetchJSON } from "@/lib/api";
import { getAuthState } from "@/lib/auth";

type ProfileResponse = {
  user: { id: string; email: string; userId: string; displayName: string };
  settings: { visibility: string };
};

const errorMessages: Record<string, string> = {
  user_id_required: "ユーザーIDを入力してください。",
  user_id_too_short: "ユーザーIDは2文字以上で入力してください。",
  user_id_too_long: "ユーザーIDは20文字以内で入力してください。",
  user_id_exists: "このユーザーIDは既に使用されています。",
  display_name_too_long: "表示名は50文字以内で入力してください。",
  email_required: "メールアドレスを入力してください。",
  invalid_email: "メールアドレスの形式が正しくありません。",
  email_exists: "このメールアドレスは既に登録されています。",
};

export default function ProfileSettingsPage() {
  const [visibility, setVisibility] = useState("public");
  const [profileForm, setProfileForm] = useState({
    userId: "",
    displayName: "",
    email: "",
  });
  const [profileError, setProfileError] = useState<string | null>(null);
  const [profileMessage, setProfileMessage] = useState<string | null>(null);
  const [settingError, setSettingError] = useState<string | null>(null);

  useEffect(() => {
    const auth = getAuthState();
    if (!auth?.userId) {
      return;
    }
    fetchJSON<ProfileResponse>(`/users/${auth.userId}`, { auth: true })
      .then((data) => {
        if (data.settings?.visibility) {
          setVisibility(data.settings.visibility);
        }
        setProfileForm({
          userId: data.user?.userId ?? "",
          displayName: data.user?.displayName ?? "",
          email: data.user?.email ?? "",
        });
      })
      .catch(() => {
        setSettingError("プロフィール情報を取得できませんでした。");
      });
  }, []);

  const handleSaveVisibility = async () => {
    setSettingError(null);
    try {
      await fetchJSON("/users/me/settings", {
        method: "PATCH",
        auth: true,
        body: JSON.stringify({ visibility }),
      });
    } catch {
      setSettingError("公開設定の保存に失敗しました。");
    }
  };

  const handleSaveProfile = async () => {
    setProfileError(null);
    setProfileMessage(null);
    try {
      await fetchJSON("/users/me", {
        method: "PATCH",
        auth: true,
        body: JSON.stringify({
          userId: profileForm.userId,
          displayName: profileForm.displayName,
          email: profileForm.email,
        }),
      });
      setProfileMessage("プロフィールを更新しました。");
    } catch (err) {
      const message = err instanceof Error ? err.message.trim() : "";
      setProfileError(errorMessages[message] || "プロフィールの更新に失敗しました。");
    }
  };

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <Link className="text-xs text-[#5c5d63]" href="/settings">
          ← 設定に戻る
        </Link>
        <h1 className="mt-3 font-[var(--font-display)] text-3xl">
          プロフィール公開
        </h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          公開範囲の変更と表示設定。
        </p>
      </section>

      <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
        <h2 className="font-[var(--font-display)] text-2xl">
          ユーザー情報
        </h2>
        <div className="mt-4 flex flex-col gap-3 text-sm">
          <label className="flex flex-col gap-2">
            ユーザーID
            <input
              className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              value={profileForm.userId}
              onChange={(event) =>
                setProfileForm((prev) => ({ ...prev, userId: event.target.value }))
              }
            />
            <span className="text-xs text-[#5c5d63]">
              2〜20文字、ログインや識別に使用（変更可能）
            </span>
          </label>
          <label className="flex flex-col gap-2">
            表示名
            <input
              className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              value={profileForm.displayName}
              onChange={(event) =>
                setProfileForm((prev) => ({
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
              value={profileForm.email}
              onChange={(event) =>
                setProfileForm((prev) => ({ ...prev, email: event.target.value }))
              }
            />
          </label>
        </div>
        {profileError ? (
          <p className="mt-3 text-xs text-red-600">{profileError}</p>
        ) : null}
        {profileMessage ? (
          <p className="mt-3 text-xs text-[#5c5d63]">{profileMessage}</p>
        ) : null}
        <button
          className="mt-4 rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white"
          type="button"
          onClick={handleSaveProfile}
        >
          プロフィールを保存
        </button>
      </section>

      <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
        <div className="flex flex-col gap-3 text-sm">
          <label className="flex items-center justify-between rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3">
            <span>公開</span>
            <input
              type="radio"
              name="visibility"
              value="public"
              checked={visibility === "public"}
              onChange={() => setVisibility("public")}
            />
          </label>
          <label className="flex items-center justify-between rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3">
            <span>フォロワー限定</span>
            <input
              type="radio"
              name="visibility"
              value="followers"
              checked={visibility === "followers"}
              onChange={() => setVisibility("followers")}
            />
          </label>
        </div>
        {settingError ? (
          <p className="mt-3 text-xs text-red-600">{settingError}</p>
        ) : null}
        <button
          className="mt-4 rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white"
          type="button"
          onClick={handleSaveVisibility}
        >
          変更を保存
        </button>
      </section>
    </div>
  );
}
