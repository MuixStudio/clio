"use client"

import { useRouter } from "next/navigation"
import { LoginCard } from "@/components/auth"

import { GithubIcon, GoogleIcon } from "@/components/icons"
import { LoginMethod } from "@/components/auth/types"

export default function LoginPage() {
  const router = useRouter()

  const methods: LoginMethod[] = [
    {
      type: "oidc",
      name: "OIDC",
      icon: GoogleIcon,
      provider: {
        name: "Google",
        id: "google",
      },
    },
    {
      type: "oidc",
      name: "OIDC",
      icon: GithubIcon,
      provider: {
        name: "Github",
        id: "github",
      },
    },
    {
      type: "password",
      name: "Password",
      onForgotPassword: () => router.push("/recovery"),
    },
    {
      type: "email-code",
      name: "Email",
      onSendCode: async (email) => console.log("Send code to", email),
      onVerify: async (email, code) => console.log("Verify", { email, code }),
    },
    {
      type: "saml",
      name: "SAML",
      onAction: async (slug) => console.log("SSO with org", slug),
    },
  ]

  return (
    <LoginCard
      methods={methods}
      onNavigateToRegister={() => router.push("/signup")}
    />
  )
}
