"use client";
import {
  Avatar,
  Badge,
  Dropdown,
  DropdownItem,
  DropdownMenu,
  DropdownSection,
  DropdownTrigger,
  Link,
  Navbar,
  NavbarBrand,
  NavbarContent,
  NavbarItem,
} from "@heroui/react";

import { AdmgoLogo } from "@/components/icons";

export default function Header() {
  return (
    <Navbar
      isBordered
      className="bg-foreground-100 items-start justify-start overflow-x-auto"
      height="3rem"
      maxWidth="full"
    >
      <NavbarBrand className="flex-none grow-0">
        <AdmgoLogo />
      </NavbarBrand>
      <NavbarContent className="w-full flex-1 grow">
        <NavbarItem>
          <Link href="/home">home</Link>
        </NavbarItem>
        <NavbarItem>
          <Link href="/cmdb">cmdb</Link>
        </NavbarItem>
        <NavbarItem>1.2</NavbarItem>
        <NavbarItem>1.3</NavbarItem>
      </NavbarContent>
      <NavbarContent className="grow-0 data-[justify=end]:grow-0" justify="end">
        <NavbarItem>
          <Link
            aria-current="page"
            className="flex gap-2 text-xs font-medium text-inherit"
            href="#"
          >
            Tick
          </Link>
        </NavbarItem>
        <NavbarItem>
          <Link
            aria-current="page"
            className="flex gap-2 text-xs font-medium text-inherit"
            href="#"
          >
            Open
          </Link>
        </NavbarItem>
        <NavbarItem>
          <Link
            aria-current="page"
            className="flex gap-2 text-xs font-medium text-inherit"
            href="#"
          >
            Chart
          </Link>
        </NavbarItem>
        <NavbarItem>
          <Link
            aria-current="page"
            className="flex gap-2 text-xs font-medium text-inherit"
            href="#"
          >
            Deployments
          </Link>
        </NavbarItem>
        <NavbarItem>
          <Dropdown placement="bottom-end">
            <DropdownTrigger>
              <button className="mt-1 h-8 w-8 outline-none transition-transform">
                <Badge
                  className="border-transparent"
                  color="success"
                  content=""
                  placement="bottom-right"
                  shape="circle"
                  size="sm"
                >
                  <Avatar
                    size="sm"
                    src="https://i.pravatar.cc/150?u=a04258114e29526708c"
                  />
                </Badge>
              </button>
            </DropdownTrigger>
            <DropdownMenu
              aria-label="Profile Actions"
              disabledKeys={["your"]}
              itemClasses={{
                base: [
                  "rounded-md",
                  "text-default-500",
                  "transition-opacity",
                  "data-[hover=true]:text-foreground",
                  "data-[hover=true]:bg-default-100",
                  "dark:data-[hover=true]:bg-default-50",
                  "data-[selectable=true]:focus:bg-default-50",
                  "data-[pressed=true]:opacity-70",
                  "data-[focus-visible=true]:ring-default-500",
                ],
              }}
              variant="flat"
            >
              <DropdownSection showDivider>
                <DropdownItem
                  key="your"
                  isReadOnly
                  className="h-14 gap-2 opacity-100"
                >
                  <p className="font-semibold">Signed in as</p>
                  <p className="font-semibold">me@muixstudio.com</p>
                </DropdownItem>
              </DropdownSection>
              <DropdownSection key="account" showDivider title="个人设置">
                <DropdownItem key="profile">个人资料</DropdownItem>
                <DropdownItem key="settings">设置</DropdownItem>
              </DropdownSection>
              <DropdownSection key="admin" showDivider title="管理员">
                <DropdownItem key="admin">系统设置</DropdownItem>
              </DropdownSection>
              <DropdownSection>
                <DropdownItem key="logout" color="danger">
                  注销
                </DropdownItem>
              </DropdownSection>
            </DropdownMenu>
          </Dropdown>
        </NavbarItem>
      </NavbarContent>
    </Navbar>
  );
}
