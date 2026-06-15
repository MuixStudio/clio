import { Bell, MoreHorizontal } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

import { PaletteGroup, PaletteSection } from "./palette-section"

export function CardsSection() {
  return (
    <PaletteSection
      id="cards"
      title="卡片 Card"
      description="components/ui/card.tsx — size: default / sm"
    >
      <PaletteGroup label="size: default" className="items-stretch">
        <Card className="w-80">
          <CardHeader>
            <CardTitle>团队通知</CardTitle>
            <CardDescription>查看最近的团队活动与告警</CardDescription>
            <CardAction>
              <Button variant="ghost" size="icon-sm" aria-label="more">
                <MoreHorizontal />
              </Button>
            </CardAction>
          </CardHeader>
          <CardContent className="py-4 text-sm text-muted-foreground">
            3 条新告警，2 个连接器状态变更
          </CardContent>
          <CardFooter className="justify-end gap-2">
            <Button variant="outline" size="sm">
              忽略
            </Button>
            <Button size="sm">查看</Button>
          </CardFooter>
        </Card>
      </PaletteGroup>

      <PaletteGroup label="size: sm" className="items-stretch">
        <Card size="sm" className="w-80">
          <CardHeader className="bg-red-100">
            <CardTitle>系统通知</CardTitle>
            <CardDescription>更紧凑的卡片样式</CardDescription>
            <CardAction>
              <Bell className="size-4 text-muted-foreground" />
            </CardAction>
          </CardHeader>
          <CardContent className="bg-blue-50 text-sm text-muted-foreground">
            紧凑布局适用于侧边栏与列表场景
          </CardContent>
          <CardContent className="text-sm text-muted-foreground">
            紧凑布局适用于侧边栏与列表场景
          </CardContent>
        </Card>
      </PaletteGroup>

      <PaletteGroup label="content only" className="items-stretch">
        <Card className="w-80">
          <CardContent className="py-4 text-sm text-muted-foreground">
            没有 Header / Footer 的简单卡片
          </CardContent>
        </Card>
      </PaletteGroup>
    </PaletteSection>
  )
}
