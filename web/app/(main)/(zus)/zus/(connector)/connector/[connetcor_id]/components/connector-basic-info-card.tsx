import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

import { InfoItem } from "./info-item"
import type { ConnectorItem, ConnectorProvider } from "./types"

export function ConnectorBasicInfoCard({
  connector,
  provider,
  teamName,
}: {
  connector: ConnectorItem
  provider: ConnectorProvider
  teamName?: string
}) {
  const labelEntries = Object.entries(connector.labels ?? {})
  const labels =
    labelEntries.length > 0
      ? labelEntries.map(([key, value]) => `${key}=${value}`).join(", ")
      : "-"

  return (
    <Card>
      <CardHeader>
        <CardTitle>Basic Information</CardTitle>
      </CardHeader>
      <CardContent className="grid grid-cols-2 gap-4 p-4 text-sm md:grid-cols-4">
        <InfoItem label="Instance Name" value={connector.name} />
        <InfoItem label="Type" value={provider.name} />
        <InfoItem label="Team" value={teamName ?? "-"} />
        <InfoItem label="Default Labels" value={labels} />
      </CardContent>
    </Card>
  )
}
