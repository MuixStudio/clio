import type { Connector } from "@/service/zus-connector"

import type { STATUS_FILTERS } from "./constants"

export type ConnectorStatusFilter = (typeof STATUS_FILTERS)[number]["value"]

export type ConnectorState = {
  teamId: string
  pageIndex: number
  pageSize: number
  statusFilter: ConnectorStatusFilter
  providerFilter: string
  count: number
  connectors: Connector[]
}