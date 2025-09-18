"use client";
import type { ModalProps } from "@heroui/react";

import React from "react";
import { TRANSITION_EASINGS } from "@heroui/framer-utils";
import { Drawer, DrawerBody, DrawerContent, cn } from "@heroui/react";

const SidebarDrawer = React.forwardRef<
  HTMLDivElement,
  ModalProps & {
    sidebarWidth?: number;
    sidebarPlacement?: "left" | "right";
  }
>(
  (
    {
      children,
      className,
      onOpenChange,
      isOpen,
      sidebarWidth = 288,
      classNames = {},
      sidebarPlacement = "left",
      motionProps: drawerMotionProps,
      ...props
    },
    ref,
  ) => {
    const motionProps = React.useMemo(() => {
      if (!!drawerMotionProps && typeof drawerMotionProps === "object")
        return drawerMotionProps;

      return {
        variants: {
          enter: {
            x: 0,
            transition: {
              x: {
                duration: 0.3,
                ease: TRANSITION_EASINGS.easeOut,
              },
            },
          },
          exit: {
            x: sidebarPlacement === "left" ? -sidebarWidth : sidebarWidth,
            transition: {
              x: {
                duration: 0.2,
                ease: TRANSITION_EASINGS.easeOut,
              },
            },
          },
        },
      };
    }, [sidebarWidth, sidebarPlacement, drawerMotionProps]);

    return (
      <>
        <Drawer
          ref={ref}
          {...props}
          classNames={{
            ...classNames,
            wrapper: cn("!w-[var(--sidebar-width)]", classNames?.wrapper, {
              "!items-start !justify-start ": sidebarPlacement === "left",
              "!items-end !justify-end": sidebarPlacement === "right",
            }),
            base: cn(
              "!m-0 h-full max-h-full w-[var(--sidebar-width)] p-0",
              classNames?.base,
              className,
              {
                "inset-y-0 left-0 max-h-[none] rounded-l-none !justify-start":
                  sidebarPlacement === "left",
                "inset-y-0 right-0 max-h-[none] rounded-r-none !justify-end":
                  sidebarPlacement === "right",
              },
            ),
            body: cn("p-0", classNames?.body),
            closeButton: cn("z-50", classNames?.closeButton),
          }}
          isOpen={isOpen}
          motionProps={motionProps}
          radius="none"
          scrollBehavior="inside"
          style={{
            // eslint-disable-next-line ts/ban-ts-comment
            // @ts-expect-error
            "--sidebar-width": `${sidebarWidth}px`,
          }}
          onOpenChange={onOpenChange}
        >
          <DrawerContent>
            <DrawerBody>{children}</DrawerBody>
          </DrawerContent>
        </Drawer>
      </>
    );
  },
);

SidebarDrawer.displayName = "SidebarDrawer";

export default SidebarDrawer;
