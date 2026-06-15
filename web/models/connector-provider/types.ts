import type { ComponentType } from "react"

export type Providers = Provider[]

export type ProviderDocs = {
  title: string
  Content: ComponentType
}

export type ProviderDefaults = {
  instanceName?: string
  teamId?: string
  labels?: Record<string, string>
}

export type Provider = {
  id: string
  name: string
  description: string
  icon: ComponentType<{ className?: string }>
  category: string
  subtitle?: string
  pushMethod?: string
  tags: string[]
  defaults?: ProviderDefaults
  docs?: ProviderDocs
}
