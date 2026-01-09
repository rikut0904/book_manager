"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { fetchJSON } from "@/lib/api";
import { getAuthState } from "@/lib/auth";

type SettingsResponse = {
  isAdmin?: boolean;
  settings?: {
    openaiModel?: string;
    openaiHasKey?: boolean;
  };
};

type OpenAIKeyItem = {
  id: string;
  name: string;
  maskedKey: string;
  createdAt: string;
};

type OpenAIModelResponse = {
  items: string[];
};

const defaultModel = "gpt-4o-mini";

export default function AISettingsPage() {
  const [openaiModel, setOpenAIModel] = useState(defaultModel);
  const [settingError, setSettingError] = useState<string | null>(null);
  const [settingMessage, setSettingMessage] = useState<string | null>(null);
  const [isAdmin, setIsAdmin] = useState(false);
  const [keys, setKeys] = useState<OpenAIKeyItem[]>([]);
  const [keyForm, setKeyForm] = useState({ name: "", apiKey: "" });
  const [models, setModels] = useState<string[]>([]);

  useEffect(() => {
    const auth = getAuthState();
    if (!auth?.userId) {
      return;
    }
    fetchJSON<SettingsResponse>(`/users/${auth.userId}`, { auth: true })
      .then((data) => {
        const settings = data.settings ?? {};
        setIsAdmin(Boolean(data.isAdmin));
        setOpenAIModel(settings.openaiModel || defaultModel);
      })
      .catch(() => {
        setSettingError("AI設定を取得できませんでした。");
      });
  }, []);

  useEffect(() => {
    if (!isAdmin) {
      return;
    }
    let isMounted = true;
    fetchJSON<{ items: OpenAIKeyItem[] }>("/admin/openai-keys", { auth: true })
      .then((data) => {
        if (!isMounted) {
          return;
        }
        setKeys(data.items ?? []);
      })
      .catch(() => {
        if (!isMounted) {
          return;
        }
        setKeys([]);
      });
    return () => {
      isMounted = false;
    };
  }, [isAdmin]);

  useEffect(() => {
    if (!isAdmin) {
      return;
    }
    let isMounted = true;
    fetchJSON<OpenAIModelResponse>("/admin/openai-models", { auth: true })
      .then((data) => {
        if (!isMounted) {
          return;
        }
        setModels(data.items ?? []);
      })
      .catch(() => {
        if (!isMounted) {
          return;
        }
        setModels([]);
      });
    return () => {
      isMounted = false;
    };
  }, [isAdmin]);

  const resolvedModel =
    models.length > 0 && !models.includes(openaiModel)
      ? models[0]
      : openaiModel;

  const handleSave = async () => {
    setSettingError(null);
    setSettingMessage(null);
    try {
      await fetchJSON("/users/me/settings", {
        method: "PATCH",
        auth: true,
        body: JSON.stringify({
          openaiModel: resolvedModel.trim(),
        }),
      });
      setSettingMessage("AI設定を保存しました。");
    } catch {
      setSettingError("AI設定の保存に失敗しました。");
    }
  };

  const handleCreateKey = async () => {
    setSettingError(null);
    if (!keyForm.name.trim() || !keyForm.apiKey.trim()) {
      setSettingError("キー名とAPIキーを入力してください。");
      return;
    }
    try {
      const item = await fetchJSON<OpenAIKeyItem>("/admin/openai-keys", {
        method: "POST",
        auth: true,
        body: JSON.stringify({
          name: keyForm.name.trim(),
          apiKey: keyForm.apiKey.trim(),
        }),
      });
      setKeys((prev) => [...prev, item]);
      setKeyForm({ name: "", apiKey: "" });
    } catch {
      setSettingError("APIキーの登録に失敗しました。");
    }
  };

  const handleDeleteKey = async (id: string) => {
    setSettingError(null);
    try {
      await fetchJSON(`/admin/openai-keys/${id}`, {
        method: "DELETE",
        auth: true,
      });
      setKeys((prev) => prev.filter((item) => item.id !== id));
    } catch {
      setSettingError("APIキーの削除に失敗しました。");
    }
  };

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <Link className="text-xs text-[#5c5d63]" href="/settings">
          ← 設定に戻る
        </Link>
        <h1 className="mt-3 font-[var(--font-display)] text-3xl">
          OpenAI設定
        </h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          OpenAI APIを利用したシリーズ推定の設定を行います。
        </p>
      </section>

      <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
        {!isAdmin ? (
          <p className="text-sm text-[#5c5d63]">
            この設定は管理者ユーザーのみ利用できます。
          </p>
        ) : null}
        <div className="flex flex-col gap-4 text-sm">
          <div className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3">
            <p className="text-xs text-[#5c5d63]">モデル名（共有キーで利用）</p>
            <select
              className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              value={resolvedModel}
              onChange={(event) => setOpenAIModel(event.target.value)}
              disabled={!isAdmin}
            >
              <option value={defaultModel}>{defaultModel}</option>
              {models.map((model) => (
                <option key={model} value={model}>
                  {model}
                </option>
              ))}
            </select>
            <p className="mt-2 text-xs text-[#5c5d63]">
              モデル一覧は共有APIキーで自動取得します。
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
          disabled={!isAdmin}
        >
          変更を保存
        </button>
      </section>

      {isAdmin ? (
        <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
          <h2 className="font-[var(--font-display)] text-2xl">
            共有OpenAI APIキー
          </h2>
          <p className="mt-2 text-sm text-[#5c5d63]">
            登録されたキーは全ユーザーで共通利用されます。表示は一部マスクされます。
          </p>
          <div className="mt-4 rounded-2xl border border-dashed border-[#e4d8c7] bg-white/80 p-4 text-sm text-[#5c5d63]">
            <p className="font-medium text-[#1b1c1f]">登録手順</p>
            <ol className="mt-2 list-decimal pl-5 text-xs leading-6">
              <li>キー名を入力（例: メインキー）</li>
              <li>APIキーを貼り付け</li>
              <li>追加を押す</li>
            </ol>
          </div>
          <div className="mt-4 flex flex-wrap gap-2 text-xs">
            {keys.map((item) => (
              <span
                key={item.id}
                className="flex items-center gap-2 rounded-full border border-[#e4d8c7] bg-white px-3 py-2 text-[#5c5d63]"
              >
                {item.name}: {item.maskedKey}
                <button
                  className="text-[10px] text-[#c86b3c] hover:text-[#8f3d1f]"
                  type="button"
                  onClick={() => handleDeleteKey(item.id)}
                >
                  削除
                </button>
              </span>
            ))}
          </div>
          <div className="mt-4 grid gap-3 text-sm md:grid-cols-[1fr_1fr_auto]">
            <input
              className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              placeholder="キー名"
              value={keyForm.name}
              onChange={(event) =>
                setKeyForm((prev) => ({ ...prev, name: event.target.value }))
              }
            />
            <input
              className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              placeholder="APIキー"
              type="password"
              value={keyForm.apiKey}
              onChange={(event) =>
                setKeyForm((prev) => ({ ...prev, apiKey: event.target.value }))
              }
            />
            <button
              className="rounded-full bg-[#efe5d4] px-4 py-3 text-sm text-[#1b1c1f]"
              type="button"
              onClick={handleCreateKey}
            >
              追加
            </button>
          </div>
        </section>
      ) : null}
    </div>
  );
}
