"use client"

import React from "react"
import { Separator } from "@/components/ui/separator"
import { Logo } from "@/components/icons"

function BrandPanel() {
  return (
    <div className="relative hidden flex-col justify-between bg-zinc-950 p-10 text-white lg:flex">
      {/* Top: logo + product name */}
      <div className="flex items-center gap-2.5">
        <Logo className="size-7 text-white" />
        <span className="text-base font-semibold tracking-tight">Clio</span>
      </div>

      {/* Center: headline */}
      <div className="space-y-4">
        <blockquote className="space-y-3">
          <p className="text-2xl font-medium leading-snug tracking-tight">
            &ldquo;Unified visibility across every layer of your infrastructure
            — from bare metal to cloud.&rdquo;
          </p>
          <Separator className="bg-white/20" />
          <footer className="text-sm text-zinc-400">
            CMDB &middot; Asset Intelligence Platform
          </footer>
        </blockquote>
      </div>

      {/* Bottom: version / copyright */}
      <p className="text-xs text-zinc-600">
        &copy; {new Date().getFullYear()} Clio. All rights reserved.
      </p>

      {/* Decorative gradient blob */}
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0 overflow-hidden"
      >
        <div className="absolute -top-32 -right-32 size-96 rounded-full bg-white/5 blur-3xl" />
        <div className="absolute -bottom-32 -left-32 size-96 rounded-full bg-white/5 blur-3xl" />
      </div>
    </div>
  )
}

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <div className="grid min-h-screen lg:grid-cols-2">
      <BrandPanel />

      {/* Right: form area */}
      <div className="flex items-center justify-center bg-background p-6">
        {children}
      </div>
    </div>
  )
}