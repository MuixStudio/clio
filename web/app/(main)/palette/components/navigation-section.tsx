import {
  Breadcrumb,
  BreadcrumbEllipsis,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"

import { PaletteGroup, PaletteSection } from "./palette-section"

export function NavigationSection() {
  return (
    <PaletteSection
      id="navigation"
      title="导航 Navigation"
      description="tabs.tsx, breadcrumb.tsx, pagination.tsx"
    >
      <PaletteGroup label="tabs — variant: default" className="flex-col items-stretch">
        <Tabs defaultValue="account" className="w-full max-w-md">
          <TabsList>
            <TabsTrigger value="account">账户</TabsTrigger>
            <TabsTrigger value="password">密码</TabsTrigger>
            <TabsTrigger value="team">团队</TabsTrigger>
          </TabsList>
          <TabsContent value="account" className="text-sm text-muted-foreground">
            管理账户的基本信息
          </TabsContent>
          <TabsContent value="password" className="text-sm text-muted-foreground">
            修改登录密码
          </TabsContent>
          <TabsContent value="team" className="text-sm text-muted-foreground">
            管理团队成员与权限
          </TabsContent>
        </Tabs>
      </PaletteGroup>

      <PaletteGroup label="tabs — variant: line" className="flex-col items-stretch">
        <Tabs defaultValue="overview" className="w-full max-w-md">
          <TabsList variant="line">
            <TabsTrigger value="overview">概览</TabsTrigger>
            <TabsTrigger value="alerts">告警</TabsTrigger>
            <TabsTrigger value="connector">接入源</TabsTrigger>
          </TabsList>
          <TabsContent value="overview" className="text-sm text-muted-foreground">
            团队整体运行状况
          </TabsContent>
          <TabsContent value="alerts" className="text-sm text-muted-foreground">
            最近触发的告警列表
          </TabsContent>
          <TabsContent value="connector" className="text-sm text-muted-foreground">
            已接入的告警源
          </TabsContent>
        </Tabs>
      </PaletteGroup>

      <PaletteGroup label="breadcrumb">
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem>
              <BreadcrumbLink href="#">Zus</BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbEllipsis />
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbLink href="#">Connector</BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbPage>详情</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
      </PaletteGroup>

      <PaletteGroup label="pagination">
        <Pagination className="mx-0 w-auto justify-start">
          <PaginationContent>
            <PaginationItem>
              <PaginationPrevious href="#" />
            </PaginationItem>
            <PaginationItem>
              <PaginationLink href="#">1</PaginationLink>
            </PaginationItem>
            <PaginationItem>
              <PaginationLink href="#" isActive>
                2
              </PaginationLink>
            </PaginationItem>
            <PaginationItem>
              <PaginationLink href="#">3</PaginationLink>
            </PaginationItem>
            <PaginationItem>
              <PaginationEllipsis />
            </PaginationItem>
            <PaginationItem>
              <PaginationNext href="#" />
            </PaginationItem>
          </PaginationContent>
        </Pagination>
      </PaletteGroup>
    </PaletteSection>
  )
}
