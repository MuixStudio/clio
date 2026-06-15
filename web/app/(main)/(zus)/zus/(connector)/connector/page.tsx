"use client"

import { useCallback, useEffect, useRef, useState } from "react"
import type { PaginationState } from "@tanstack/react-table"
import { toast } from "sonner"

import { useZusTeam } from "@/hooks/use-zus-team"
import {
  disableConnector,
  enableConnector,
  listConnectors,
  type Connector,
} from "@/service/zus-connector"

import { ConnectorInstancesCard } from "./components/connector-instances-card"
import { ConnectorPageHeader } from "./components/connector-page-header"
import { PROVIDER_FILTER_ALL } from "./components/constants"
import { ProviderSidebarCard } from "./components/provider-sidebar-card"
import type { ConnectorState, ConnectorStatusFilter } from "./components/types"

export default function Page() {
  const { team, isLoading: isTeamLoading } = useZusTeam()
  const [connectorState, setConnectorState] = useState<ConnectorState | null>(null)
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 10,
  })
  const [statusFilter, setStatusFilter] = useState<ConnectorStatusFilter>("all")
  const [providerFilter, setProviderFilter] = useState(PROVIDER_FILTER_ALL)
  const [error, setError] = useState<string | null>(null)
  const [updatingId, setUpdatingId] = useState<string | null>(null)
  const requestIdRef = useRef(0)

  const loadConnectors = useCallback(async () => {
    if (!team) return

    const requestId = ++requestIdRef.current
    const statusParam = statusFilter === "all" ? undefined : statusFilter
    const typeParam =
      providerFilter === PROVIDER_FILTER_ALL ? undefined : providerFilter

    try {
      const { data } = await listConnectors(
        team.id,
        pagination.pageIndex + 1,
        pagination.pageSize,
        statusParam,
        typeParam
      )
      if (requestIdRef.current !== requestId) return

      setError(null)
      setConnectorState({
        teamId: team.id,
        pageIndex: pagination.pageIndex,
        pageSize: pagination.pageSize,
        statusFilter,
        providerFilter,
        count: data.count,
        connectors: data.connectors,
      })
    } catch (err) {
      if (requestIdRef.current !== requestId) return
      console.error(err)
      setError("Failed to load connectors")
      setConnectorState({
        teamId: team.id,
        pageIndex: pagination.pageIndex,
        pageSize: pagination.pageSize,
        statusFilter,
        providerFilter,
        count: 0,
        connectors: [],
      })
    }
  }, [team, pagination.pageIndex, pagination.pageSize, statusFilter, providerFilter])

  useEffect(() => {
    void (async () => {
      await loadConnectors()
    })()
  }, [loadConnectors])

  const toggleEnabled = useCallback(
    async (connector: Connector) => {
      if (!team) return

      setUpdatingId(connector.id)
      try {
        if (connector.enabled) await disableConnector(team.id, connector.id)
        else await enableConnector(team.id, connector.id)
        toast.success(connector.enabled ? "Connector disabled" : "Connector enabled")
        await loadConnectors()
      } catch (err) {
        console.error(err)
        toast.error("Operation failed")
      } finally {
        setUpdatingId(null)
      }
    },
    [team, loadConnectors]
  )

  const isCurrentConnectorState = Boolean(
    team &&
      connectorState?.teamId === team.id &&
      connectorState.pageIndex === pagination.pageIndex &&
      connectorState.pageSize === pagination.pageSize &&
      connectorState.statusFilter === statusFilter &&
      connectorState.providerFilter === providerFilter
  )
  const connectors = isCurrentConnectorState ? connectorState?.connectors ?? [] : []
  const totalCount = isCurrentConnectorState ? connectorState?.count ?? 0 : 0
  const pageCount = Math.max(1, Math.ceil(totalCount / pagination.pageSize))
  const loading = isTeamLoading || Boolean(team && !isCurrentConnectorState)

  const handlePaginationChange = (next: PaginationState) => {
    const nextPageCount = Math.max(1, Math.ceil(totalCount / next.pageSize))
    setPagination({
      pageIndex: Math.min(Math.max(next.pageIndex, 0), nextPageCount - 1),
      pageSize: next.pageSize,
    })
  }

  const handleStatusFilterChange = (value: string) => {
    setStatusFilter(value as ConnectorStatusFilter)
    setPagination((prev) => ({ ...prev, pageIndex: 0 }))
  }

  const handleProviderFilterChange = (value: string) => {
    setProviderFilter(value)
    setPagination((prev) => ({ ...prev, pageIndex: 0 }))
  }

  return (
    <div className="flex flex-col gap-4">
      <ConnectorPageHeader
        teamName={team?.name}
        totalCount={totalCount}
        loading={loading}
      />

      <div className="grid gap-4 px-4 lg:grid-cols-[280px_1fr]">
        <ProviderSidebarCard />

        <ConnectorInstancesCard
          error={error}
          connectors={connectors}
          loading={loading}
          pagination={pagination}
          pageCount={pageCount}
          statusFilter={statusFilter}
          providerFilter={providerFilter}
          updatingId={updatingId}
          onRetry={() => void loadConnectors()}
          onPaginationChange={handlePaginationChange}
          onStatusFilterChange={handleStatusFilterChange}
          onProviderFilterChange={handleProviderFilterChange}
          onToggleEnabled={toggleEnabled}
        />
      </div>
    </div>
  )
}