"use client"

import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { Slot } from "radix-ui"

import { cn } from "@/lib/utils"
import { Input } from "@/components/ui/input"
import { Separator } from "@/components/ui/separator"
import { Skeleton } from "@/components/ui/skeleton"
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip"

const SECONDARY_SIDEBAR_MIN_WIDTH = 158
const SECONDARY_SIDEBAR_MAX_WIDTH = 268
const SECONDARY_SIDEBAR_DEFAULT_WIDTH = 222
const SECONDARY_SIDEBAR_KEYBOARD_STEP = 16

type SecondarySidebarContextProps = {
  width: number
  minWidth: number
  maxWidth: number
  isResizing: boolean
  startResizing: (event: React.PointerEvent<HTMLButtonElement>) => void
  resizeByKeyboard: (event: React.KeyboardEvent<HTMLButtonElement>) => void
}

const SecondarySidebarContext =
  React.createContext<SecondarySidebarContextProps | null>(null)

function useSecondarySidebar() {
  const context = React.useContext(SecondarySidebarContext)

  if (!context) {
    throw new Error(
      "useSecondarySidebar must be used within a SecondarySidebarProvider."
    )
  }

  return context
}

type SecondarySidebarProviderProps = React.ComponentProps<"div"> & {
  defaultWidth?: number
  minWidth?: number
  maxWidth?: number
}

function clampWidth(width: number, minWidth: number, maxWidth: number) {
  return Math.min(maxWidth, Math.max(minWidth, width))
}

function SecondarySidebarProvider({
  defaultWidth = SECONDARY_SIDEBAR_DEFAULT_WIDTH,
  minWidth = SECONDARY_SIDEBAR_MIN_WIDTH,
  maxWidth = SECONDARY_SIDEBAR_MAX_WIDTH,
  className,
  style,
  children,
  ...props
}: SecondarySidebarProviderProps) {
  const [width, setWidth] = React.useState(() =>
    clampWidth(defaultWidth, minWidth, maxWidth)
  )
  const [isResizing, setIsResizing] = React.useState(false)

  const widthRef = React.useRef(width)
  const startXRef = React.useRef(0)
  const startWidthRef = React.useRef(width)

  React.useEffect(() => {
    widthRef.current = width
  }, [width])

  React.useEffect(() => {
    if (!isResizing) {
      return
    }

    const handlePointerMove = (event: PointerEvent) => {
      const nextWidth = startWidthRef.current + event.clientX - startXRef.current
      setWidth(clampWidth(nextWidth, minWidth, maxWidth))
    }

    const stopResizing = () => {
      setIsResizing(false)
      document.body.style.removeProperty("cursor")
      document.body.style.removeProperty("user-select")
    }

    document.body.style.cursor = "col-resize"
    document.body.style.userSelect = "none"
    window.addEventListener("pointermove", handlePointerMove)
    window.addEventListener("pointerup", stopResizing, { once: true })
    window.addEventListener("pointercancel", stopResizing, { once: true })

    return () => {
      document.body.style.removeProperty("cursor")
      document.body.style.removeProperty("user-select")
      window.removeEventListener("pointermove", handlePointerMove)
      window.removeEventListener("pointerup", stopResizing)
      window.removeEventListener("pointercancel", stopResizing)
    }
  }, [isResizing, minWidth, maxWidth])

  const startResizing = React.useCallback(
    (event: React.PointerEvent<HTMLButtonElement>) => {
      event.preventDefault()
      startXRef.current = event.clientX
      startWidthRef.current = widthRef.current
      setIsResizing(true)
    },
    []
  )

  const resizeByKeyboard = React.useCallback(
    (event: React.KeyboardEvent<HTMLButtonElement>) => {
      if (event.key === "ArrowLeft") {
        event.preventDefault()
        setWidth((currentWidth) =>
          clampWidth(
            currentWidth - SECONDARY_SIDEBAR_KEYBOARD_STEP,
            minWidth,
            maxWidth
          )
        )
      }

      if (event.key === "ArrowRight") {
        event.preventDefault()
        setWidth((currentWidth) =>
          clampWidth(
            currentWidth + SECONDARY_SIDEBAR_KEYBOARD_STEP,
            minWidth,
            maxWidth
          )
        )
      }
    },
    [minWidth, maxWidth]
  )

  const contextValue = React.useMemo<SecondarySidebarContextProps>(
    () => ({
      width,
      minWidth,
      maxWidth,
      isResizing,
      startResizing,
      resizeByKeyboard,
    }),
    [width, minWidth, maxWidth, isResizing, startResizing, resizeByKeyboard]
  )

  return (
    <SecondarySidebarContext.Provider value={contextValue}>
      <div
        data-slot="secondary-sidebar-wrapper"
        className={cn("flex h-full min-h-0 w-full overflow-hidden", className)}
        style={style}
        {...props}
      >
        {children}
      </div>
    </SecondarySidebarContext.Provider>
  )
}

