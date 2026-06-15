"use client"

import { useId } from "react"
import { Monitor, Moon, Sun } from "lucide-react"
import { useTheme } from "next-themes"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip"

const items = [
  { value: "light", label: "Light", icon: Sun },
  { value: "system", label: "System", icon: Monitor },
  { value: "dark", label: "Dark", icon: Moon },
]

export function ThemeSwitcher({ className }: { className?: string }) {
  const id = useId()
  const { resolvedTheme, setTheme } = useTheme()

  return (
    <RadioGroup
      className={`flex gap-0.5 ${className ?? ""}`}
      value={resolvedTheme ?? "system"}
      onValueChange={setTheme}
    >
      {items.map((item) => {
        const Icon = item.icon
        const isActive = (resolvedTheme ?? "system") === item.value
        return (
          <Tooltip key={`${id}-${item.value}`}>
            <TooltipTrigger asChild>
              <label
                className={`relative flex cursor-pointer items-center justify-center rounded-md p-1.5 transition-colors ${
                  isActive
                    ? "bg-primary text-primary-foreground"
                    : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
                }`}
              >
                <RadioGroupItem
                  id={`${id}-${item.value}`}
                  value={item.value}
                  className="sr-only"
                ></RadioGroupItem>
                <Icon className="size-4"></Icon>
              </label>
            </TooltipTrigger>
            <TooltipContent side="bottom">{item.label}</TooltipContent>
          </Tooltip>
        )
      })}
    </RadioGroup>
  )
}
