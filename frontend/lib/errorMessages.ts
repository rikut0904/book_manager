export const commonErrorMessages: Record<string, string> = {
  email_required: "メールアドレスを入力してください。",
  invalid_email: "メールアドレスの形式が正しくありません。",
  email_exists: "このメールアドレスは既に登録されています。",
  display_name_too_long: "ユーザー名は50文字以内で入力してください。",
};

export const signupErrorMessages: Record<string, string> = {
  ...commonErrorMessages,
  password_required: "パスワードを入力してください。",
  password_too_short: "パスワードは8文字以上で入力してください。",
  user_id_required: "ユーザーIDを入力してください。",
  user_id_too_short: "ユーザーIDは2文字以上で入力してください。",
  user_id_too_long: "ユーザーIDは20文字以内で入力してください。",
  user_id_exists: "このユーザーIDは既に使用されています。",
};
