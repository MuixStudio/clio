"use client"

import * as React from "react"
import { PanelLeftClose, PanelLeftOpen } from "lucide-react"

import { NavMain, useMenu } from "./nav-main"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarRail,
  useSidebar,
} from "@/components/ui/sidebar"
import { Button } from "@/components/ui/button"

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const { toggleSidebar, open } = useSidebar()
  const { menu } = useMenu()

  return (
    <Sidebar
      collapsible="icon"
      className="top-(--navbar-height) h-[calc(100svh-var(--navbar-height))]!"
      {...props}
    >
      {menu?.header && (
        <SidebarHeader>
          {menu.header.element ?? (
            <SidebarGroupLabel>{menu.header.title}</SidebarGroupLabel>
          )}
        </SidebarHeader>
      )}
      <SidebarContent>
        <NavMain />
      </SidebarContent>
      {menu?.footer?.element && (
        <SidebarFooter>{menu.footer.element}</SidebarFooter>
      )}
      <SidebarFooter className="hidden md:flex">
        <Button variant="outline" onClick={toggleSidebar}>
          {open ? <PanelLeftClose /> : <PanelLeftOpen />}
        </Button>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
