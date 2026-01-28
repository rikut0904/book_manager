import { getAuthState } from "./auth";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

type FetchOptions = RequestInit & {
  auth?: boolean;
};

export class APIError extends Error {
  status: number;
  data: Record<string, unknown>;

  constructor(message: string, status: number, data: Record<string, unknown> = {}) {
    super(message);
    this.name = "APIError";
    this.status = status;
    this.data = data;
  }
}

export async function fetchJSON<T>(
  path: string,
  options: FetchOptions = {}
): Promise<T> {
  const headers = new Headers(options.headers);
  headers.set("Content-Type", "application/json");

  if (options.auth) {
    const auth = getAuthState();
    if (auth?.accessToken) {
      headers.set("Authorization", `Bearer ${auth.accessToken}`);
    }
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const text = await response.text();
    let message = "Request failed";
    let data: Record<string, unknown> = {};
    try {
      const json = JSON.parse(text);
      message = json.message || json.error || message;
      data = json;
    } catch {
      message = text || message;
    }
    throw new APIError(message, response.status, data);
  }

  return (await response.json()) as T;
}
