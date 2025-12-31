# APIサンプル（Request/Response）

## サインアップ
POST /auth/signup

req:
```
{
  "email": "user@example.com",
  "password": "********",
  "username": "rikut"
}
```

res:
```
{
  "accessToken": "...",
  "refreshToken": "...",
  "user": {"id": "...", "email": "...", "username": "rikut"}
}
```

## ISBN検索
GET /isbn/lookup?isbn=978...

res:
```
{
  "isbn13": "978...",
  "title": "...",
  "authors": ["..."],
  "publisher": "...",
  "publishedDate": "2024-01-01",
  "thumbnailUrl": "...",
  "source": "google"
}
```

## 所蔵登録
POST /user-books

req:
```
{
  "bookId": "...",
  "note": "メモ",
  "acquiredAt": "2024-01-01"
}
```

res:
```
{
  "id": "...",
  "bookId": "...",
  "note": "メモ",
  "acquiredAt": "2024-01-01"
}
```

## シリーズ上書き
PATCH /user-series/override

req:
```
{
  "bookId": "...",
  "seriesId": "...",
  "volumeNumber": 3
}
```

res:
```
{"ok": true}
```

## お気に入り登録
POST /favorites

req:
```
{"type": "series", "seriesId": "..."}
```

res:
```
{"id": "...", "type": "series", "seriesId": "..."}
```

## 次に買う本（手入力）
POST /next-to-buy/manual

req:
```
{
  "title": "...",
  "seriesName": "...",
  "volumeNumber": 5,
  "note": "..."
}
```

res:
```
{"id": "..."}
```

## おすすめ投稿
POST /recommendations

req:
```
{"bookId": "...", "comment": "面白かった"}
```

res:
```
{"id": "..."}
```

## プロフィール取得
GET /users/{id}

res:
```
{
  "user": {"id": "...", "username": "..."},
  "stats": {"ownedCount": 10, "seriesCount": 3, "followers": 2, "following": 5},
  "recommendations": ["..."]
}
```

## 書誌情報の報告
POST /book-reports

req:
```
{
  "bookId": "...",
  "suggestion": "正しいタイトル...",
  "note": "備考"
}
```

res:
```
{"ok": true}
```
