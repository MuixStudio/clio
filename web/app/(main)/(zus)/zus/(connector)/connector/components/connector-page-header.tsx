import Link from "next/link"
import { Plus } from "lucide-react"

import {
  HeaderBar,
  HeaderBarActions,
  HeaderBarDescription,
  HeaderBarLeft,
  HeaderBarTitle,
} from "@/app/(main)/components/header-bar/header-bar"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"

export function ConnectorPageHeader({
  teamName,
  totalCount,
  loading,
}: {
  teamName?: string
  totalCount: number
  loading: boolean
}) {
  return (
    <HeaderBar>
      <HeaderBarLeft>
        <HeaderBarTitle className="flex items-center gap-2">
          Connectors
          {!loading && <Badge variant="outline">{totalCount}</Badge>}
        </HeaderBarTitle>
        <HeaderBarDescription>
          {teamName
            ? `Alert source instances connected to ${teamName}`
            : "Select a team to view connectors"}
        </HeaderBarDescription>
      </HeaderBarLeft>
      <HeaderBarActions>
        <Button asChild>
          <Link href="/zus/connector/new">
            <Plus />
            Add Connector
          </Link>
        </Button>
      </HeaderBarActions>
    </HeaderBar>
  )
}
