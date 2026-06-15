"use client"

import { useCallback, useEffect, useState } from "react"

import { listTeams, type Team } from "@/service/zus-team"

const TEAM_ID_KEY = "zus_team_id"
const TEAM_CHANGE_EVENT = "zus-team-change"

function readStoredTeamId(): string | null {
  return localStorage.getItem(TEAM_ID_KEY)
}

function writeStoredTeamId(teamId: string) {
  localStorage.setItem(TEAM_ID_KEY, teamId)
  window.dispatchEvent(new CustomEvent(TEAM_CHANGE_EVENT, { detail: teamId }))
}

/**
 * Fetches the current user's teams and tracks which one is selected.
 * The selection is persisted to localStorage and kept in sync across
 * every component using this hook (and across browser tabs).
 */
export function useZusTeam() {
  const [teams, setTeams] = useState<Team[]>([])
  const [teamId, setTeamId] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    let cancelled = false

    setIsLoading(true)
    listTeams()
      .then(({ data }) => {
        if (cancelled) return
        setTeams(data.teams)
        setTeamId((current) => {
          if (current && data.teams.some((t) => t.id === current)) return current
          const stored = readStoredTeamId()
          if (stored && data.teams.some((t) => t.id === stored)) return stored
          return data.teams[0]?.id ?? null
        })
      })
      .finally(() => {
        if (!cancelled) setIsLoading(false)
      })

    return () => {
      cancelled = true
    }
  }, [])

  useEffect(() => {
    const onTeamChange = (e: Event) => {
      const id = (e as CustomEvent<string>).detail ?? readStoredTeamId()
      if (id) setTeamId(id)
    }
    const onStorage = (e: StorageEvent) => {
      if (e.key === TEAM_ID_KEY && e.newValue) setTeamId(e.newValue)
    }

    window.addEventListener(TEAM_CHANGE_EVENT, onTeamChange)
    window.addEventListener("storage", onStorage)
    return () => {
      window.removeEventListener(TEAM_CHANGE_EVENT, onTeamChange)
      window.removeEventListener("storage", onStorage)
    }
  }, [])

  const selectTeam = useCallback((id: string) => {
    setTeamId(id)
    writeStoredTeamId(id)
  }, [])

  const team = teams.find((t) => t.id === teamId) ?? null

  return { team, teams, isLoading, selectTeam }
}