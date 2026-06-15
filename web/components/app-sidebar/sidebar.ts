import type { LucideIcon } from "lucide-react";

import React from "react";

export interface NavItem {
  title: string;
  url: string;
  /** LucideIcon component — resolved from icon registry */
  icon?: LucideIcon;
  isActive?: boolean;
  items?: {
    title: string;
    url: string;
    icon?: LucideIcon;
  }[];
}

export interface NavGroup {
  /** Optional label shown above the group */
  name?: string;
  items: NavItem[];
}

export interface Menu {
  /** Rendered by AppSidebar in the SidebarHeader slot */
  header?: MenuHeader;
  /** Rendered by AppSidebar in the SidebarFooter slot */
  footer?: MenuFooter;
  navGroups: NavGroup[];
}

export interface MenuHeader {
  title: string;
  element?: React.ReactNode;
}

export interface MenuFooter {
  element?: React.ReactNode;
}

// ─── API types (icon as string, serializable over JSON) ──────────────────────

export interface ApiNavItem {
  title: string;
  url: string;
  /** Icon name, e.g. "Settings2". Must exist in sidebarIconRegistry. */
  icon?: string;
  items?: {
    title: string;
    url: string;
    icon?: string;
  }[];
}

export interface ApiNavGroup {
  name?: string;
  items: ApiNavItem[];
}

export interface ApiMenu {
  navGroups: ApiNavGroup[];
}