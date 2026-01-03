"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { fetchJSON } from "@/lib/api";
import { getAuthState } from "@/lib/auth";

type SettingsResponse = {
  settings?: {
    geminiEnabled?: boolean;
    geminiModel?: string;
    geminiHasKey?: boolean;
  };
};

const defaultModel = "gemini-1.5-flash";

export default function AISettingsPage() {
  const [geminiEnabled, setGeminiEnabled] = useState(false);
  const [geminiModel, setGeminiModel] = useState(defaultModel);
  const [geminiApiKey, setGeminiApiKey] = useState("");
  const [geminiHasKey, setGeminiHasKey] = useState(false);
  const [settingError, setSettingError] = useState<string | null>(null);
  const [settingMessage, setSettingMessage] = useState<string | null>(null);

  useEffect(() => {
    const auth = getAuthState();
    if (!auth?.userId) {
      return;
    }
    fetchJSON<SettingsResponse>(`/users/${auth.userId}`, { auth: true })
      .then((data) => {
        const settings = data.settings ?? {};
        setGeminiEnabled(Boolean(settings.geminiEnabled));
        setGeminiModel(settings.geminiModel || defaultModel);
        setGeminiHasKey(Boolean(settings.geminiHasKey));
      })
      .catch(() => {
        setSettingError("AI設定を取得できませんでした。");
      });
  }, []);

  const handleSave = async () => {
    setSettingError(null);
    setSettingMessage(null);
    try {
      await fetchJSON("/users/me/settings", {
        method: "PATCH",
        auth: true,
        body: JSON.stringify({
          geminiEnabled,
          geminiModel: geminiModel.trim(),
          geminiApiKey: geminiApiKey.trim(),
        }),
      });
      setSettingMessage("AI設定を保存しました。");
      if (geminiApiKey.trim()) {
        setGeminiHasKey(true);
        setGeminiApiKey("");
      }
    } catch {
      setSettingError("AI設定の保存に失敗しました。");
    }
  };

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <Link className="text-xs text-[#5c5d63]" href="/settings">
          ← 設定に戻る
        </Link>
        <h1 className="mt-3 font-[var(--font-display)] text-3xl">
          AIシリーズ判定
        </h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          Gemini APIを利用したシリーズ推定の設定を行います。
        </p>
      </section>

      <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
        <div className="flex flex-col gap-4 text-sm">
          <label className="flex items-center justify-between rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3">
            <span>Geminiの利用に同意して有効化</span>
            <input
              type="checkbox"
              checked={geminiEnabled}
              onChange={(event) => setGeminiEnabled(event.target.checked)}
            />
          </label>
          <div className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3">
            <p className="text-xs text-[#5c5d63]">モデル名</p>
            <input
              className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              list="gemini-models"
              value={geminiModel}
              onChange={(event) => setGeminiModel(event.target.value)}
              placeholder={defaultModel}
            />
            <datalist id="gemini-models">
              <option value="gemini-1.5-flash" />
              <option value="gemini-1.5-pro" />
              <option value="gemini-2.0-flash" />
            </datalist>
          </div>
          <div className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3">
            <p className="text-xs text-[#5c5d63]">APIキー（任意）</p>
            <input
              className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              type="password"
              value={geminiApiKey}
              onChange={(event) => setGeminiApiKey(event.target.value)}
              placeholder="環境変数でも設定できます"
            />
            <p className="mt-2 text-xs text-[#5c5d63]">
              登録状況: {geminiHasKey ? "登録済み" : "未登録"}
            </p>
          </div>
        </div>
        {settingError ? (
          <p className="mt-3 text-xs text-red-600">{settingError}</p>
        ) : null}
        {settingMessage ? (
          <p className="mt-3 text-xs text-[#5c5d63]">{settingMessage}</p>
        ) : null}
        <button
          className="mt-4 rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white"
          type="button"
          onClick={handleSave}
        >
          変更を保存
        </button>
      </section>
    </div>
  );
}
