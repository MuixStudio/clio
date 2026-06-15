"use client"

import { useContext } from "react"
import { SessionContext } from "./session-provider"

export function useSession() {
  return useContext(SessionContext)
}