"use client"

import * as React from "react"
import Link from "next/link"
import { ChevronRight } from "lucide-react"

import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible"
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from "@/components/ui/sidebar"
import { Skeleton } from "@/components/ui/skeleton"
import { useIsActive, useActivePath } from "@/hooks/use-active-item"
import { Menu, NavItem } from "@/components/app-sidebar/sidebar"

// ─── Menu Context ─────────────────────────────────────────────────────────────

type MenuContextProps = {
  menu?: Menu
  loading: boolean
  /**
   * Register or unregister a menu for a route segment.
   * routeKey must be a stable path prefix (e.g. "/cmdb", "/settings").
   * The menu with the longest registered key wins (deepest route takes priority).
   * Pass undefined to unregister (call in useEffect cleanup).
   */
  setMenu: (menu: Menu | undefined, routeKey: string) => void
  setLoading: (loading: boolean) => void
}

const MenuContext = React.createContext<MenuContextProps | null>(null)

export function useMenu() {
  const ctx = React.useContext(MenuContext)
  if (!ctx) throw new Error("useMenu must be used within a MenuProvider.")
  return ctx
}

/**
 * Provides menu state for the sidebar.
 * Each route layout registers its menu with a stable routeKey.
 * The deepest (longest) key wins. Cleanup removes only that layout's entry,
 * so navigating back restores the parent route's menu automatically.
 */
export function MenuProvider({
  menu: initialMenu,
  children,
}: {
  menu?: Menu
  children: React.ReactNode
}) {
  const [menu, setMenuState] = React.useState<Menu | undefined>(initialMenu)
  const [loading, setLoading] = React.useState(false)
  // Key: stable route prefix. Value: the menu registered for that prefix.
  const registryRef = React.useRef<Map<string, Menu>>(new Map())

  const setMenu = React.useCallback(
    (newMenu: Menu | undefined, routeKey: string) => {
      if (newMenu === undefined) {
        registryRef.current.delete(routeKey)
      } else {
        registryRef.current.set(routeKey, newMenu)
      }

      if (registryRef.current.size === 0) {
        setMenuState(undefined)
        return
      }

      // Show the menu registered for the longest (deepest) route key.
      const deepestKey = Array.from(registryRef.current.keys()).reduce(
        (a, b) => (b.length >= a.length ? b : a),
      )
      setMenuState(registryRef.current.get(deepestKey))
    },
    [],
  )

  const value = React.useMemo<MenuContextProps>(
    () => ({ menu, loading, setMenu, setLoading }),
    [menu, loading, setMenu],
  )

  return <MenuContext.Provider value={value}>{children}</MenuContext.Provider>
}

// ─── Skeleton ─────────────────────────────────────────────────────────────────

const SKELETON_WIDTHS = [88, 112, 72, 96, 80] as const

function NavMainSkeleton() {
  return (
    <SidebarGroup>
      <SidebarGroupLabel>
        <Skeleton className="h-3 w-20" />
      </SidebarGroupLabel>
      <SidebarMenu>
        {SKELETON_WIDTHS.map((width, i) => (
          <SidebarMenuItem key={i}>
            <SidebarMenuButton>
              <Skeleton className="size-4 shrink-0 rounded-sm" />
              <Skeleton className="h-3.5 rounded-sm" style={{ width }} />
            </SidebarMenuButton>
          </SidebarMenuItem>
        ))}
      </SidebarMenu>
    </SidebarGroup>
  )
}

// ─── Menu items ───────────────────────────────────────────────────────────────

function MenuSubItem({ subItem }: { subItem: { title: string; url: string } }) {
  const isActive = useIsActive(subItem.url)
  return (
    <SidebarMenuSubItem>
      <SidebarMenuSubButton asChild isActive={isActive}>
        <Link href={subItem.url}>
          <span>{subItem.title}</span>
        </Link>
      </SidebarMenuSubButton>
    </SidebarMenuSubItem>
  )
}

function MenuItem({ item }: { item: NavItem }) {
  const isLeaf = !item.items || item.items.length === 0
  // Leaf nodes use exact match; parent nodes use prefix match so they stay
  // highlighted when a child route is active.
  const isActive = useIsActive(item.url, isLeaf)
  const pathname = useActivePath()

  // Leaf node — no sub-items, render as a plain link
  if (isLeaf) {
    return (
      <SidebarMenuItem>
        <SidebarMenuButton asChild isActive={isActive} tooltip={item.title}>
          <Link href={item.url}>
            {item.icon && <item.icon />}
            <span>{item.title}</span>
          </Link>
        </SidebarMenuButton>
      </SidebarMenuItem>
    )
  }

  const subItems = item.items!
  const hasActiveChild = subItems.some((sub) => pathname.startsWith(sub.url))
  const defaultOpen = item.isActive || isActive || hasActiveChild

  return (
    <Collapsible asChild className="group/collapsible" defaultOpen={defaultOpen}>
      <SidebarMenuItem>
        <CollapsibleTrigger asChild>
          <SidebarMenuButton isActive={isActive} tooltip={item.title}>
            {item.icon && <item.icon />}
            <span>{item.title}</span>
            <ChevronRight className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
          </SidebarMenuButton>
        </CollapsibleTrigger>
        <CollapsibleContent>
          <SidebarMenuSub>
            {subItems.map((sub) => (
              <MenuSubItem key={sub.title} subItem={sub} />
            ))}
          </SidebarMenuSub>
        </CollapsibleContent>
      </SidebarMenuItem>
    </Collapsible>
  )
}

// ─── NavMain ──────────────────────────────────────────────────────────────────

/**
 * Renders the sidebar navigation groups. Reads the current menu from MenuContext;
 * pass `menu` to override the context value (useful for static sidebars).
 *
 * navGroups — each group gets its own SidebarGroup with optional label.
 *
 * Menu.header and Menu.footer are rendered by AppSidebar in the dedicated
 * SidebarHeader/SidebarFooter slots, not here.
 */
export function NavMain({ menu: overrideMenu }: { menu?: Menu } = {}) {
  const { menu: contextMenu, loading } = useMenu()
  const menu = overrideMenu ?? contextMenu

  if (loading) return <NavMainSkeleton />
  if (!menu) return null

  return (
    <>
      {menu.navGroups.map((group, i) => (
        <SidebarGroup key={group.name ?? i}>
          {group.name && <SidebarGroupLabel>{group.name}</SidebarGroupLabel>}
          <SidebarMenu>
            {group.items.map((item) => (
              <MenuItem key={item.title} item={item} />
            ))}
          </SidebarMenu>
        </SidebarGroup>
      ))}
    </>
  )
}