"use client"

import { useRouter } from "next/navigation"
import { RecoveryCard } from "@/components/auth/recovery-card"

const LOGO = (
  <div className="flex size-10 items-center justify-center rounded-xl bg-primary text-primary-foreground font-bold text-lg">
    C
  </div>
)

export default function RecoveryPage() {
  const router = useRouter()

  return (
    <RecoveryCard
      method={{
        onSendCode: async (email) => console.log("Send code to", email),
        onVerifyCode: async (email, code) =>
          console.log("Verify code", { email, code }),
        onReset: async (password) =>
          console.log("Reset password", { password }),
      }}
      logo={LOGO}
      onNavigateToLogin={() => router.push("/signin")}
    />
  )
}
