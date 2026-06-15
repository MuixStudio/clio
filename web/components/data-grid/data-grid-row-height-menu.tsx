"use client"

import type { Table } from "@tanstack/react-table"
import {
  AlignVerticalSpaceAroundIcon,
  ChevronsDownUpIcon,
  ChevronDownIcon,
  EqualIcon,
  MinusIcon,
} from "lucide-react"
import * as React from "react"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

const rowHeights = [
  {
    label: "Short",
    value: "short" as const,
    icon: MinusIcon,
  },
  {
    label: "Medium",
    value: "medium" as const,
    icon: EqualIcon,
  },
  {
    label: "Tall",
    value: "tall" as const,
    icon: AlignVerticalSpaceAroundIcon,
  },
  {
    label: "Extra Tall",
    value: "extra-tall" as const,
    icon: ChevronsDownUpIcon,
  },
] as const

interface DataGridRowHeightMenuProps<TData> {
  table: Table<TData>
  disabled?: boolean
  className?: string
}

export function DataGridRowHeightMenu<TData>({
  table,
  disabled,
  className,
}: DataGridRowHeightMenuProps<TData>) {
  const rowHeight = table.options.meta?.rowHeight
  const onRowHeightChange = table.options.meta?.onRowHeightChange

  const selectedRowHeight = React.useMemo(() => {
    return (
      rowHeights.find((opt) => opt.value === rowHeight) ?? {
        label: "Short",
        value: "short" as const,
        icon: MinusIcon,
      }
    )
  }, [rowHeight])

  const SelectedIcon = selectedRowHeight.icon

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild disabled={disabled}>
        <Button
          variant="outline"
          size="sm"
          className={cn("w-full justify-between rounded font-bold", className)}
        >
          <SelectedIcon />
          {selectedRowHeight.label}
          <ChevronDownIcon className="ml-auto size-4 opacity-50" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start" className="min-w-36">
        <DropdownMenuRadioGroup
          value={rowHeight}
          onValueChange={onRowHeightChange as (value: string) => void}
        >
          {rowHeights.map((option) => {
            const OptionIcon = option.icon
            return (
              <DropdownMenuRadioItem key={option.value} value={option.value}>
                <OptionIcon className="size-4" />
                {option.label}
              </DropdownMenuRadioItem>
            )
          })}
        </DropdownMenuRadioGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
