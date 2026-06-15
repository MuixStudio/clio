"use client"

import { useRouter } from "next/navigation"
import { RegisterCard } from "@/components/auth"
import type { RegisterMethod } from "@/components/auth"

const METHODS: RegisterMethod[] = [
  {
    type: "password",
    name: "Password",
  },
  {
    type: "email-code",
    name: "Email",
    onSendCode: async (email) => console.log("Send code to", email),
    onVerify: async (email, code) => console.log("Verify", { email, code }),
  },
]

export default function RegisterPage() {
  const router = useRouter()

  return (
    <RegisterCard
      methods={METHODS}
      onNavigateToLogin={() => router.push("/signin")}
    />
  )
}