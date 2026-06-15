import { Power, Trash2 } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"

export function ConnectorDangerCard({
  enabled,
  isUpdating,
  onToggleEnabled,
  onRequestDelete,
}: {
  enabled: boolean
  isUpdating: boolean
  onToggleEnabled: () => void
  onRequestDelete: () => void
}) {
  return (
    <Card className="ring-destructive">
      <CardContent className="flex flex-wrap items-center justify-between gap-4 p-4">
        <div>
          <div className="text-sm font-medium text-destructive">Danger Zone</div>
          <div className="text-xs text-muted-foreground">
            Disabling stops receiving pushes; deleting removes this instance (historical
            alerts are preserved).
          </div>
        </div>
        <div className="flex gap-2">
          <Button type="button" variant="outline" onClick={onToggleEnabled} disabled={isUpdating}>
            <Power />
            {enabled ? "Disable" : "Enable"}
          </Button>
          <Button
            type="button"
            variant="outline"
            onClick={onRequestDelete}
            disabled={isUpdating}
          >
            <Trash2 />
            Delete Connector
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}
