import type { DateRange } from "react-day-picker"

const ZERO_TIME_PREFIX = "0001-01-01"

export function formatDateTime(value?: string) {
  if (!value || value.startsWith(ZERO_TIME_PREFIX)) return "-"
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return "-"
  return date.toLocaleString("en-US", { hour12: false })
}

export function formatDateShort(date?: Date) {
  if (!date) return ""
  return date.toLocaleDateString("en-US")
}

export function formatDateRangeLabel(range?: DateRange) {
  if (!range?.from && !range?.to) return null
  if (range.from && range.to) {
    return `${formatDateShort(range.from)} - ${formatDateShort(range.to)}`
  }
  return formatDateShort(range.from ?? range.to)
}

export function toStartOfDayISOString(date?: Date) {
  if (!date) return undefined
  const d = new Date(date)
  d.setHours(0, 0, 0, 0)
  return d.toISOString()
}

export function toEndOfDayISOString(date?: Date) {
  if (!date) return undefined
  const d = new Date(date)
  d.setHours(23, 59, 59, 999)
  return d.toISOString()
}

export function formatRelativeTime(value?: string) {
  if (!value || value.startsWith(ZERO_TIME_PREFIX)) return "-"
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return "-"

  const diffSec = Math.floor((Date.now() - date.getTime()) / 1000)
  if (diffSec < 60) return "just now"
  const diffMin = Math.floor(diffSec / 60)
  if (diffMin < 60) return `${diffMin}m ago`
  const diffHour = Math.floor(diffMin / 60)
  if (diffHour < 24) return `${diffHour}h ago`
  const diffDay = Math.floor(diffHour / 24)
  if (diffDay < 30) return `${diffDay}d ago`
  return formatDateTime(value)
}