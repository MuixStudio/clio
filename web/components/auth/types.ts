import type { ReactNode } from "react"
import type React from "react"

export type PasswordMethod = {
  type: "password"
  name: "Password"
  onForgotPassword?: () => void
}

export type EmailCodeMethod = {
  type: "email-code"
  name: "Email"
  onSendCode: (email: string) => void | Promise<void>
  onVerify: (email: string, code: string) => void | Promise<void>
}

export type SAMLMethod = {
  type: "saml"
  name: "SAML"
  onAction: (organizationSlug: string) => void | Promise<void>
}

export type OIDCMethod = {
  type: "oidc"
  name: "OIDC"
  provider: {
    name: string
    id: "github" | "google"
  }
  icon: React.ComponentType<{ size?: number; className?: string }>
}

export type RecoveryMethod = {
  onSendCode: (email: string) => void | Promise<void>
  onVerifyCode: (email: string, code: string) => void | Promise<void>
  onReset: (password: string) => void | Promise<void>
}

export type RecoveryCardProps = {
  method: RecoveryMethod
  logo?: ReactNode | string
  title?: string
  description?: string
  onNavigateToLogin?: () => void
  accountAgreementUrl?: string
  className?: string
}

export type LoginMethod = PasswordMethod | EmailCodeMethod | SAMLMethod | OIDCMethod
export type RegisterMethod = EmailCodeMethod | PasswordMethod

export type LoginCardProps = {
  methods: LoginMethod[]
  logo?: ReactNode | string
  title?: string
  description?: string
  onNavigateToRegister?: () => void
  accountAgreementUrl?: string
  className?: string
}

export type RegisterCardProps = {
  methods: RegisterMethod[]
  logo?: ReactNode | string
  title?: string
  description?: string
  onNavigateToLogin?: () => void
  accountAgreementUrl?: string
  className?: string
}
