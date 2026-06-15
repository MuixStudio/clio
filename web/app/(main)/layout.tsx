"use client"

import { AppNavbar } from "@/components/app-navbar/app-navbar"
import {
  SidebarInset,
  SidebarProvider,
  useSidebar,
} from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/app-sidebar/app-sidebar"
import { cn } from "@/lib/utils"
import { MenuProvider } from "@/components/app-sidebar/nav-main"

export default function MainLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <SidebarProvider className="flex flex-col">
      <AppNavbar />
      <div className="flex min-w-0 flex-1">
        <MenuProvider>
          <AppSidebar />
          <SidebarInset className="min-w-0 flex-1">
            <I>{children}</I>
          </SidebarInset>
        </MenuProvider>
      </div>
    </SidebarProvider>
  )
}

function I({ children }: { children: React.ReactNode }) {
  const { open } = useSidebar()

  return (
    <div
      data-open={open ? "true" : undefined}
      className={cn(
        "h-[calc(100svh-var(--navbar-height))] overflow-auto transition-[width] duration-200 ease-linear",
        "w-full min-w-0 md:w-[calc(100svw-var(--sidebar-width-icon))] md:data-open:w-[calc(100svw-var(--sidebar-width))]"
      )}
    >
      {children}
    </div>
  )
}
