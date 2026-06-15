"use client"

import * as React from "react"
import Link from "next/link"
import {
  CircleAlertIcon,
  HomeIcon,
  FileTextIcon,
  LayoutGridIcon,
  PanelLeftIcon,
} from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  Navbar,
  NavbarBrand,
  NavbarContent,
  NavbarItem,
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
} from "@/components/ui/navbar"
import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
  navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu"

import { NavbarUser } from "./navbar-user"
import { SidebarTrigger, useSidebar } from "@/components/ui/sidebar"
import { TeamSelect } from "@/components/zus/team-select"

export function AppNavbar() {
  const { toggleSidebar } = useSidebar()
  return (
    <NavbarProvider>
      <Navbar position="sticky" blur="md">
        {/* Sidebar Mobile Toggle */}
        <NavbarContent className="flex md:hidden" justify="start">
          <NavbarItem>
            <Button
              // className="h-8 w-8"
              variant="ghost"
              // size="icon"
              onClick={toggleSidebar}
            >
              <PanelLeftIcon />
            </Button>
          </NavbarItem>
        </NavbarContent>

        {/* Logo */}
        <NavbarBrand>
          <Link className="flex items-center gap-2 font-semibold" href="/">
            <div className="flex size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
              <HomeIcon className="size-4" />
            </div>
            <span className="hidden sm:inline-block">Clio</span>
          </Link>
        </NavbarBrand>

        {/* Desktop Navigation - Center */}
        <NavbarContent className="hidden md:flex" justify="center">
          <NavbarItem>
            <DesktopNavigation />
          </NavbarItem>
        </NavbarContent>
        <NavbarContent justify="end">
          <NavbarItem>
            <TeamSelect />
          </NavbarItem>
        </NavbarContent>
        {/* Spacer pushes actions to the right */}
        <NavbarSpacer />

        {/* Right Actions */}
        <NavbarContent justify="end">
          <NavbarItem>
            <NavbarUser />
          </NavbarItem>
        </NavbarContent>
      </Navbar>
    </NavbarProvider>
  )
}

function DesktopNavigation() {
  return (
    <NavigationMenu viewport={false}>
      <NavigationMenuList>
        <NavigationMenuItem>
          <NavigationMenuLink asChild className={navigationMenuTriggerStyle()}>
            <Link href="/docs">Documentation</Link>
          </NavigationMenuLink>
        </NavigationMenuItem>
        <NavigationMenuItem>
          <NavigationMenuTrigger>Components</NavigationMenuTrigger>
          <NavigationMenuContent>
            <ul className="w-72">
              <li>
                <NavigationMenuLink asChild>
                  <Link href="#">
                    <div className="flex flex-col">
                      <div className="font-medium">All Components</div>
                      <div className="text-muted-foreground">
                        Browse all components in the library.
                      </div>
                    </div>
                  </Link>
                </NavigationMenuLink>
                <NavigationMenuLink asChild>
                  <Link href="#">
                    <div className="flex flex-col">
                      <div className="font-medium">Primitives</div>
                      <div className="text-muted-foreground">
                        Low-level building blocks.
                      </div>
                    </div>
                  </Link>
                </NavigationMenuLink>
                <NavigationMenuLink asChild>
                  <Link href="#">
                    <div className="flex flex-col">
                      <div className="font-medium">Patterns</div>
                      <div className="text-muted-foreground">
                        Common UI patterns and layouts.
                      </div>
                    </div>
                  </Link>
                </NavigationMenuLink>
              </li>
            </ul>
          </NavigationMenuContent>
        </NavigationMenuItem>
        <NavigationMenuItem>
          <NavigationMenuTrigger>Resources</NavigationMenuTrigger>
          <NavigationMenuContent>
            <ul className="grid w-50">
              <li>
                <NavigationMenuLink asChild>
                  <Link className="flex-row items-center gap-2" href="#">
                    <CircleAlertIcon />
                    Blog
                  </Link>
                </NavigationMenuLink>
                <NavigationMenuLink asChild>
                  <Link className="flex-row items-center gap-2" href="#">
                    <CircleAlertIcon />
                    Changelog
                  </Link>
                </NavigationMenuLink>
                <NavigationMenuLink asChild>
                  <Link href="#">
                    <CircleAlertIcon />
                    GitHub
                  </Link>
                </NavigationMenuLink>
              </li>
            </ul>
          </NavigationMenuContent>
        </NavigationMenuItem>
      </NavigationMenuList>
    </NavigationMenu>
  )
}
