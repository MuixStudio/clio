"use client"

import type { Menu } from "@/components/app-sidebar/sidebar"
import { KeyRound, User } from "lucide-react"
import { SidebarMenuLoader } from "@/components/app-sidebar/sidebar-menu-loader"

const MENU: Menu = {
  header: { title: "CMDB" },
  navGroups: [
    {
      items: [
        {
          title: "Dashboard",
          url: "/cmdb",
          icon: User,
        },
        {
          title: "Models",
          url: "/cmdb/models",
          icon: KeyRound,
        },
      ],
    },
  ],
}

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <>
      <SidebarMenuLoader menu={MENU} routeKey="/cmdb" />
      <div className="mx-auto max-w-362 p-6">
        <div className="w-full">{children}</div>
      </div>
    </>
  )
}