import { Badge } from "@/components/ui/badge"
import { cn } from "@/lib/utils"

export function StatusBadge({ enabled }: { enabled: boolean }) {
  return (
    <Badge
      className={cn(
        "shrink-0",
        enabled
          ? "border-success/20 bg-success/15 text-success"
          : "bg-muted text-muted-foreground"
      )}
    >
      <span
        className={cn(
          "size-1.5 rounded-full",
          enabled ? "bg-success" : "bg-muted-foreground"
        )}
      />
      {enabled ? "Enabled" : "Disabled"}
    </Badge>
  )
}
