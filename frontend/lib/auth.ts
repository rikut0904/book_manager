export type AuthState = {
  accessToken: string;
  refreshToken: string;
  userId: string;
  displayName?: string;
};

const STORAGE_KEY = "book_manager_auth";

export function getAuthState(): AuthState | null {
  if (typeof window === "undefined") {
    return null;
  }
  const raw = window.localStorage.getItem(STORAGE_KEY);
  if (!raw) {
    return null;
  }
  try {
    return JSON.parse(raw) as AuthState;
  } catch {
    return null;
  }
}

export function setAuthState(state: AuthState) {
  if (typeof window === "undefined") {
    return;
  }
  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
  window.dispatchEvent(new Event("auth-changed"));
}

export function updateAuthState(patch: Partial<AuthState>) {
  if (typeof window === "undefined") {
    return;
  }
  const current = getAuthState();
  if (!current) {
    return;
  }
  setAuthState({ ...current, ...patch });
}

export function clearAuthState() {
  if (typeof window === "undefined") {
    return;
  }
  window.localStorage.removeItem(STORAGE_KEY);
  window.dispatchEvent(new Event("auth-changed"));
}
