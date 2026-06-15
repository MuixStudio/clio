import Link from "next/link"
import { Plug } from "lucide-react"

import { Button } from "@/components/ui/button"

export function ConnectorStatePanel({
  type,
  message,
  onRetry,
  onReset,
}: {
  type: "error" | "empty" | "filtered-empty"
  message?: string
  onRetry?: () => void
  onReset?: () => void
}) {
  if (type === "error") {
    return (
      <div className="flex flex-col items-center gap-3 p-12 text-center text-sm text-muted-foreground">
        <Plug className="size-8" />
        <span>{message}</span>
        <Button type="button" variant="outline" onClick={onRetry}>
          Retry
        </Button>
      </div>
    )
  }

  if (type === "filtered-empty") {
    return (
      <div className="flex flex-col items-center gap-3 p-16 text-center text-sm text-muted-foreground">
        <Plug className="size-8" />
        No connectors match the current filters
        <Button type="button" variant="outline" onClick={onReset}>
          Reset filters
        </Button>
      </div>
    )
  }

  return (
    <div className="flex flex-col items-center gap-3 p-16 text-center text-sm text-muted-foreground">
      <Plug className="size-8" />
      No connectors yet
      <Button asChild variant="outline">
        <Link href="/zus/connect">Add a connector</Link>
      </Button>
    </div>
  )
}
