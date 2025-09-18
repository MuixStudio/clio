"use client";

import React from "react";
import { Button, Divider, ScrollShadow, cn } from "@heroui/react";
import { useShallow } from "zustand/react/shallow";
import { Icon } from "@iconify/react";

import SidebarDrawer from "./sidebar/sidebar-drawer";
import Sidebar from "./sidebar/sidebar";
import MenuHeader from "./menu-header";
import MenuFooter from "./menu-footer";

import { useMenuStore } from "@/components/menu/store";

type MenuProps = {
  isCompact?: boolean;
  className?: string;
  onToggle?: () => void;
  ref?: React.Ref<HTMLAnchorElement>;
};

const Menu = React.forwardRef<HTMLAnchorElement, MenuProps>(
  ({ isCompact, onToggle, className, ...props }, ref) => {
    const {
      title,
      items,
      customFooterContent,
      customHeaderContent,
      description,
      isOpen,
      onOpenChange,
      isLoaded,
    } = useMenuStore(
      useShallow((state) => ({
        title: state.title,
        items: state.items,
        isOpen: state.isOpen,
        onOpenChange: state.onOpenChange,
        customHeaderContent: state.customHeaderContent,
        description: state.description,
        customFooterContent: state.customFooterContent,
        isLoaded: state.isLoaded,
      })),
    );

    return (
      <>
        <SidebarDrawer
          className="overflow-y-auto"
          isOpen={isOpen}
          sidebarWidth={300}
          onOpenChange={onOpenChange}
        >
          <div className="relative flex h-full w-56 flex-col">
            {/* Header */}
            <div className="h-24 shrink-0 grow-0 overflow-y-auto">
              <MenuHeader
                customHeaderContent={customHeaderContent}
                description={description}
                isLoaded={isLoaded}
                title={title}
              />
            </div>
            <Divider />
            {/* Sidebar */}
            <ScrollShadow className="h-full max-h-full py-6">
              <Sidebar
                ref={ref}
                items={items}
                isLoaded={isLoaded}
                // defaultSelectedKey={items[0].key}
              />
            </ScrollShadow>
            {/* Footer */}
            <div className="max-h-36 shrink-0 grow-0 overflow-y-auto">
              <MenuFooter customFooterContent={customFooterContent} />
            </div>
          </div>
        </SidebarDrawer>

        <div
          className={cn(
            "!border-r-small border-divider flex h-full w-full flex-col",
            {
              "items-center": isCompact,
            },
            className,
          )}
        >
          <div
            className={cn(
              "h-24 w-full shrink-0 grow-0 overflow-y-auto",

              {
                "flex justify-center items-center": isCompact,
              },
            )}
          >
            <MenuHeader
              customHeaderContent={customHeaderContent}
              description={description}
              isCompact={isCompact}
              isLoaded={isLoaded}
              title={title}
            />
          </div>

          <div className="shrink-0">
            <Button
              className="min-w-0"
              fullWidth={true}
              variant="light"
              onPress={onToggle}
            >
              <Icon
                className="text-default-500"
                height={24}
                icon={
                  isCompact ? "tabler:arrow-bar-right" : "tabler:arrow-bar-left"
                }
                width={24}
              />
            </Button>
          </div>

          <ScrollShadow className="h-full w-full py-6">
            <Sidebar
              isCompact={isCompact}
              isLoaded={isLoaded}
              items={items}
              // defaultSelectedKey={items[0].key}
              ref={ref}
            />
          </ScrollShadow>

          {/* Footer */}
          <div className="max-h-36 shrink-0 grow-0 overflow-y-auto">
            <MenuFooter customFooterContent={customFooterContent} />
          </div>
        </div>
      </>
    );
  },
);

export default React.memo(Menu);
