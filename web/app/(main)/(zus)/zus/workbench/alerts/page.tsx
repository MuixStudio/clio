"use client"

import { useEffect, useMemo, useState } from "react"

import {
  HeaderBar,
  HeaderBarDescription,
  HeaderBarLeft,
  HeaderBarTitle,
} from "@/app/(main)/components/header-bar/header-bar"
import { Accordion } from "@/components/ui/accordion"
import { Input } from "@/components/ui/input"
import { useZusTeam } from "@/hooks/use-zus-team"
import type { AlertStatus } from "@/service/zus-alert"
import { listConnectors, type Connector } from "@/service/zus-connector"

import { AlertFilters } from "./components/alert-filters"
import { AlertSeveritySection } from "./components/alert-severity-section"
import {
  CONNECTOR_FILTER_ALL,
  DEFAULT_ALERT_FILTERS,
  SEVERITY_SECTIONS,
  STATUS_FILTER_ALL,
} from "./components/constants"
import type { AlertFiltersValue } from "./components/types"
import { toEndOfDayISOString, toStartOfDayISOString } from "./components/utils"

export default function Page() {
  const { team } = useZusTeam()
  const [connectors, setConnectors] = useState<Connector[]>([])
  const [filters, setFilters] = useState<AlertFiltersValue>(DEFAULT_ALERT_FILTERS)
  const [draftFilters, setDraftFilters] = useState<AlertFiltersValue>(
    DEFAULT_ALERT_FILTERS
  )
  const [search, setSearch] = useState("")
  const [openSections, setOpenSections] = useState<string[]>(() =>
    SEVERITY_SECTIONS.filter((section) => section.defaultOpen).map(
      (section) => section.severity
    )
  )

  useEffect(() => {
    if (!team) return
    let cancelled = false

    listConnectors(team.id, 1, 100)
      .then(({ data }) => {
        if (!cancelled) setConnectors(data.connectors)
      })
      .catch((err) => console.error(err))

    return () => {
      cancelled = true
    }
  }, [team])

  const connectorNames = useMemo(
    () => Object.fromEntries(connectors.map((c) => [c.id, c.name])),
    [connectors]
  )

  const sectionFilters = useMemo(
    () => ({
      status:
        filters.status === STATUS_FILTER_ALL
          ? undefined
          : (filters.status as AlertStatus),
      connector_id:
        filters.connectorId === CONNECTOR_FILTER_ALL
          ? undefined
          : filters.connectorId,
      start_at: toStartOfDayISOString(filters.dateRange?.from),
      end_at: toEndOfDayISOString(filters.dateRange?.to),
    }),
    [filters]
  )

  const hasActiveFilters =
    filters.status !== STATUS_FILTER_ALL ||
    filters.connectorId !== CONNECTOR_FILTER_ALL ||
    Boolean(filters.dateRange?.from || filters.dateRange?.to)

  const applyFilters = () => setFilters(draftFilters)

  const clearFilters = () => {
    setDraftFilters(DEFAULT_ALERT_FILTERS)
    setFilters(DEFAULT_ALERT_FILTERS)
  }

  return (
    <>
      {/*header*/}
      <HeaderBar>
        <HeaderBarLeft>
          <HeaderBarTitle>Alerts</HeaderBarTitle>
          <HeaderBarDescription>View and filter all alerts</HeaderBarDescription>
        </HeaderBarLeft>
      </HeaderBar>

      {/*alerts*/}
      <div className="flex flex-col gap-3 p-4">
        <div className="flex flex-wrap items-center gap-2">
          <Input
            placeholder="Search by entity, IP, source..."
            value={search}
            onChange={(event) => setSearch(event.target.value)}
            className="h-8 w-60"
          />
          <AlertFilters
            value={draftFilters}
            connectors={connectors}
            hasActiveFilters={hasActiveFilters}
            onChange={setDraftFilters}
            onApply={applyFilters}
            onClear={clearFilters}
          />
        </div>

        {team && (
          <Accordion
            type="multiple"
            value={openSections}
            onValueChange={setOpenSections}
            className="flex flex-col gap-3"
          >
            {SEVERITY_SECTIONS.map((section) => (
              <AlertSeveritySection
                key={section.severity}
                teamId={team.id}
                severity={section.severity}
                code={section.code}
                label={section.label}
                dotClassName={section.dotClassName}
                headerClassName={section.headerClassName}
                open={openSections.includes(section.severity)}
                filters={sectionFilters}
                search={search}
                connectorNames={connectorNames}
              />
            ))}
          </Accordion>
        )}
      </div>
    </>
  )
}