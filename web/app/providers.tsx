"use client"

import type { ThemeProviderProps } from "next-themes"
import { TooltipProvider } from "@/components/ui/tooltip"

import * as React from "react"
import { ThemeProvider } from "@/components/theme-provider"

import { Toaster } from "@/components/ui/sonner"

export interface ProvidersProps {
  children: React.ReactNode
  themeProps?: ThemeProviderProps
}

export function Providers({ children, themeProps }: ProvidersProps) {
  return (
    <ThemeProvider
      {...themeProps}
      attribute="class"
      defaultTheme="system"
      enableSystem
      disableTransitionOnChange
    >
      <TooltipProvider>{children}</TooltipProvider>
      <Toaster />
    </ThemeProvider>
  )
}
