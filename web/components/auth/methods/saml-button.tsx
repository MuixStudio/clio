"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { getUserOIDCSSOUrl } from "@/service/auth"
import type { SAMLMethod } from "../types"

export function SAMLButton({ method }: { method: SAMLMethod }) {
  const [loading, setLoading] = useState(false)

  // async function handleClick() {
  //   setLoading(true)
  //   try {
  //     const { authorization_url } = await getUserOIDCSSOUrl(
  //       "oidc",
  //       method.provider.id
  //     )
  //     window.location.href = authorization_url
  //     // setLoading(false)
  //   } finally {
  //     setLoading(false)
  //   }
  // }

  return (
    <Button
      className="w-full gap-2 bg-foreground text-background"
      // onClick={handleClick}
      disabled={loading}
    >
      Continue with SSO
    </Button>
  )
}