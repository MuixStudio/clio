import { Info } from "lucide-react"

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

export function ConnectorStatsCard() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Recent Alerts & Delivery Stats</CardTitle>
        <CardDescription>Alert and delivery statistics for this connector</CardDescription>
      </CardHeader>
      <CardContent className="p-4">
        <div className="flex flex-col items-center gap-2 rounded-lg border border-dashed p-10 text-center text-sm text-muted-foreground">
          <Info className="size-5" />
          Per-connector alert / delivery statistics API is not available yet
        </div>
      </CardContent>
    </Card>
  )
}
