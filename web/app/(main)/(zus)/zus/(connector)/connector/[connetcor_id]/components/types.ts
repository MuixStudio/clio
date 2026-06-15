import type { Provider } from "@/models/connector-provider/types"

export type ConnectorItem = {
  id: string
  type: string
  name: string
  token?: string
  labels: Record<string, string>
  enabled: boolean
  created_at: string
  updated_at?: string
}

export type ConnectorProvider = Provider
