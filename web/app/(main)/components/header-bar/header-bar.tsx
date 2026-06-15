import * as React from "react"
import { Slot } from "@radix-ui/react-slot"

import { cn } from "@/lib/utils"

function HeaderBar({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn(
        "flex items-center gap-4 border-b border-border bg-background px-3 py-2",
        className
      )}
      data-slot="header-bar"
      {...props}
    />
  )
}

function HeaderBarLeft({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn("flex flex-col gap-1", className)}
      data-slot="header-bar-left"
      {...props}
    />
  )
}

function HeaderBarTitle({
  className,
  asChild = false,
  ...props
}: React.ComponentProps<"div"> & { asChild?: boolean }) {
  const Comp = asChild ? Slot : "div"
  return (
    <Comp
      className={cn("text-lg font-semibold leading-none tracking-tight", className)}
      data-slot="header-bar-title"
      {...props}
    />
  )
}

function HeaderBarDescription({
  className,
  asChild = false,
  ...props
}: React.ComponentProps<"div"> & { asChild?: boolean }) {
  const Comp = asChild ? Slot : "div"
  return (
    <Comp
      className={cn("text-xs text-muted-foreground", className)}
      data-slot="header-bar-description"
      {...props}
    />
  )
}

function HeaderBarActions({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn("ml-auto flex items-center gap-2", className)}
      data-slot="header-bar-actions"
      {...props}
    />
  )
}

function HeaderBarItem({
  className,
  asChild = false,
  ...props
}: React.ComponentProps<"div"> & { asChild?: boolean }) {
  const Comp = asChild ? Slot : "div"
  return (
    <Comp
      className={cn("flex items-center", className)}
      data-slot="header-bar-item"
      {...props}
    />
  )
}

export {
  HeaderBar,
  HeaderBarLeft,
  HeaderBarTitle,
  HeaderBarDescription,
  HeaderBarActions,
  HeaderBarItem,
}
