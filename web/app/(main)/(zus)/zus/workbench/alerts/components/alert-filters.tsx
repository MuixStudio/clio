import { CalendarIcon, ListFilter, X } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Calendar } from "@/components/ui/calendar"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import type { Connector } from "@/service/zus-connector"

import { CONNECTOR_FILTER_ALL, STATUS_OPTIONS } from "./constants"
import type { AlertFiltersValue } from "./types"
import { formatDateRangeLabel } from "./utils"

export function AlertFilters({
  value,
  connectors,
  hasActiveFilters,
  onChange,
  onApply,
  onClear,
}: {
  value: AlertFiltersValue
  connectors: Connector[]
  hasActiveFilters: boolean
  onChange: (value: AlertFiltersValue) => void
  onApply: () => void
  onClear: () => void
}) {
  const dateRangeLabel = formatDateRangeLabel(value.dateRange)

  return (
    <>
      <Select
        value={value.status}
        onValueChange={(status) => onChange({ ...value, status })}
      >
        <SelectTrigger size="sm" className="w-28">
          <SelectValue placeholder="Status" />
        </SelectTrigger>
        <SelectContent position="popper">
          <SelectGroup>
            {STATUS_OPTIONS.map((option) => (
              <SelectItem key={option.value} value={option.value}>
                {option.label}
              </SelectItem>
            ))}
          </SelectGroup>
        </SelectContent>
      </Select>

      <Select
        value={value.connectorId}
        onValueChange={(connectorId) => onChange({ ...value, connectorId })}
      >
        <SelectTrigger size="sm" className="w-44">
          <SelectValue placeholder="Connector" />
        </SelectTrigger>
        <SelectContent position="popper">
          <SelectGroup>
            <SelectItem value={CONNECTOR_FILTER_ALL}>All Connectors</SelectItem>
            {connectors.map((connector) => (
              <SelectItem key={connector.id} value={connector.id}>
                {connector.name}
              </SelectItem>
            ))}
          </SelectGroup>
        </SelectContent>
      </Select>

      <Popover>
        <PopoverTrigger asChild>
          <Button variant="outline" size="sm" className="font-normal">
            <CalendarIcon />
            {dateRangeLabel ?? "Date Range"}
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-auto p-0" align="start">
          <Calendar
            autoFocus
            captionLayout="dropdown"
            mode="range"
            selected={value.dateRange}
            onSelect={(dateRange) => onChange({ ...value, dateRange })}
          />
        </PopoverContent>
      </Popover>

      <Button type="button" size="sm" onClick={onApply}>
        <ListFilter />
        Filter
      </Button>

      {hasActiveFilters && (
        <Button type="button" variant="outline" size="sm" onClick={onClear}>
          <X />
          Clear Filters
        </Button>
      )}
    </>
  )
}
