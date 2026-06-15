import { CheckCircle2, Sparkles } from "lucide-react"

import { Badge } from "@/components/ui/badge"

import { PaletteGroup, PaletteSection } from "./palette-section"

const VARIANTS = [
  "default",
  "secondary",
  "destructive",
  "outline",
  "ghost",
  "link",
] as const

export function BadgesSection() {
  return (
    <PaletteSection
      id="badges"
      title="徽标 Badge"
      description="components/ui/badge.tsx — variant"
    >
      <PaletteGroup label="variants">
        {VARIANTS.map((variant) => (
          <Badge key={variant} variant={variant}>
            {variant}
          </Badge>
        ))}
      </PaletteGroup>
      <PaletteGroup label="with icon">
        <Badge>
          <Sparkles /> New
        </Badge>
        <Badge variant="secondary">
          <CheckCircle2 /> Active
        </Badge>
        <Badge variant="destructive">Error</Badge>
        <Badge variant="outline">Outline</Badge>
      </PaletteGroup>
    </PaletteSection>
  )
}