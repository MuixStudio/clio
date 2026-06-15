"use client"

import { User, KeyRound } from "lucide-react"
import { SidebarMenuLoader } from "@/components/app-sidebar/sidebar-menu-loader"
import type { Menu } from "@/components/app-sidebar/sidebar"

const SETTINGS_MENU: Menu = {
  header: { title: "Settings" },
  navGroups: [
    {
      items: [
        {
          title: "Profile",
          url: "/settings/profile",
          icon: User,
        },
        {
          title: "Account",
          url: "/settings/account",
          icon: KeyRound,
        },
      ],
    },
  ],
}

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <>
      <SidebarMenuLoader menu={SETTINGS_MENU} routeKey="/settings" />
      <div className="mx-auto max-w-362 p-6">
        <div className="w-full">{children}</div>
      </div>
    </>
  )
}