import type { AlertSeverity } from "@/service/zus-alert"

import type { AlertFiltersValue } from "./types"

export const STATUS_FILTER_ALL = "all"
export const CONNECTOR_FILTER_ALL = "all"

export const STATUS_OPTIONS = [
  { label: "All Statuses", value: STATUS_FILTER_ALL },
  { label: "Firing", value: "firing" },
  { label: "Resolved", value: "resolved" },
] as const

export const DEFAULT_ALERT_FILTERS: AlertFiltersValue = {
  status: STATUS_FILTER_ALL,
  connectorId: CONNECTOR_FILTER_ALL,
  dateRange: undefined,
}

export const SEVERITY_SECTIONS: {
  severity: AlertSeverity
  code: string
  label: string
  defaultOpen: boolean
  dotClassName: string
  headerClassName: string
}[] = [
  {
    severity: "critical",
    code: "P0",
    label: "Critical",
    defaultOpen: true,
    dotClassName: "bg-destructive",
    headerClassName: "bg-destructive/10",
  },
  {
    severity: "high",
    code: "P1",
    label: "High",
    defaultOpen: true,
    dotClassName: "bg-orange-500",
    headerClassName: "bg-orange-500/10",
  },
  {
    severity: "medium",
    code: "P2",
    label: "Medium",
    defaultOpen: false,
    dotClassName: "bg-warning",
    headerClassName: "bg-warning/10",
  },
  {
    severity: "low",
    code: "P3",
    label: "Low",
    defaultOpen: false,
    dotClassName: "bg-info",
    headerClassName: "bg-info/10",
  },
  {
    severity: "informational",
    code: "P4",
    label: "Info",
    defaultOpen: false,
    dotClassName: "bg-muted-foreground",
    headerClassName: "bg-muted/50",
  },
]