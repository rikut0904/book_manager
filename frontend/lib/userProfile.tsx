"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";

import { fetchJSON } from "@/lib/api";
import { getAuthState, updateAuthState } from "@/lib/auth";

type UserProfile = {
  id: string;
  email: string;
  userId: string;
  displayName: string;
};

type UserProfileState = {
  profile: UserProfile | null;
  isLoading: boolean;
  error: string | null;
  refresh: () => Promise<void>;
};

const UserProfileContext = createContext<UserProfileState | null>(null);

export function UserProfileProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    const auth = getAuthState();
    if (!auth?.accessToken) {
      setProfile(null);
      setIsLoading(false);
      return;
    }
    setIsLoading(true);
    setError(null);
    try {
      const data = await fetchJSON<{ user: UserProfile }>("/users/profile", {
        auth: true,
      });
      setProfile(data.user ?? null);
      if (data.user?.displayName) {
        updateAuthState({ displayName: data.user.displayName });
      }
    } catch {
      setError("ユーザー情報を取得できませんでした。");
      setProfile(null);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  useEffect(() => {
    const handler = () => {
      refresh();
    };
    window.addEventListener("auth-changed", handler);
    return () => {
      window.removeEventListener("auth-changed", handler);
    };
  }, [refresh]);

  const value = useMemo(
    () => ({ profile, isLoading, error, refresh }),
    [profile, isLoading, error, refresh]
  );

  return (
    <UserProfileContext.Provider value={value}>
      {children}
    </UserProfileContext.Provider>
  );
}

export function useUserProfile() {
  const ctx = useContext(UserProfileContext);
  if (!ctx) {
    throw new Error("useUserProfile must be used within UserProfileProvider");
  }
  return ctx;
}
