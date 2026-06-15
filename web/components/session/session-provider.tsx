"use client"

import { createContext, useCallback, useEffect, useState } from "react"
import type { Session } from "./session-types"

const ACCESS_TOKEN_KEY = "access_token"

type SessionState =
  | { state: "authenticated"; session: Session }
  | { state: "unauthenticated" }
  | { state: "error"; error: Error }

export type SessionContextData = {
  isLoading: boolean
  initialized: boolean
  session: Session | null
  error: Error | undefined
  refetch: () => void
}

export const SessionContext = createContext<SessionContextData>({
  session: null,
  isLoading: false,
  initialized: false,
  error: undefined,
  refetch: () => {},
})

function decodeJwt<T>(token: string): T {
  const base64Url = token.split(".")[1]
  if (!base64Url) throw new Error("Invalid JWT: missing payload")
  const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/")
  const json = decodeURIComponent(
    atob(base64)
      .split("")
      .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
      .join(""),
  )
  return JSON.parse(json) as T
}

function readSessionFromStorage(): Session | null {
  const token = localStorage.getItem(ACCESS_TOKEN_KEY)
  if (!token) return null
  const session = decodeJwt<Session>(token)
  if (!session.active) return null
  if (new Date(session.expires_at) < new Date()) return null
  return session
}

function resolveState(session: Session | null): SessionState {
  return session
    ? { state: "authenticated", session }
    : { state: "unauthenticated" }
}

/**
 * Reads the access_token JWT from localStorage, decodes the session payload,
 * and provides it via context. Syncs automatically across browser tabs via
 * the storage event.
 *
 * @example
 * <SessionProvider>
 *   <App />
 * </SessionProvider>
 */
export function SessionProvider({ children }: React.PropsWithChildren) {
  const [sessionState, setSessionState] = useState<SessionState>(() => {
    try {
      return resolveState(readSessionFromStorage())
    } catch (error) {
      return { state: "error", error: error as Error }
    }
  })

  const refetch = useCallback(() => {
    try {
      setSessionState(resolveState(readSessionFromStorage()))
    } catch (error) {
      setSessionState({ state: "error", error: error as Error })
    }
  }, [])

  // sync session state when another tab writes/removes the token
  useEffect(() => {
    const handleStorage = (e: StorageEvent) => {
      if (e.key === ACCESS_TOKEN_KEY) refetch()
    }
    window.addEventListener("storage", handleStorage)
    return () => window.removeEventListener("storage", handleStorage)
  }, [refetch])

  return (
    <SessionContext.Provider
      value={{
        session:
          sessionState.state === "authenticated" ? sessionState.session : null,
        error:
          sessionState.state === "error" ? sessionState.error : undefined,
        isLoading: false,
        initialized: true,
        refetch,
      }}
    >
      {children}
    </SessionContext.Provider>
  )
}