"use client"

import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@/lib/utils"


const toolbarVariants = cva(
  "flex h-(--toolbar-height) w-full items-center border-b border-border bg-background",
  {
    variants: {
      position: {
        static: "relative",
        sticky: "sticky top-0 z-50",
        fixed: "fixed top-0 right-0 left-0 z-50",
      },
      maxWidth: {
        full: "max-w-full",
        xl: "mx-auto max-w-7xl",
        lg: "mx-auto max-w-5xl",
        md: "mx-auto max-w-3xl",
      },
      blur: {
        none: "",
        sm: "bg-background/80 backdrop-blur-sm",
        md: "bg-background/80 backdrop-blur-md",
        lg: "bg-background/80 backdrop-blur-lg",
      },
    },
    defaultVariants: {
      position: "static",
      maxWidth: "full",
      blur: "none",
    },
  }
)

function Toolbar({
  className,
  position,
  maxWidth,
  blur,
  children,
  ...props
}: React.ComponentProps<"header"> & VariantProps<typeof toolbarVariants>) {
  return (
    <header
      className={cn(toolbarVariants({ position, maxWidth, blur }), className)}
      data-slot="toolbar"
      {...props}
    >
      <div className="flex h-full w-full items-center gap-2 px-4 md:px-6">
        {children}
      </div>
    </header>
  )
}


const toolbarContentVariants = cva("flex items-center gap-1", {
  variants: {
    justify: {
      start: "justify-start",
      center: "justify-center",
      end: "ml-auto justify-end",
    },
  },
  defaultVariants: {
    justify: "start",
  },
})

function ToolbarContent({
  className,
  justify,
  children,
  ...props
}: React.ComponentProps<"div"> & VariantProps<typeof toolbarContentVariants>) {
  return (
    <div
      className={cn(toolbarContentVariants({ justify }), className)}
      data-slot="toolbar-content"
      {...props}
    >
      {children}
    </div>
  )
}

const toolbarItemVariants = cva(
  "flex items-center text-sm font-medium text-foreground/80 transition-colors hover:text-foreground",
  {
    variants: {
      isActive: {
        true: "text-foreground",
        false: "",
      },
    },
    defaultVariants: {
      isActive: false,
    },
  }
)

function ToolbarItem({
  className,
  asChild = false,
  isActive = false,
  ...props
}: React.ComponentProps<"div"> & {
  asChild?: boolean
  isActive?: boolean
}) {
  const Comp = asChild ? Slot : "div"

  return (
    <Comp
      className={cn(toolbarItemVariants({ isActive }), className)}
      data-active={isActive}
      data-slot="toolbar-item"
      {...props}
    />
  )
}

function ToolbarLink({
  className,
  asChild = false,
  isActive = false,
  ...props
}: React.ComponentProps<"a"> & {
  asChild?: boolean
  isActive?: boolean
}) {
  const Comp = asChild ? Slot : "a"

  return (
    <Comp
      className={cn(
        "flex h-9 items-center rounded-md px-3 text-sm font-medium text-foreground/80 transition-colors hover:bg-accent hover:text-foreground focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none",
        isActive && "bg-accent text-foreground",
        className
      )}
      data-active={isActive}
      data-slot="toolbar-link"
      {...props}
    />
  )
}

function ToolbarSeparator({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn("my-2 h-px bg-border", className)}
      data-slot="toolbar-separator"
      {...props}
    />
  )
}

function ToolbarSpacer({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn("flex-1", className)}
      data-slot="toolbar-spacer"
      {...props}
    />
  )
}

export {
  Toolbar,
  ToolbarContent,
  ToolbarItem,
  ToolbarLink,
  ToolbarSeparator,
  ToolbarSpacer,
}
