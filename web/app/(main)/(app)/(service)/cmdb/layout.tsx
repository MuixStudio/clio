"use client";
import { usePathname } from "next/navigation";

import { useMenuStore } from "@/components/menu/store";
import { menu } from "@/app/(main)/(app)/(service)/cmdb/menu";

export default function CMDBLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();

  // useMenuStore.setState(
  //   {
  //     title: menu.title,
  //     items: menu.items,
  //     defaultSelectedKey: pathname,
  //     description: menu.description,
  //     customHeaderContent: menu.title ? menu.customHeaderContent : undefined,
  //     customFooterContent: menu.customFooterContent,
  //     isLoaded: true,
  //   },
  // )
  setTimeout(() => {
    useMenuStore.setState({
      title: menu.title,
      items: menu.items,
      defaultSelectedKey: menu.defaultSelectedKey,
      description: menu.description,
      customHeaderContent: menu.title ? menu.customHeaderContent : undefined,
      customFooterContent: menu.customFooterContent,
      isLoaded: true,
    });
  }, 2000);

  return <div>{children}</div>;
}
