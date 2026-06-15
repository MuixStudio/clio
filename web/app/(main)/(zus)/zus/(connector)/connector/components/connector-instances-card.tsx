import type { PaginationState } from "@tanstack/react-table"

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { DataTableSkeleton } from "@/components/data-table/data-table-skeleton"
import type { Connector } from "@/service/zus-connector"

import { ConnectorStatePanel } from "./connector-state-panel"
import { ConnectorTable } from "./connector-table"
import type { ConnectorStatusFilter } from "./types"

export function ConnectorInstancesCard({
  error,
  connectors,
  loading,
  pagination,
  pageCount,
  statusFilter,
  providerFilter,
  updatingId,
  onRetry,
  onPaginationChange,
  onStatusFilterChange,
  onProviderFilterChange,
  onToggleEnabled,
}: {
  error: string | null
  connectors: Connector[]
  loading: boolean
  pagination: PaginationState
  pageCount: number
  statusFilter: ConnectorStatusFilter
  providerFilter: string
  updatingId: string | null
  onRetry: () => void
  onPaginationChange: (pagination: PaginationState) => void
  onStatusFilterChange: (value: string) => void
  onProviderFilterChange: (value: string) => void
  onToggleEnabled: (connector: Connector) => void
}) {
  return (
    <Card className="min-w-0">
      <CardHeader>
        <CardTitle>Connector Instances</CardTitle>
        <CardDescription>
          View connectors by team with pagination, and enable or disable instances
        </CardDescription>
      </CardHeader>
      <CardContent className="p-1.5">
        {error ? (
          <ConnectorStatePanel type="error" message={error} onRetry={onRetry} />
        ) : loading ? (
          <DataTableSkeleton columnCount={7} filterCount={2} />
        ) : (
          <ConnectorTable
            connectors={connectors}
            pagination={pagination}
            pageCount={pageCount}
            statusFilter={statusFilter}
            providerFilter={providerFilter}
            updatingId={updatingId}
            onPaginationChange={onPaginationChange}
            onStatusFilterChange={onStatusFilterChange}
            onProviderFilterChange={onProviderFilterChange}
            onToggleEnabled={onToggleEnabled}
          />
        )}
      </CardContent>
    </Card>
  )
}