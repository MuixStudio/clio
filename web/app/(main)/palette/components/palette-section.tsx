import * as React from "react"

import { cn } from "@/lib/utils"

export function PaletteSection({
  id,
  title,
  description,
  className,
  children,
}: {
  id: string
  title: string
  description?: string
  className?: string
  children: React.ReactNode
}) {
  return (
    <section id={id} className="scroll-mt-20 space-y-3">
      <div className="space-y-1">
        <h2 className="text-lg font-semibold tracking-tight">{title}</h2>
        {description ? (
          <p className="text-sm text-muted-foreground">{description}</p>
        ) : null}
      </div>
      <div
        className={cn(
          "flex flex-col gap-6 rounded-xl border border-border bg-card p-4",
          className
        )}
      >
        {children}
      </div>
    </section>
  )
}

export function PaletteGroup({
  label,
  className,
  children,
}: {
  label?: string
  className?: string
  children: React.ReactNode
}) {
  return (
    <div className="space-y-2.5">
      {label ? (
        <h3 className="text-xs font-medium tracking-wide text-muted-foreground uppercase">
          {label}
        </h3>
      ) : null}
      <div className={cn("flex flex-wrap items-center gap-2", className)}>
        {children}
      </div>
    </div>
  )
}