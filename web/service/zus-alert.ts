import { get } from "./base"

export type AlertSeverity = "critical" | "high" | "medium" | "low" | "informational"
export type AlertStatus = "firing" | "resolved"

export type Alert = {
  id: string
  fingerprint: string
  source: string
  connector_id: string
  team_id: string
  severity: AlertSeverity
  status: AlertStatus
  starts_at: string
  ends_at: string
  entity_name: string
  entity_ip: string
  entity_svc: string
  labels: Record<string, string>
  annotations: Record<string, string>
  created_at: string
  updated_at: string
}

type ApiResponse<T> = { message: string; data: T }

export type ListAlertsParams = {
  page?: number
  page_size?: number
  status?: AlertStatus
  severity?: AlertSeverity
  connector_id?: string
  start_at?: string
  end_at?: string
}

export const listAlerts = (team_id: string, params: ListAlertsParams) => {
  return get<ApiResponse<{ count: number; alerts: Alert[] }>>(
    `/api/v1/team/${team_id}/alerts`,
    { params }
  )
}

export const getAlertHistory = (
  team_id: string,
  alert_id: string,
  page?: number,
  page_size?: number
) => {
  return get<ApiResponse<{ count: number; alerts: Alert[] }>>(
    `/api/v1/team/${team_id}/alert/${alert_id}/history`,
    { params: { page, page_size } }
  )
}