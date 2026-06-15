"use client"

import * as React from "react"
import {
  getCoreRowModel,
  useReactTable,
  type PaginationState,
  type Updater,
} from "@tanstack/react-table"

import { DataTable } from "@/components/data-table/data-table"
import { DataTableViewOptions } from "@/components/data-table/data-table-view-options"
import type { Connector } from "@/service/zus-connector"

import { getConnectorColumns } from "./columns"
import { ConnectorFilters } from "./connector-filters"
import { ConnectorStatePanel } from "./connector-state-panel"
import { PROVIDER_FILTER_ALL } from "./constants"
import type { ConnectorStatusFilter } from "./types"

export function ConnectorTable({
  connectors,
  pagination,
  pageCount,
  statusFilter,
  providerFilter,
  updatingId,
  onPaginationChange,
  onStatusFilterChange,
  onProviderFilterChange,
  onToggleEnabled,
}: {
  connectors: Connector[]
  pagination: PaginationState
  pageCount: number
  statusFilter: ConnectorStatusFilter
  providerFilter: string
  updatingId: string | null
  onPaginationChange: (pagination: PaginationState) => void
  onStatusFilterChange: (value: string) => void
  onProviderFilterChange: (value: string) => void
  onToggleEnabled: (connector: Connector) => void
}) {
  const columns = React.useMemo(
    () => getConnectorColumns({ updatingId, onToggleEnabled }),
    [updatingId, onToggleEnabled]
  )

  const table = useReactTable({
    data: connectors,
    columns,
    pageCount,
    state: { pagination },
    manualPagination: true,
    enableSorting: false,
    getRowId: (row) => row.id,
    onPaginationChange: (updater: Updater<PaginationState>) => {
      onPaginationChange(
        typeof updater === "function" ? updater(pagination) : updater
      )
    },
    getCoreRowModel: getCoreRowModel(),
  })

  const hasActiveFilters =
    statusFilter !== "all" || providerFilter !== PROVIDER_FILTER_ALL

  const emptyState =
    connectors.length === 0 ? (
      hasActiveFilters ? (
        <ConnectorStatePanel
          type="filtered-empty"
          onReset={() => {
            onStatusFilterChange("all")
            onProviderFilterChange(PROVIDER_FILTER_ALL)
          }}
        />
      ) : (
        <ConnectorStatePanel type="empty" />
      )
    ) : undefined

  return (
    <DataTable table={table} emptyState={emptyState}>
      <div className="flex items-center justify-between gap-2 p-1">
        <ConnectorFilters
          statusFilter={statusFilter}
          providerFilter={providerFilter}
          onStatusFilterChange={onStatusFilterChange}
          onProviderFilterChange={onProviderFilterChange}
        />
        <DataTableViewOptions table={table} />
      </div>
    </DataTable>
  )
}