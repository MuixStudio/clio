import { Badge } from "@/components/ui/badge"
import { cn } from "@/lib/utils"

export function AlertStatusBadge({ status }: { status: string }) {
  const firing = status === "firing"

  return (
    <Badge
      className={cn(
        "shrink-0",
        firing
          ? "border-destructive/20 bg-destructive/15 text-destructive"
          : "border-success/20 bg-success/15 text-success"
      )}
    >
      <span
        className={cn(
          "size-1.5 rounded-full",
          firing ? "bg-destructive" : "bg-success"
        )}
      />
      {firing ? "Firing" : "Resolved"}
    </Badge>
  )
}