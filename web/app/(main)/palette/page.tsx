"use client"

import { BadgesSection } from "./components/badges-section"
import { ButtonsSection } from "./components/buttons-section"
import { CardsSection } from "./components/cards-section"
import { ColorsSection } from "./components/colors-section"
import { DataDisplaySection } from "./components/data-display-section"
import { InputsSection } from "./components/inputs-section"
import { NavigationSection } from "./components/navigation-section"
import { OverlaysSection } from "./components/overlays-section"
import { SelectionSection } from "./components/selection-section"

const NAV_ITEMS = [
  { id: "colors", label: "颜色" },
  { id: "buttons", label: "按钮" },
  { id: "badges", label: "徽标" },
  { id: "inputs", label: "输入" },
  { id: "selection", label: "选择" },
  { id: "cards", label: "卡片" },
  { id: "data-display", label: "数据展示" },
  { id: "navigation", label: "导航" },
  { id: "overlays", label: "弹层" },
]

export default function PalettePage() {
  return (
    <div className="flex flex-col gap-8 pb-16">
      <div className="space-y-3">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold tracking-tight">组件调色盘</h1>
          <p className="text-sm text-muted-foreground">
            预览 components/ui 下 shadcn 组件的全部样式变体，用于调整全局主题与样式。
          </p>
        </div>
        <nav className="flex flex-wrap gap-1.5">
          {NAV_ITEMS.map((item) => (
            <a
              key={item.id}
              href={`#${item.id}`}
              className="rounded-md border border-border bg-background px-2.5 py-1 text-xs text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
            >
              {item.label}
            </a>
          ))}
        </nav>
      </div>

      <ColorsSection />
      <ButtonsSection />
      <BadgesSection />
      <InputsSection />
      <SelectionSection />
      <CardsSection />
      <DataDisplaySection />
      <NavigationSection />
      <OverlaysSection />
    </div>
  )
}
