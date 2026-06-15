"use client"

import type * as React from "react"
import {
  BotIcon,
  ChevronRightIcon,
  CommandIcon,
  HashIcon,
  InboxIcon,
  SearchIcon,
  SparklesIcon,
} from "lucide-react"

import { cn } from "@/lib/utils"
import {
  SECONDARY_SIDEBAR_MAX_WIDTH,
  SECONDARY_SIDEBAR_MIN_WIDTH,
  SecondarySidebar,
  SecondarySidebarContent,
  SecondarySidebarFooter,
  SecondarySidebarGroup,
  SecondarySidebarGroupContent,
  SecondarySidebarGroupLabel,
  SecondarySidebarHeader,
  SecondarySidebarMenu,
  SecondarySidebarMenuButton,
  SecondarySidebarMenuItem,
  SecondarySidebarResizeHandle,
} from "@/components/ui/secondary-sidebar"

const workbenchItems = [
  { label: "Overview", icon: InboxIcon, active: true },
  { label: "Incidents", icon: SparklesIcon },
  { label: "Channels", icon: HashIcon },
  { label: "Agents", icon: BotIcon },
]

export function WorkbenchSidebar({
  className,
  ...props
}: React.ComponentProps<typeof SecondarySidebar>) {
  return (
    <SecondarySidebar className={cn("group/workbench-sidebar", className)} {...props}>
      <SecondarySidebarHeader>
        <div className="flex items-center gap-2 rounded-xl border border-sidebar-border/70 bg-sidebar-accent/30 p-2 shadow-xs">
          <div className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground shadow-sm">
            <CommandIcon className="size-4" />
          </div>
          <div className="grid min-w-0 flex-1 text-left leading-tight">
            <span className="truncate text-sm font-semibold">Workbench</span>
            <span className="truncate text-xs text-sidebar-foreground/60">
              Zus operations
            </span>
          </div>
        </div>
      </SecondarySidebarHeader>

      <SecondarySidebarContent>
        <SecondarySidebarGroup>
          <div className="px-1 pb-2">
            <div className="flex h-8 items-center gap-2 rounded-lg border border-sidebar-border/70 bg-background/60 px-2 text-xs text-sidebar-foreground/55 shadow-xs">
              <SearchIcon className="size-3.5 shrink-0" />
              <span className="truncate">Search workbench</span>
            </div>
          </div>

          <SecondarySidebarGroupLabel>Space</SecondarySidebarGroupLabel>
          <SecondarySidebarGroupContent>
            <SecondarySidebarMenu>
              {workbenchItems.map((item) => (
                <SecondarySidebarMenuItem key={item.label}>
                  <SecondarySidebarMenuButton isActive={item.active}>
                    <item.icon />
                    <span>{item.label}</span>
                    {item.active ? <ChevronRightIcon className="ml-auto" /> : null}
                  </SecondarySidebarMenuButton>
                </SecondarySidebarMenuItem>
              ))}
            </SecondarySidebarMenu>
          </SecondarySidebarGroupContent>
        </SecondarySidebarGroup>
      </SecondarySidebarContent>

      <SecondarySidebarFooter>
        <div className="rounded-xl border border-sidebar-border/70 bg-sidebar-accent/25 p-3">
          <div className="text-xs font-medium">Resizable menu</div>
          <div className="mt-1 text-xs leading-5 text-sidebar-foreground/55">
            Drag the right edge to resize between {SECONDARY_SIDEBAR_MIN_WIDTH}px
            and {SECONDARY_SIDEBAR_MAX_WIDTH}px.
          </div>
        </div>
      </SecondarySidebarFooter>

      <SecondarySidebarResizeHandle aria-label="Resize workbench sidebar" />
    </SecondarySidebar>
  )
}