function SecondarySidebar({
  className,
  style,
  children,
  ...props
}: React.ComponentProps<"aside">) {
  const { width, minWidth, maxWidth, isResizing } = useSecondarySidebar()

  return (
    <aside
      data-slot="secondary-sidebar"
      className={cn(
        "relative flex h-full shrink-0 flex-col border-r border-sidebar-border bg-sidebar text-sidebar-foreground",
        isResizing && "select-none",
        className
      )}
      style={{
        width,
        minWidth,
        maxWidth,
        ...style,
      }}
      {...props}
    >
      {children}
    </aside>
  )
}

function SecondarySidebarInset({
  className,
  ...props
}: React.ComponentProps<"main">) {
  return (
    <main
      data-slot="secondary-sidebar-inset"
      className={cn(
        "relative flex min-w-0 flex-1 flex-col overflow-hidden bg-background",
        className
      )}
      {...props}
    />
  )
}

function SecondarySidebarInput({
  className,
  ...props
}: React.ComponentProps<typeof Input>) {
  return (
    <Input
      data-slot="secondary-sidebar-input"
      data-secondary-sidebar="input"
      className={cn("h-8 w-full bg-background shadow-none", className)}
      {...props}
    />
  )
}

function SecondarySidebarHeader({
  className,
  ...props
}: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="secondary-sidebar-header"
      data-secondary-sidebar="header"
      className={cn(
        "flex flex-col gap-2 border-b border-sidebar-border/80 p-3",
        className
      )}
      {...props}
    />
  )
}

function SecondarySidebarFooter({
  className,
  ...props
}: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="secondary-sidebar-footer"
      data-secondary-sidebar="footer"
      className={cn(
        "flex flex-col gap-2 border-t border-sidebar-border/80 p-3",
        className
      )}
      {...props}
    />
  )
}

function SecondarySidebarSeparator({
  className,
  ...props
}: React.ComponentProps<typeof Separator>) {
  return (
    <Separator
      data-slot="secondary-sidebar-separator"
      data-secondary-sidebar="separator"
      className={cn("mx-2 w-auto bg-sidebar-border", className)}
      {...props}
    />
  )
}

function SecondarySidebarContent({
  className,
  ...props
}: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="secondary-sidebar-content"
      data-secondary-sidebar="content"
      className={cn(
        "no-scrollbar flex min-h-0 flex-1 flex-col gap-0 overflow-auto bg-[radial-gradient(circle_at_20%_0%,hsl(var(--sidebar-accent)/0.55),transparent_28rem)]",
        className
      )}
      {...props}
    />
  )
}

function SecondarySidebarGroup({
  className,
  ...props
}: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="secondary-sidebar-group"
      data-secondary-sidebar="group"
      className={cn("relative flex w-full min-w-0 flex-col p-2", className)}
      {...props}
    />
  )
}

function SecondarySidebarGroupLabel({
  className,
  asChild = false,
  ...props
}: React.ComponentProps<"div"> & { asChild?: boolean }) {
  const Comp = asChild ? Slot.Root : "div"

  return (
    <Comp
      data-slot="secondary-sidebar-group-label"
      data-secondary-sidebar="group-label"
      className={cn(
        "flex h-8 shrink-0 items-center rounded-md px-2 text-xs font-medium text-sidebar-foreground/70 ring-sidebar-ring outline-hidden transition-colors focus-visible:ring-2 [&>svg]:size-4 [&>svg]:shrink-0",
        className
      )}
      {...props}
    />
  )
}

function SecondarySidebarGroupAction({
  className,
  asChild = false,
  ...props
}: React.ComponentProps<"button"> & { asChild?: boolean }) {
  const Comp = asChild ? Slot.Root : "button"

  return (
    <Comp
      data-slot="secondary-sidebar-group-action"
      data-secondary-sidebar="group-action"
      className={cn(
        "absolute top-3.5 right-3 flex aspect-square w-5 items-center justify-center rounded-md p-0 text-sidebar-foreground ring-sidebar-ring outline-hidden transition-colors after:absolute after:-inset-2 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 md:after:hidden [&>svg]:size-4 [&>svg]:shrink-0",
        className
      )}
      {...props}
    />
  )
}

function SecondarySidebarGroupContent({
  className,
  ...props
}: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="secondary-sidebar-group-content"
      data-secondary-sidebar="group-content"
      className={cn("w-full text-sm", className)}
      {...props}
    />
  )
}

