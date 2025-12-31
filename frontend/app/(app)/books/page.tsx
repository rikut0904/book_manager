import Link from "next/link";

const books = [
  {
    id: "bk_001",
    title: "透明な約束",
    author: "冬原 千景",
    series: "星街メトロ",
    volume: 4,
    note: "サイン本",
  },
  {
    id: "bk_002",
    title: "海辺の標本室",
    author: "浅井 まどか",
    series: "海辺の研究録",
    volume: 1,
    note: "装丁が好き",
  },
  {
    id: "bk_003",
    title: "火曜の地図帳",
    author: "鐘ヶ江 透",
    series: "徒歩旅行記",
    volume: 2,
    note: "次巻予約済み",
  },
];

export default function BooksPage() {
  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-3xl border border-[#e4d8c7] bg-white/80 p-6 shadow-sm">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <p className="text-xs uppercase tracking-[0.3em] text-[#c86b3c]">
              Library
            </p>
            <h1 className="mt-2 font-[var(--font-display)] text-3xl">
              所蔵一覧
            </h1>
          </div>
          <div className="flex flex-wrap gap-3">
            <button className="rounded-full border border-[#e4d8c7] px-4 py-2 text-xs text-[#5c5d63]">
              タグ
            </button>
            <button className="rounded-full border border-[#e4d8c7] px-4 py-2 text-xs text-[#5c5d63]">
              シリーズ
            </button>
            <Link
              className="rounded-full bg-[#1b1c1f] px-4 py-2 text-xs font-medium text-white"
              href="/books/new"
            >
              書籍登録
            </Link>
          </div>
        </div>
        <div className="mt-6 flex flex-wrap gap-3">
          <input
            className="w-full rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm outline-none transition focus:border-[#c86b3c] md:flex-1"
            placeholder="タイトル・著者・ISBNで検索"
          />
          <select className="rounded-2xl border border-[#e4d8c7] bg-white px-4 py-3 text-sm text-[#5c5d63]">
            <option>全て</option>
            <option>お気に入り</option>
            <option>最近追加</option>
          </select>
        </div>
      </section>

      <section className="grid gap-4 lg:grid-cols-3">
        {books.map((book) => (
          <Link
            key={book.id}
            className="group rounded-3xl border border-[#e4d8c7] bg-white/70 p-5 shadow-sm transition hover:-translate-y-1 hover:shadow-md"
            href={`/books/${book.id}`}
          >
            <div className="flex items-center justify-between text-xs text-[#5c5d63]">
              <span>{book.series}</span>
              <span>Vol.{book.volume}</span>
            </div>
            <h2 className="mt-4 font-[var(--font-display)] text-xl text-[#1b1c1f]">
              {book.title}
            </h2>
            <p className="mt-2 text-sm text-[#5c5d63]">{book.author}</p>
            <div className="mt-4 rounded-2xl bg-[#f6f1e7] px-3 py-2 text-xs text-[#5c5d63]">
              メモ: {book.note}
            </div>
            <p className="mt-4 text-xs text-[#c86b3c]">詳細を見る →</p>
          </Link>
        ))}
      </section>
    </div>
  );
}
