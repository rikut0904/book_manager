# ER（関係）概要

## 関係一覧
- users 1 - 1 profile_settings
- users 1 - n user_books
- users 1 - n recommendations
- users 1 - n favorites
- users 1 - n next_to_buy_manual
- users n - n users (follows)

- books 1 - n user_books
- books 1 - 1 book_series_auto
- books 1 - n recommendations

- series 1 - n book_series_auto
- series 1 - n user_book_series_override


## シリーズ解決ロジック
- series_id = coalesce(user_book_series_override.series_id, book_series_auto.series_id)
- volume_number = coalesce(user_book_series_override.volume_number, book_series_auto.volume_number)
