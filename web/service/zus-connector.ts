import { del, get, patch, post } from "./base"

export type Connector = {
  id: string
  type: string
  name: string
  token?: string
  labels: Record<string, string>
  enabled: boolean
  created_at: string
  updated_at: string
}

type ApiResponse<T> = { message: string; data: T }

export const listConnectors = (
  team_id: string,
  page?: number,
  page_size?: number,
  status?: "enable" | "disable",
  type?: string
) => {
  return get<ApiResponse<{ count: number; connectors: Connector[] }>>(
    `/api/v1/team/${team_id}/connectors`,
    { params: { page, page_size, status, type } }
  )
}

export const deleteConnector = (team_id: string, connector_id: string) => {
  return del<{ message: string }>(
    `/api/v1/team/${team_id}/connector/${connector_id}`
  )
}

export const createConnector = (
  name: string,
  type: string,
  labels: object,
  team_id: string
) => {
  return post<ApiResponse<Connector>>(`/api/v1/team/${team_id}/connector`, {
    body: { name, type, labels: labels },
  })
}

export const countConnectors = (
  team_id: string,
  type: string,
  status?: "enable" | "disable"
) => {
  return get<ApiResponse<{ type: string; count: number }>>(
    `/api/v1/team/${team_id}/connector/count`,
    {
      params: status ? { type, status } : { type },
    }
  )
}

export const enableConnector = (team_id: string, connector_id: string) => {
  return patch<{ message: string }>(
    `/api/v1/team/${team_id}/connector/${connector_id}`,
    {
      body: {
        enabled: true,
      },
    }
  )
}

export const disableConnector = (team_id: string, connector_id: string) => {
  return patch<{ message: string }>(
    `/api/v1/team/${team_id}/connector/${connector_id}`,
    {
      body: {
        enabled: false,
      },
    }
  )
}
