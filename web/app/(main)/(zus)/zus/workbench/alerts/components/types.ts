import type { DateRange } from "react-day-picker"

export type AlertFiltersValue = {
  status: string
  connectorId: string
  dateRange: DateRange | undefined
}