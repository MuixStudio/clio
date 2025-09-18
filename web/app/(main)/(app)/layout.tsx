"use client";
import React from "react";
import { useMediaQuery } from "usehooks-ts";
import { cn } from "@heroui/react";

import Menu from "@/components/menu/menu";
import AppHeader from "@/components/app-header/app-header";

export default function AppLayout({ children }: { children: React.ReactNode }) {
  const [isCollapsed, setIsCollapsed] = React.useState(false);
  const isTablet = useMediaQuery("(max-width: 1024px)");
  const isMobile = useMediaQuery("(max-width: 640px)");
  const [isCompact, setIsCompact] = React.useState(isCollapsed || isTablet);

  React.useMemo(() => {
    setIsCompact(isTablet);
  }, [isTablet]);
  const [isMounted, setIsMounted] = React.useState(false);

  React.useEffect(() => {
    setIsMounted(true);
  }, []);

  const onToggle = () => {
    setIsCompact((prev) => !prev);
  };

  if (!isMounted) return null;

  return (
    <>
      <div
        className={cn(
          "transition-width h-full w-56 flex-col",
          {
            "w-16 items-center": isCompact,
          },
          {
            "w-0 items-center": isMobile,
          },
        )}
      >
        <Menu
          className="bg-foreground-50"
          isCompact={isCompact}
          onToggle={onToggle}
        />
      </div>
      <div className="flex h-full w-full flex-1 flex-col">
        <div className="shrink-0">
          <AppHeader />
        </div>
        <main className="h-full w-full overflow-y-auto">
          <div className="m-auto h-full w-full max-w-[1296px]">
            <div className="p-6 w-full h-full">{children}</div>
          </div>
        </main>
      </div>
    </>
  );
}
