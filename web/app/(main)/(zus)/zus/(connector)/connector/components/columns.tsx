"use client"

import Link from "next/link"
import type { ColumnDef } from "@tanstack/react-table"
import { Plug, Power } from "lucide-react"

import { DataTableColumnHeader } from "@/components/data-table/data-table-column-header"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import type { Connector } from "@/service/zus-connector"

import { StatusBadge } from "./status-badge"
import { formatDate, formatLabels, getProvider } from "./utils"

export function getConnectorColumns({
  updatingId,
  onToggleEnabled,
}: {
  updatingId: string | null
  onToggleEnabled: (connector: Connector) => void
}): ColumnDef<Connector>[] {
  return [
    {
      accessorKey: "name",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} label="Connector" />
      ),
      cell: ({ row }) => {
        const connector = row.original
        const provider = getProvider(connector.type)
        const Icon = provider?.icon ?? Plug
        return (
          <div className="flex items-center gap-3">
            <span className="flex size-9 shrink-0 items-center justify-center rounded-md bg-muted text-foreground">
              <Icon className="size-5" />
            </span>
            <span className="min-w-0">
              <span className="block truncate font-medium">
                {connector.name}
              </span>
              <span className="block truncate text-xs text-muted-foreground">
                Instance ID {connector.id}
              </span>
            </span>
          </div>
        )
      },
      meta: { label: "Connector" },
    },
    {
      accessorKey: "type",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} label="Provider" />
      ),
      cell: ({ row }) => {
        const provider = getProvider(row.original.type)
        return <Badge variant="outline">{provider?.name ?? row.original.type}</Badge>
      },
      meta: { label: "Provider" },
    },
    {
      id: "labels",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} label="Default Labels" />
      ),
      cell: ({ row }) => {
        const labels = formatLabels(row.original.labels)
        if (labels.length === 0) {
          return <span className="text-muted-foreground">-</span>
        }
        return (
          <div className="flex max-w-60 flex-wrap gap-1">
            {labels.map(([key, value]) => (
              <Badge key={key} variant="secondary" className="font-mono text-[10px]">
                {key}={value}
              </Badge>
            ))}
          </div>
        )
      },
      enableSorting: false,
      enableHiding: false,
    },
    {
      accessorKey: "enabled",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} label="Status" />
      ),
      cell: ({ row }) => <StatusBadge enabled={row.original.enabled} />,
      meta: { label: "Status" },
    },
    {
      accessorKey: "created_at",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} label="Created At" />
      ),
      cell: ({ row }) => (
        <span className="text-muted-foreground">
          {formatDate(row.original.created_at)}
        </span>
      ),
      meta: { label: "Created At" },
    },
    {
      accessorKey: "updated_at",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} label="Updated At" />
      ),
      cell: ({ row }) => (
        <span className="text-muted-foreground">
          {formatDate(row.original.updated_at)}
        </span>
      ),
      meta: { label: "Updated At" },
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const connector = row.original
        return (
          <div className="flex justify-end gap-2">
            <Button asChild variant="outline" size="sm">
              <Link href={`/zus/connector/${connector.id}`}>View</Link>
            </Button>
            <Button
              type="button"
              variant="outline"
              size="sm"
              disabled={updatingId === connector.id}
              onClick={() => onToggleEnabled(connector)}
            >
              <Power />
              {connector.enabled ? "Disable" : "Enable"}
            </Button>
          </div>
        )
      },
      enableHiding: false,
      size: 180,
    },
  ]
}