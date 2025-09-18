"use client";
import React from "react";
import { Navbar, NavbarContent, NavbarItem } from "@heroui/navbar";
import { BreadcrumbItem, Breadcrumbs } from "@heroui/breadcrumbs";
import { Button } from "@heroui/button";
import { Icon } from "@iconify/react";

import { useMenuStore } from "@/components/menu/store";

const AppHeader: React.FC = () => {
  const onOpen = useMenuStore((store) => store.onOpen);

  return (
    <>
      <Navbar
        isBordered
        className="bg-foreground-50"
        classNames={{ wrapper: "pr-6 pl-0" }}
        height="2.25rem"
        maxWidth="full"
      >
        <NavbarContent className="gap-0">
          <NavbarItem className="h-full w-9 flex sm:hidden">
            <Button
              isIconOnly
              as="div"
              className="flex h-full w-full rounded-none"
              size="sm"
              variant="light"
              onPress={onOpen}
            >
              <Icon
                className="text-default-500"
                height={24}
                icon="solar:hamburger-menu-outline"
                width={24}
              />
            </Button>
          </NavbarItem>
          <NavbarItem className="ml-2">
            <Breadcrumbs size="sm">
              <BreadcrumbItem>Home</BreadcrumbItem>
            </Breadcrumbs>
          </NavbarItem>
        </NavbarContent>
      </Navbar>
    </>
  );
};

export default AppHeader;
