"use client"

import { useEffect } from "react"
import { useMenu } from "@/components/app-sidebar/nav-main"
import type { Menu } from "@/components/app-sidebar/sidebar"

type SidebarMenuLoaderProps = {
  /** The menu config to register for this route segment. */
  menu: Menu
  /**
   * Stable route prefix that uniquely identifies this layout segment,
   * e.g. "/cmdb" or "/settings". Used as the registry key — the deepest
   * (longest) registered key wins and is shown in the sidebar.
   */
  routeKey: string
}

/**
 * Drop into any route layout to register that segment's sidebar menu.
 * Renders nothing — side-effect only.
 *
 * @example
 * // app/(main)/cmdb/layout.tsx
 * <SidebarMenuLoader menu={CMDB_MENU} routeKey="/cmdb" />
 */
export function SidebarMenuLoader({ menu, routeKey }: SidebarMenuLoaderProps) {
  const { setMenu } = useMenu()

  useEffect(() => {
    setMenu(menu, routeKey)
    return () => setMenu(undefined, routeKey)
  }, [menu, routeKey, setMenu])

  return null
}