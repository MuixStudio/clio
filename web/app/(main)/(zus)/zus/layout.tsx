"use client"

import type { Menu } from "@/components/app-sidebar/sidebar"
import { Users, BadgePlus, Cable, Siren, LayoutDashboard } from "lucide-react"
import { SidebarMenuLoader } from "@/components/app-sidebar/sidebar-menu-loader"
import { TeamSelect } from "@/components/zus/team-select"

const MENU: Menu = {
  header: { title: "zus" },
  navGroups: [
    {
      name: "All",
      items: [
        {
          title: "Overview",
          url: "/zus",
          icon: LayoutDashboard,
        },
        {
          title: "Workbench",
          url: "/zus/workbench",
          icon: Siren,
        },
      ],
    },
    {
      name: "Provider",
      items: [
        {
          title: "Connector",
          url: "/zus/connector",
          icon: Cable,
        },
      ],
    },
    {
      name: "Team",
      items: [
        {
          title: "Members",
          url: "/zus/members",
          icon: Users,
        },
      ],
    },
  ],
  // footer: { element: <TeamSelect /> },
}

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <>
      <SidebarMenuLoader menu={MENU} routeKey="/cmdb" />
      <div className="h-full w-full">{children}</div>
    </>
  )
}
