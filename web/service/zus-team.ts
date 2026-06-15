import { del, get, patch, post } from "./base"

export type Team = {
  id: string
  name: string
  description: string
  created_at: string
  updated_at: string
}

type ApiResponse<T> = { message: string; data: T }

export const listTeams = (page?: number, page_size?: number) => {
  return get<ApiResponse<{ count: number; teams: Team[] }>>(`/api/v1/teams`, {
    params: { page, page_size },
  })
}

export const createTeam = (name: string, description?: string) => {
  return post<ApiResponse<Team>>(`/api/v1/team`, {
    body: { name, description },
  })
}

export const updateTeam = (
  team_id: string,
  name?: string,
  description?: string
) => {
  return patch<{ message: string }>(`/api/v1/team/${team_id}`, {
    body: { name, description },
  })
}

export const deleteTeam = (team_id: string) => {
  return del<{ message: string }>(`/api/v1/team/${team_id}`)
}

export type TeamMember = {
  id: string
  team_id: string
  user_id: string
  email: string
  joined_at: string
}

export const listMembers = (team_id: string) => {
  return get<ApiResponse<{ count: number; members: TeamMember[] }>>(
    `/api/v1/team/${team_id}/members`
  )
}

export const addMember = (team_id: string, user_id: string) => {
  return post<{ message: string }>(`/api/v1/team/${team_id}/member`, {
    body: { user_id },
  })
}

export const removeMember = (team_id: string, member_id: string) => {
  return del<{ message: string }>(
    `/api/v1/team/${team_id}/member/${member_id}`
  )
}