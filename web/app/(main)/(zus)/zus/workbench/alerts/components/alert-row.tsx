"use client"

import { useRouter } from "next/navigation"

import type { Alert } from "@/service/zus-alert"

import { AlertStatusBadge } from "./status-badge"
import { formatRelativeTime } from "./utils"

export function AlertRow({
  alert,
  connectorName,
}: {
  alert: Alert
  connectorName?: string
}) {
  const router = useRouter()
  const goToHistory = () => router.push(`/zus/alert/history?alert_id=${alert.id}`)

  const subtitle = [
    alert.entity_svc,
    alert.entity_ip,
    alert.source && `Source ${alert.source}`,
    connectorName,
    formatRelativeTime(alert.starts_at),
  ]
    .filter(Boolean)
    .join(" · ")

  return (
    <div
      role="button"
      tabIndex={0}
      onClick={goToHistory}
      onKeyDown={(event) => {
        if (event.key === "Enter") goToHistory()
      }}
      className="flex cursor-pointer items-center gap-3 px-3 py-2.5 hover:bg-muted/50"
    >
      <AlertStatusBadge status={alert.status} />
      <div className="min-w-0 flex-1">
        <div className="truncate text-sm font-medium">
          {alert.entity_name || alert.fingerprint}
        </div>
        {subtitle && (
          <div className="truncate text-xs text-muted-foreground">
            {subtitle}
          </div>
        )}
      </div>
    </div>
  )
}