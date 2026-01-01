"use client";

import { useState } from "react";

import { fetchJSON } from "@/lib/api";

type LookupResult = {
  id: string;
  isbn13: string;
  title: string;
  authors: string[];
  publisher: string;
  publishedDate: string;
  thumbnailUrl: string;
  source: string;
};

export default function BookNewPage() {
  const [isbn, setIsbn] = useState("");
  const [result, setResult] = useState<LookupResult | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [manualError, setManualError] = useState<string | null>(null);
  const [manualSuccess, setManualSuccess] = useState<string | null>(null);
  const [manualForm, setManualForm] = useState({
    title: "",
    authors: "",
    isbn13: "",
  });

  const handleLookup = async () => {
    setError(null);
    setResult(null);
    if (!isbn.trim()) {
      setError("ISBNを入力してください。");
      return;
    }
    try {
      const data = await fetchJSON<LookupResult>(
        `/isbn/lookup?isbn=${encodeURIComponent(isbn.trim())}`,
        { auth: true }
      );
      setResult(data);
    } catch {
      setError("ISBN検索に失敗しました。");
    }
  };

  const handleManualCreate = async () => {
    setManualError(null);
    setManualSuccess(null);
    if (!manualForm.title.trim()) {
      setManualError("タイトルを入力してください。");
      return;
    }
    try {
      const data = await fetchJSON<LookupResult>("/books", {
        method: "POST",
        auth: true,
        body: JSON.stringify({
          title: manualForm.title.trim(),
          authors: manualForm.authors
            .split(",")
            .map((item) => item.trim())
            .filter(Boolean),
          isbn13: manualForm.isbn13.trim(),
        }),
      });
      setManualSuccess(`登録しました: ${data.title}`);
      setManualForm({ title: "", authors: "", isbn13: "" });
    } catch {
      setManualError("手動登録に失敗しました。");
    }
  };

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
          Register
        </p>
        <h1 className="mt-2 font-[var(--font-display)] text-3xl">
          書籍登録
        </h1>
        <p className="mt-2 text-sm text-[#5c5d63]">
          ISBN検索で書誌を自動取得。見つからない場合は手動で登録できます。
        </p>
      </section>

      <div className="grid gap-6 lg:grid-cols-[1.2fr_0.8fr]">
        <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
          <h2 className="font-[var(--font-display)] text-2xl">ISBN検索</h2>
          <p className="mt-2 text-sm text-[#5c5d63]">
            13桁のISBNを入力してください。バーコード入力は後から追加予定です。
          </p>
          <div className="mt-4 flex flex-wrap gap-3">
            <input
              className="flex-1 rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
              placeholder="978-4-00-123456-7"
              value={isbn}
              onChange={(event) => setIsbn(event.target.value)}
            />
            <button
              className="rounded-full bg-[#1b1c1f] px-5 py-3 text-sm font-medium text-white"
              type="button"
              onClick={handleLookup}
            >
              取得
            </button>
          </div>
          {error ? <p className="mt-3 text-sm text-red-600">{error}</p> : null}
          <div className="mt-6 rounded-2xl border border-dashed border-[#e4d8c7] bg-white/70 p-4 text-sm text-[#5c5d63]">
            <p className="font-medium text-[#1b1c1f]">取得結果</p>
            {result ? (
              <>
                <p className="mt-2">タイトル: {result.title}</p>
                <p>著者: {result.authors?.join(" / ") || "未登録"}</p>
                <p>
                  出版社: {result.publisher || "未登録"} /{" "}
                  {result.publishedDate || "未登録"}
                </p>
                <p className="mt-2 text-xs text-[#c86b3c]">
                  登録ID: {result.id}
                </p>
              </>
            ) : (
              <p className="mt-2">まだ結果がありません。</p>
            )}
          </div>
        </section>

        <section className="rounded-3xl border border-[#e4d8c7] bg-white/70 p-6 shadow-sm">
          <h2 className="font-[var(--font-display)] text-2xl">手動登録</h2>
          <div className="mt-4 flex flex-col gap-3 text-sm">
            <label className="text-[#1b1c1f]">
              タイトル
              <input
                className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
                value={manualForm.title}
                onChange={(event) =>
                  setManualForm((prev) => ({
                    ...prev,
                    title: event.target.value,
                  }))
                }
              />
            </label>
            <label className="text-[#1b1c1f]">
              著者
              <input
                className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
                placeholder="著者A, 著者B"
                value={manualForm.authors}
                onChange={(event) =>
                  setManualForm((prev) => ({
                    ...prev,
                    authors: event.target.value,
                  }))
                }
              />
            </label>
            <label className="text-[#1b1c1f]">
              ISBN (任意)
              <input
                className="mt-2 w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c]"
                value={manualForm.isbn13}
                onChange={(event) =>
                  setManualForm((prev) => ({
                    ...prev,
                    isbn13: event.target.value,
                  }))
                }
              />
            </label>
            {manualError ? (
              <p className="text-xs text-red-600">{manualError}</p>
            ) : null}
            {manualSuccess ? (
              <p className="text-xs text-emerald-700">{manualSuccess}</p>
            ) : null}
            <button
              className="rounded-full bg-[#c86b3c] px-5 py-3 text-sm font-medium text-white"
              type="button"
              onClick={handleManualCreate}
            >
              所蔵に追加
            </button>
          </div>
        </section>
      </div>
    </div>
  );
}
