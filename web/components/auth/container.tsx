"use client"

import type { ReactNode } from "react"
import { cn } from "@/lib/utils"
import { PrivacyNotice } from "@/components/auth/_privacy-notice"

type ContainerProps = {
  logo?: ReactNode | string
  title?: string
  description?: string
  children: ReactNode
  accountAgreementUrl?: string
  className?: string
}

function Logo({ value }: { value: ReactNode | string }) {
  if (typeof value === "string") {
    if (value.startsWith("http") || value.startsWith("/")) {
      return <img src={value} alt="" className="size-10 object-contain" />
    }
    return <span className="text-sm font-medium">{value}</span>
  }
  return <>{value}</>
}

export function Container({
  logo,
  title,
  description,
  children,
  className,
  accountAgreementUrl = "/account-agreement",
}: ContainerProps) {
  return (
    <div className={cn("w-full max-w-md p-6", className)}>
      <div className="mb-6 flex flex-col items-center gap-2 text-center">
        {logo && (
          <div className="mb-1">
            <Logo value={logo} />
          </div>
        )}
        {title && (
          <h1 className="text-xl font-semibold tracking-tight">{title}</h1>
        )}
        {description && (
          <p className="text-sm text-muted-foreground">{description}</p>
        )}
      </div>

      <div className="flex flex-col gap-4">{children}</div>

      <PrivacyNotice accountAgreementUrl={accountAgreementUrl} />
    </div>
  )
}
