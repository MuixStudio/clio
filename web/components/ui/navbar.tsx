"use client"

import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { cva, type VariantProps } from "class-variance-authority"
import { MenuIcon, XIcon } from "lucide-react"

import { useIsMobile } from "@/hooks/use-mobile"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet"

type NavbarContextProps = {
  isMenuOpen: boolean
  setIsMenuOpen: (open: boolean) => void
  toggleMenu: () => void
  isMobile: boolean
}

const NavbarContext = React.createContext<NavbarContextProps | null>(null)

function useNavbar() {
  const context = React.useContext(NavbarContext)

  if (!context) {
    throw new Error("useNavbar must be used within a NavbarProvider.")
  }

  return context
}

function NavbarProvider({
  children,
  className,
  ...props
}: React.ComponentProps<"div">) {
  const isMobile = useIsMobile()
  const [isMenuOpen, setIsMenuOpen] = React.useState(false)

  const toggleMenu = React.useCallback(() => {
    setIsMenuOpen((open) => !open)
  }, [])

  // Close menu when switching from mobile to desktop
  React.useEffect(() => {
    if (!isMobile && isMenuOpen) {
      setIsMenuOpen(false)
    }
  }, [isMobile, isMenuOpen])

  const contextValue = React.useMemo<NavbarContextProps>(
    () => ({
      isMenuOpen,
      setIsMenuOpen,
      toggleMenu,
      isMobile,
    }),
    [isMenuOpen, setIsMenuOpen, toggleMenu, isMobile]
  )

  return (
    <NavbarContext.Provider value={contextValue}>
      <div
        className={cn("flex relative h-(--navbar-height) w-full", className)}
        data-slot="navbar-wrapper"
        {...props}
      >
        {children}
      </div>
    </NavbarContext.Provider>
  )
}

const navbarVariants = cva(
  "flex h-(--navbar-height) w-full items-center border-b border-border bg-background",
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

function Navbar({
  className,
  position,
  maxWidth,
  blur,
  children,
  ...props
}: React.ComponentProps<"header"> & VariantProps<typeof navbarVariants>) {
  return (
    <header
      className={cn(navbarVariants({ position, maxWidth, blur }), className)}
      data-slot="navbar"
      {...props}
    >
      <div className="flex h-full bg-navbar w-full items-center gap-4 px-4 md:px-6">
        {children}
      </div>
    </header>
  )
}

function NavbarBrand({
  className,
  asChild = false,
  ...props
}: React.ComponentProps<"div"> & { asChild?: boolean }) {
  const Comp = asChild ? Slot : "div"

  return (
    <Comp
      className={cn("flex shrink-0 items-center gap-2", className)}
      data-slot="navbar-brand"
      {...props}
    />
  )
}

const navbarContentVariants = cva("flex items-center gap-1", {
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

function NavbarContent({
  className,
  justify,
  children,
  ...props
}: React.ComponentProps<"div"> & VariantProps<typeof navbarContentVariants>) {
  return (
    <div
      className={cn(navbarContentVariants({ justify }), className)}
      data-slot="navbar-content"
      {...props}
    >
      {children}
    </div>
  )
}

const navbarItemVariants = cva(
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

function NavbarItem({
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
      className={cn(navbarItemVariants({ isActive }), className)}
      data-active={isActive}
      data-slot="navbar-item"
      {...props}
    />
  )
}

function NavbarLink({
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
      data-slot="navbar-link"
      {...props}
    />
  )
}

function NavbarToggle({
  className,
  ...props
}: React.ComponentProps<typeof Button>) {
  const { toggleMenu, isMenuOpen } = useNavbar()

  return (
    <Button
      aria-label={isMenuOpen ? "Close menu" : "Open menu"}
      className={cn("size-9 md:hidden", className)}
      data-slot="navbar-toggle"
      size="icon"
      variant="ghost"
      onClick={toggleMenu}
      {...props}
    >
      {isMenuOpen ? (
        <XIcon className="size-5" />
      ) : (
        <MenuIcon className="size-5" />
      )}
    </Button>
  )
}

function NavbarMenu({
  className,
  children,
  side = "left",
  ...props
}: React.ComponentProps<"div"> & {
  side?: "left" | "right"
}) {
  const { isMenuOpen, setIsMenuOpen, isMobile } = useNavbar()

  if (!isMobile) {
    return null
  }

  return (
    <Sheet open={isMenuOpen} onOpenChange={setIsMenuOpen}>
      <SheetContent
        className={cn("w-80 p-0", className)}
        data-slot="navbar-menu"
        side={side}
        {...props}
      >
        <SheetHeader className="sr-only">
          <SheetTitle>Navigation Menu</SheetTitle>
          <SheetDescription>Mobile navigation menu</SheetDescription>
        </SheetHeader>
        <div className="flex h-full flex-col overflow-y-auto">{children}</div>
      </SheetContent>
    </Sheet>
  )
}

function NavbarMenuHeader({
  className,
  ...props
}: React.ComponentProps<"div">) {
  return (
    <div
      className={cn(
        "flex h-(--navbar-height) shrink-0 items-center border-b border-border px-4",
        className
      )}
      data-slot="navbar-menu-header"
      {...props}
    />
  )
}

function NavbarMenuContent({
  className,
  ...props
}: React.ComponentProps<"div">) {
  return (
    <div
      className={cn("flex flex-1 flex-col gap-1 p-4", className)}
      data-slot="navbar-menu-content"
      {...props}
    />
  )
}

function NavbarMenuFooter({
  className,
  ...props
}: React.ComponentProps<"div">) {
  return (
    <div
      className={cn("mt-auto border-t border-border p-4", className)}
      data-slot="navbar-menu-footer"
      {...props}
    />
  )
}

const navbarMenuItemVariants = cva(
  "flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-foreground/80 transition-colors hover:bg-accent hover:text-foreground focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none",
  {
    variants: {
      isActive: {
        true: "bg-accent text-foreground",
        false: "",
      },
    },
    defaultVariants: {
      isActive: false,
    },
  }
)

function NavbarMenuItem({
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
      className={cn(navbarMenuItemVariants({ isActive }), className)}
      data-active={isActive}
      data-slot="navbar-menu-item"
      {...props}
    />
  )
}

function NavbarMenuGroup({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn("flex flex-col gap-1", className)}
      data-slot="navbar-menu-group"
      {...props}
    />
  )
}

function NavbarMenuGroupLabel({
  className,
  ...props
}: React.ComponentProps<"div">) {
  return (
    <div
      className={cn(
        "px-3 py-2 text-xs font-semibold tracking-wider text-muted-foreground uppercase",
        className
      )}
      data-slot="navbar-menu-group-label"
      {...props}
    />
  )
}

function NavbarSeparator({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn("my-2 h-px bg-border", className)}
      data-slot="navbar-separator"
      {...props}
    />
  )
}

function NavbarSpacer({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn("flex-1", className)}
      data-slot="navbar-spacer"
      {...props}
    />
  )
}

export {
  Navbar,
  NavbarBrand,
  NavbarContent,
  NavbarItem,
  NavbarLink,
  NavbarMenu,
  NavbarMenuContent,
  NavbarMenuFooter,
  NavbarMenuGroup,
  NavbarMenuGroupLabel,
  NavbarMenuHeader,
  NavbarMenuItem,
  NavbarProvider,
  NavbarSeparator,
  NavbarSpacer,
  NavbarToggle,
  useNavbar,
}
