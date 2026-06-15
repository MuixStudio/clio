import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { providers } from "@/models/connector-provider/constants"

import { PROVIDER_FILTER_ALL, STATUS_FILTERS } from "./constants"
import type { ConnectorStatusFilter } from "./types"

export function ConnectorFilters({
  statusFilter,
  providerFilter,
  onStatusFilterChange,
  onProviderFilterChange,
}: {
  statusFilter: ConnectorStatusFilter
  providerFilter: string
  onStatusFilterChange: (value: string) => void
  onProviderFilterChange: (value: string) => void
}) {
  return (
    <div className="flex flex-wrap items-center gap-2">
      <Select value={statusFilter} onValueChange={onStatusFilterChange}>
        <SelectTrigger size="sm" className="w-32">
          <SelectValue placeholder="Status" />
        </SelectTrigger>
        <SelectContent position="popper">
          <SelectGroup>
            {STATUS_FILTERS.map((filter) => (
              <SelectItem key={filter.value} value={filter.value}>
                {filter.label}
              </SelectItem>
            ))}
          </SelectGroup>
        </SelectContent>
      </Select>
      <Select value={providerFilter} onValueChange={onProviderFilterChange}>
        <SelectTrigger size="sm" className="w-40">
          <SelectValue placeholder="Provider" />
        </SelectTrigger>
        <SelectContent position="popper">
          <SelectGroup>
            <SelectItem value={PROVIDER_FILTER_ALL}>All Providers</SelectItem>
            {providers.map((provider) => (
              <SelectItem key={provider.id} value={provider.id}>
                {provider.name}
              </SelectItem>
            ))}
          </SelectGroup>
        </SelectContent>
      </Select>
    </div>
  )
}
