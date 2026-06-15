"use client"

import { Suspense, useEffect } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { exchangeCode } from "@/service/auth"

function AuthCallbackContent() {
  const router = useRouter()
  const params = useSearchParams()

  useEffect(() => {
    const code = params.get("code")
    if (!code) {
      router.replace("/signin")
      return
    }

    exchangeCode(code)
      .then(() => {
        router.replace("/")
      })
      .catch(() => {
        router.replace("/signin")
      })
  }, [params, router])

  return <p>登录中...</p>
}

export default function AuthCallback() {
  return (
    <Suspense fallback={<p>登录中...</p>}>
      <AuthCallbackContent />
    </Suspense>
  )
}
