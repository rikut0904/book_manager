import { getAuthState } from "./auth";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

type FetchOptions = RequestInit & {
  auth?: boolean;
};

export async function fetchJSON<T>(
  path: string,
  options: FetchOptions = {}
): Promise<T> {
  const headers = new Headers(options.headers);
  headers.set("Content-Type", "application/json");

  if (options.auth) {
    const auth = getAuthState();
    if (auth?.userId) {
      headers.set("X-User-Id", auth.userId);
    }
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || "Request failed");
  }

  return (await response.json()) as T;
}