function SecondarySidebarMenu({
  className,
  ...props
}: React.ComponentProps<"ul">) {
  return (
    <ul
      data-slot="secondary-sidebar-menu"
      data-secondary-sidebar="menu"
      className={cn("flex w-full min-w-0 flex-col gap-0", className)}
      {...props}
    />
  )
}

function SecondarySidebarMenuItem({
  className,
  ...props
}: React.ComponentProps<"li">) {
  return (
    <li
      data-slot="secondary-sidebar-menu-item"
      data-secondary-sidebar="menu-item"
      className={cn("group/menu-item relative", className)}
      {...props}
    />
  )
}

const secondarySidebarMenuButtonVariants = cva(
  "peer/menu-button group/menu-button flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left text-sm ring-sidebar-ring outline-hidden transition-[width,height,padding,background-color,color] group-has-data-[secondary-sidebar=menu-action]/menu-item:pr-8 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 data-open:hover:bg-sidebar-accent data-open:hover:text-sidebar-accent-foreground data-active:bg-sidebar-accent data-active:font-medium data-active:text-sidebar-accent-foreground [&_svg]:size-4 [&_svg]:shrink-0 [&>span:last-child]:truncate",
  {
    variants: {
      variant: {
        default: "hover:bg-sidebar-accent hover:text-sidebar-accent-foreground",
        outline:
          "bg-background shadow-[0_0_0_1px_hsl(var(--sidebar-border))] hover:bg-sidebar-accent hover:text-sidebar-accent-foreground hover:shadow-[0_0_0_1px_hsl(var(--sidebar-accent))]",
      },
      size: {
        default: "h-8 text-sm",
        sm: "h-7 text-xs",
        lg: "h-12 text-sm",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
)

function SecondarySidebarMenuButton({
  asChild = false,
  isActive = false,
  variant = "default",
  size = "default",
  tooltip,
  className,
  ...props
}: React.ComponentProps<"button"> & {
  asChild?: boolean
  isActive?: boolean
  tooltip?: string | React.ComponentProps<typeof TooltipContent>
} & VariantProps<typeof secondarySidebarMenuButtonVariants>) {
  const Comp = asChild ? Slot.Root : "button"

  const button = (
    <Comp
      data-slot="secondary-sidebar-menu-button"
      data-secondary-sidebar="menu-button"
      data-size={size}
      data-active={isActive}
      className={cn(
        secondarySidebarMenuButtonVariants({ variant, size }),
        className
      )}
      {...props}
    />
  )

  if (!tooltip) {
    return button
  }

  if (typeof tooltip === "string") {
    tooltip = {
      children: tooltip,
    }
  }

  return (
    <Tooltip>
      <TooltipTrigger asChild>{button}</TooltipTrigger>
      <TooltipContent side="right" align="center" {...tooltip} />
    </Tooltip>
  )
}

function SecondarySidebarMenuAction({
  className,
  asChild = false,
  showOnHover = false,
  ...props
}: React.ComponentProps<"button"> & {
  asChild?: boolean
  showOnHover?: boolean
}) {
  const Comp = asChild ? Slot.Root : "button"

  return (
    <Comp
      data-slot="secondary-sidebar-menu-action"
      data-secondary-sidebar="menu-action"
      className={cn(
        "absolute top-1.5 right-1 flex aspect-square w-5 items-center justify-center rounded-md p-0 text-sidebar-foreground ring-sidebar-ring outline-hidden transition-colors peer-hover/menu-button:text-sidebar-accent-foreground peer-data-[size=default]/menu-button:top-1.5 peer-data-[size=lg]/menu-button:top-2.5 peer-data-[size=sm]/menu-button:top-1 after:absolute after:-inset-2 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 md:after:hidden [&>svg]:size-4 [&>svg]:shrink-0",
        showOnHover &&
          "group-focus-within/menu-item:opacity-100 group-hover/menu-item:opacity-100 peer-data-active/menu-button:text-sidebar-accent-foreground aria-expanded:opacity-100 md:opacity-0",
        className
      )}
      {...props}
    />
  )
}

function SecondarySidebarMenuBadge({
  className,
  ...props
}: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="secondary-sidebar-menu-badge"
      data-secondary-sidebar="menu-badge"
      className={cn(
        "pointer-events-none absolute right-1 flex h-5 min-w-5 items-center justify-center rounded-md px-1 text-xs font-medium text-sidebar-foreground tabular-nums select-none peer-hover/menu-button:text-sidebar-accent-foreground peer-data-[size=default]/menu-button:top-1.5 peer-data-[size=lg]/menu-button:top-2.5 peer-data-[size=sm]/menu-button:top-1 peer-data-active/menu-button:text-sidebar-accent-foreground",
        className
      )}
      {...props}
    />
  )
}

function SecondarySidebarMenuSkeleton({
  className,
  showIcon = false,
  ...props
}: React.ComponentProps<"div"> & {
  showIcon?: boolean
}) {
  const [width] = React.useState(() => {
    return `${Math.floor(Math.random() * 40) + 50}%`
  })

  return (
    <div
      data-slot="secondary-sidebar-menu-skeleton"
      data-secondary-sidebar="menu-skeleton"
      className={cn("flex h-8 items-center gap-2 rounded-md px-2", className)}
      {...props}
    >
      {showIcon && (
        <Skeleton
          className="size-4 rounded-md"
          data-secondary-sidebar="menu-skeleton-icon"
        />
      )}
      <Skeleton
        className="h-4 max-w-(--skeleton-width) flex-1"
        data-secondary-sidebar="menu-skeleton-text"
        style={
          {
            "--skeleton-width": width,
          } as React.CSSProperties
        }
      />
    </div>
  )
}

function SecondarySidebarMenuSub({
  className,
  ...props
}: React.ComponentProps<"ul">) {
  return (
    <ul
      data-slot="secondary-sidebar-menu-sub"
      data-secondary-sidebar="menu-sub"
      className={cn(
        "mx-3.5 flex min-w-0 translate-x-px flex-col gap-1 border-l border-sidebar-border px-2.5 py-0.5",
        className
      )}
      {...props}
    />
  )
}

function SecondarySidebarMenuSubItem({
  className,
  ...props
}: React.ComponentProps<"li">) {
  return (
    <li
      data-slot="secondary-sidebar-menu-sub-item"
      data-secondary-sidebar="menu-sub-item"
      className={cn("group/menu-sub-item relative", className)}
      {...props}
    />
  )
}

function SecondarySidebarMenuSubButton({
  asChild = false,
  size = "md",
  isActive = false,
  className,
  ...props
}: React.ComponentProps<"a"> & {
  asChild?: boolean
  size?: "sm" | "md"
  isActive?: boolean
}) {
  const Comp = asChild ? Slot.Root : "a"

  return (
    <Comp
      data-slot="secondary-sidebar-menu-sub-button"
      data-secondary-sidebar="menu-sub-button"
      data-size={size}
      data-active={isActive}
      className={cn(
        "flex h-7 min-w-0 -translate-x-px items-center gap-2 overflow-hidden rounded-md px-2 text-sidebar-foreground ring-sidebar-ring outline-hidden hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 data-[size=md]:text-sm data-[size=sm]:text-xs data-active:bg-sidebar-accent data-active:text-sidebar-accent-foreground [&>span:last-child]:truncate [&>svg]:size-4 [&>svg]:shrink-0 [&>svg]:text-sidebar-accent-foreground",
        className
      )}
      {...props}
    />
  )
}

function SecondarySidebarResizeHandle({
  className,
  "aria-label": ariaLabel = "Resize secondary sidebar",
  ...props
}: React.ComponentProps<"button">) {
  const { isResizing, startResizing, resizeByKeyboard } = useSecondarySidebar()

  return (
    <button
      type="button"
      role="separator"
      aria-orientation="vertical"
      aria-label={ariaLabel}
      data-slot="secondary-sidebar-resize-handle"
      className={cn(
        "absolute inset-y-0 -right-1 z-20 w-2 cursor-col-resize touch-none outline-none",
        "after:absolute after:inset-y-0 after:left-1/2 after:w-px after:-translate-x-1/2 after:bg-transparent after:transition-colors",
        "hover:after:bg-sidebar-ring focus-visible:after:bg-sidebar-ring",
        isResizing && "after:bg-sidebar-ring",
        className
      )}
      onPointerDown={startResizing}
      onKeyDown={resizeByKeyboard}
      {...props}
    />
  )
}

export {
  SecondarySidebar,
  SecondarySidebarContent,
  SecondarySidebarFooter,
  SecondarySidebarGroup,
  SecondarySidebarGroupAction,
  SecondarySidebarGroupContent,
  SecondarySidebarGroupLabel,
  SecondarySidebarHeader,
  SecondarySidebarInput,
  SecondarySidebarInset,
  SecondarySidebarMenu,
  SecondarySidebarMenuAction,
  SecondarySidebarMenuBadge,
  SecondarySidebarMenuButton,
  SecondarySidebarMenuItem,
  SecondarySidebarMenuSkeleton,
  SecondarySidebarMenuSub,
  SecondarySidebarMenuSubButton,
  SecondarySidebarMenuSubItem,
  SecondarySidebarProvider,
  SecondarySidebarResizeHandle,
  SecondarySidebarSeparator,
  useSecondarySidebar,
  SECONDARY_SIDEBAR_DEFAULT_WIDTH,
  SECONDARY_SIDEBAR_MAX_WIDTH,
  SECONDARY_SIDEBAR_MIN_WIDTH,
}