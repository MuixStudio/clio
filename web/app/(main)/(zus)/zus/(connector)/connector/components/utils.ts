import { providers } from "@/models/connector-provider/constants"

export function getProvider(type: string) {
  return providers.find((provider) => provider.id === type)
}

export function formatDate(value?: string) {
  if (!value) return "-"
  return new Date(value).toLocaleString("en-US", { hour12: false })
}

export function formatLabels(labels: Record<string, string>) {
  return Object.entries(labels ?? {})
}
