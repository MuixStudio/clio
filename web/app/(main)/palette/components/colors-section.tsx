import { cn } from "@/lib/utils"

import { PaletteSection } from "./palette-section"

const COLORS: { name: string; bg: string; fg: string }[] = [
  { name: "background", bg: "bg-background", fg: "text-foreground" },
  { name: "foreground", bg: "bg-foreground", fg: "text-background" },
  { name: "card", bg: "bg-card", fg: "text-card-foreground" },
  { name: "popover", bg: "bg-popover", fg: "text-popover-foreground" },
  { name: "primary", bg: "bg-primary", fg: "text-primary-foreground" },
  { name: "secondary", bg: "bg-secondary", fg: "text-secondary-foreground" },
  { name: "muted", bg: "bg-muted", fg: "text-muted-foreground" },
  { name: "accent", bg: "bg-accent", fg: "text-accent-foreground" },
  { name: "destructive", bg: "bg-destructive", fg: "text-white" },
  { name: "success", bg: "bg-success", fg: "text-success-foreground" },
  { name: "warning", bg: "bg-warning", fg: "text-warning-foreground" },
  { name: "info", bg: "bg-info", fg: "text-info-foreground" },
  { name: "border", bg: "bg-border", fg: "text-foreground" },
  { name: "input", bg: "bg-input", fg: "text-foreground" },
  { name: "ring", bg: "bg-ring", fg: "text-foreground" },
]

export function ColorsSection() {
  return (
    <PaletteSection
      id="colors"
      title="颜色 Colors"
      description="styles/globals.css 中 @theme 定义的全局色彩变量"
    >
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-5">
        {COLORS.map((color) => (
          <div
            key={color.name}
            className={cn(
              "flex h-20 flex-col justify-between rounded-lg border border-border p-2.5 text-xs font-medium",
              color.bg,
              color.fg
            )}
          >
            <span>{color.name}</span>
            <span className="font-mono opacity-70">--{color.name}</span>
          </div>
        ))}
      </div>
    </PaletteSection>
  )
}