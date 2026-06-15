import { MoreHorizontal, Power, Send, Trash2 } from "lucide-react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

import { StatusBadge } from "./status-badge"
import type { ConnectorItem, ConnectorProvider } from "./types"

export function ConnectorSummaryCard({
  connector,
  provider,
  createdAtLabel,
  isUpdating,
  onToggleEnabled,
  onRequestDelete,
}: {
  connector: ConnectorItem
  provider: ConnectorProvider
  createdAtLabel: string
  isUpdating: boolean
  onToggleEnabled: () => void
  onRequestDelete: () => void
}) {
  const Icon = provider.icon

  return (
    <Card>
      <CardContent className="flex flex-wrap items-center gap-4 p-4">
        <span className="flex size-11 shrink-0 items-center justify-center rounded-lg bg-muted text-foreground">
          <Icon className="size-5" />
        </span>
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <span className="font-semibold">{connector.name}</span>
            <Badge variant="outline">{provider.name}</Badge>
          </div>
          <div className="mt-1 text-xs text-muted-foreground">
            Instance ID {connector.id} · Created at {createdAtLabel}
          </div>
        </div>
        <StatusBadge enabled={connector.enabled} />
        <Button type="button" variant="outline" disabled title="Send test event API is not available yet">
          <Send />
          Send Test
        </Button>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button type="button" variant="outline" size="icon">
              <MoreHorizontal />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={onToggleEnabled} disabled={isUpdating}>
              <Power />
              {connector.enabled ? "Disable" : "Enable"}
            </DropdownMenuItem>
            <DropdownMenuItem
              variant="destructive"
              disabled={isUpdating}
              onSelect={(event) => {
                event.preventDefault()
                onRequestDelete()
              }}
            >
              <Trash2 />
              Delete Connector
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </CardContent>
    </Card>
  )
}
