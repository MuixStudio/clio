"use client";
import type { MenuProps } from "@/components/menu/store";

import { Icon } from "@iconify/react";

import {
  SidebarItem,
  SidebarItemType,
} from "@/components/menu/sidebar/sidebar";

export const items: SidebarItem[] = [
  {
    key: "menu-1",
    href: "/cmdb/menu-1",
    icon: "solar:home-2-linear",
    title: "overview",
  },
  {
    key: "menu-2",
    href: "/cmdb/menu-2",
    icon: "solar:widget-2-outline",
    type: SidebarItemType.Nest,
    title: "list",
    items: [
      {
        key: "menu-3",
        href: "/cmdb/menu-1",
        icon: "solar:home-2-linear",
        title: "overview",
      },
      {
        key: "menu-4",
        href: "/cmdb/menu-1",
        icon: "solar:home-2-linear",
        title: "overview",
      },
    ],
    endContent: (
      <Icon
        className="text-default-400"
        icon="solar:add-circle-line-duotone"
        width={24}
      />
    ),
  },
  {
    key: "menu-5",
    href: "/cmdb/menu-2",
    icon: "solar:widget-2-outline",
    title: "list",
    items: [
      {
        key: "menu-6",
        href: "/cmdb/menu-1",
        icon: "solar:home-2-linear",
        title: "overview",
      },
      {
        key: "menu-7",
        href: "/cmdb/menu-1",
        icon: "solar:home-2-linear",
        title: "overview",
      },
    ],
    endContent: (
      <Icon
        className="text-default-400"
        icon="solar:add-circle-line-duotone"
        width={24}
      />
    ),
  },

    {
        key: "menu-5",
        href: "/cmdb/menu-2",
        icon: "solar:widget-2-outline",
        title: "list",
        items: [
            {
                key: "menu-6",
                href: "/cmdb/menu-1",
                icon: "solar:home-2-linear",
                title: "overview",
            },
            {
                key: "menu-7",
                href: "/cmdb/menu-1",
                icon: "solar:home-2-linear",
                title: "overview",
            },
        ],
        endContent: (
            <Icon
                className="text-default-400"
                icon="solar:add-circle-line-duotone"
                width={24}
            />
        ),
    },
  {
    key: "tasks",
    href: "#",
    icon: "solar:checklist-minimalistic-outline",
    title: "task",
    endContent: (
      <Icon
        className="text-default-400"
        icon="solar:add-circle-line-duotone"
        width={24}
      />
    ),
  },
];

export const menu: MenuProps = {
  title: "On-call",
  description: "一站式告警响应平台",
  defaultSelectedKey: "menu-1",
  items,
  customFooterContent: <div>ppp</div>,
};
