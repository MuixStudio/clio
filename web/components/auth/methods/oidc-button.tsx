"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { getUserOIDCSSOUrl } from "@/service/auth"
import type { OIDCMethod } from "../types"

export function OidcButton({ method }: { method: OIDCMethod }) {
  const [loading, setLoading] = useState(false)
  const Icon = method.icon

  async function handleClick() {
    setLoading(true)
    try {
      const res = await getUserOIDCSSOUrl("oidc", method.provider.id)
      window.location.href = res?.redirect_url
    } catch {
      setLoading(false)
    }
  }

  return (
    <Button
      variant="outline"
      className="w-full gap-2"
      onClick={handleClick}
      disabled={loading}
    >
      <Icon size={16} />
      Continue with {method.provider.name}
    </Button>
  )
}
