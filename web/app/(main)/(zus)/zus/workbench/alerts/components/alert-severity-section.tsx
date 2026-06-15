"use client"

import { useCallback, useEffect, useRef, useState } from "react"
import { Loader2, Siren } from "lucide-react"

import {
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { cn } from "@/lib/utils"
import {
  listAlerts,
  type Alert,
  type AlertSeverity,
  type AlertStatus,
} from "@/service/zus-alert"

import { AlertRow } from "./alert-row"

const PAGE_SIZE = 10

type SectionFilters = {
  status?: AlertStatus
  connector_id?: string
  start_at?: string
  end_at?: string
}

function SectionLabel({
  code,
  label,
  dotClassName,
}: {
  code: string
  label: string
  dotClassName: string
}) {
  return (
    <span className="inline-flex items-center gap-1.5 rounded-full border px-2 py-1 text-xs font-medium">
      <span className={cn("size-2 rounded-full", dotClassName)} />
      {code} {label}
    </span>
  )
}

export function AlertSeveritySection({
  teamId,
  severity,
  code,
  label,
  dotClassName,
  headerClassName,
  open,
  filters,
  search,
  connectorNames,
}: {
  teamId: string
  severity: AlertSeverity
  code: string
  label: string
  dotClassName: string
  headerClassName: string
  open: boolean
  filters: SectionFilters
  search: string
  connectorNames: Record<string, string>
}) {
  const [alerts, setAlerts] = useState<Alert[]>([])
  const [page, setPage] = useState(0)
  const [hasMore, setHasMore] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const requestIdRef = useRef(0)
  const containerRef = useRef<HTMLDivElement>(null)
  const sentinelRef = useRef<HTMLDivElement>(null)

  const fetchPage = useCallback(
    async (pageToLoad: number, replace: boolean) => {
      const requestId = ++requestIdRef.current
      setLoading(true)
      try {
        const { data } = await listAlerts(teamId, {
          page: pageToLoad,
          page_size: PAGE_SIZE,
          severity,
          ...filters,
        })
        if (requestIdRef.current !== requestId) return
        setError(null)
        setAlerts((prev) =>
          replace ? data.alerts : [...prev, ...data.alerts]
        )
        setHasMore(data.alerts.length === PAGE_SIZE)
        setPage(pageToLoad)
      } catch (err) {
        if (requestIdRef.current !== requestId) return
        console.error(err)
        setError("Failed to load alerts")
        if (replace) {
          setAlerts([])
          setHasMore(false)
        }
      } finally {
        if (requestIdRef.current === requestId) setLoading(false)
      }
    },
    [teamId, severity, filters]
  )

  useEffect(() => {
    void (async () => {
      await fetchPage(1, true)
    })()
  }, [fetchPage])

  useEffect(() => {
    if (!open || !hasMore || loading) return
    const container = containerRef.current
    const sentinel = sentinelRef.current
    if (!container || !sentinel) return

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0]?.isIntersecting) void fetchPage(page + 1, false)
      },
      { root: container, threshold: 0 }
    )
    observer.observe(sentinel)
    return () => observer.disconnect()
  }, [open, hasMore, loading, page, fetchPage])

  const query = search.trim().toLowerCase()
  const filteredAlerts = query
    ? alerts.filter((alert) =>
        [
          alert.entity_name,
          alert.entity_ip,
          alert.entity_svc,
          alert.source,
          alert.fingerprint,
        ].some((value) => value?.toLowerCase().includes(query))
      )
    : alerts

  return (
    <AccordionItem value={severity} className="rounded-lg border">
      <AccordionTrigger
        className={cn(
          "items-center px-3 hover:no-underline data-[state=open]:rounded-b-none",
          headerClassName
        )}
      >
        <SectionLabel code={code} label={label} dotClassName={dotClassName} />
      </AccordionTrigger>
      <AccordionContent className="h-auto border-t border-dashed pb-0">
        <div ref={containerRef} className="max-h-160 overflow-y-auto">
          {loading && alerts.length === 0 ? (
            <div className="flex flex-col gap-2 p-3">
              {Array.from({ length: 3 }).map((_, index) => (
                <Skeleton key={index} className="h-12 w-full" />
              ))}
            </div>
          ) : error ? (
            <div className="flex flex-col items-center gap-2 p-6 text-center text-sm text-muted-foreground">
              <Siren className="size-5" />
              {error}
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => void fetchPage(1, true)}
              >
                Retry
              </Button>
            </div>
          ) : filteredAlerts.length === 0 ? (
            <div className="p-6 text-center text-sm text-muted-foreground">
              No alerts at this severity
            </div>
          ) : (
            <div className="divide-y">
              {filteredAlerts.map((alert) => (
                <AlertRow
                  key={alert.id}
                  alert={alert}
                  connectorName={connectorNames[alert.connector_id]}
                />
              ))}
              {hasMore && (
                <div
                  ref={sentinelRef}
                  className="flex items-center justify-center p-2"
                >
                  {loading && (
                    <Loader2 className="size-4 animate-spin text-muted-foreground" />
                  )}
                </div>
              )}
            </div>
          )}
        </div>
      </AccordionContent>
    </AccordionItem>
  )
}