import { Heart, Loader2, Mail, Plus, Trash2 } from "lucide-react"

import { Button } from "@/components/ui/button"

import { PaletteGroup, PaletteSection } from "./palette-section"

const VARIANTS = [
  "default",
  "secondary",
  "destructive",
  "outline",
  "ghost",
  "link",
] as const

const SIZES = ["lg", "default", "sm", "xs"] as const

const ICON_SIZES = ["icon-lg", "icon", "icon-sm", "icon-xs"] as const

export function ButtonsSection() {
  return (
    <PaletteSection
      id="buttons"
      title="按钮 Button"
      description="components/ui/button.tsx — variant × size"
    >
      {VARIANTS.map((variant) => (
        <PaletteGroup key={variant} label={variant}>
          {SIZES.map((size) => (
            <Button key={size} variant={variant} size={size}>
              Button
            </Button>
          ))}
          {ICON_SIZES.map((size) => (
            <Button
              key={size}
              variant={variant}
              size={size}
              aria-label="favorite"
            >
              <Heart />
            </Button>
          ))}
          <Button variant={variant} disabled>
            Disabled
          </Button>
        </PaletteGroup>
      ))}
      <PaletteGroup label="with icon">
        <Button>
          <Mail /> Email
        </Button>
        <Button variant="secondary">
          <Plus /> Add item
        </Button>
        <Button variant="destructive">
          <Trash2 /> Delete
        </Button>
        <Button variant="outline" disabled>
          <Loader2 className="animate-spin" /> Loading
        </Button>
      </PaletteGroup>
    </PaletteSection>
  )
}