"use client"

import { Separator } from "@/components/ui/separator"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { PasswordForm } from "./methods/password-form"
import { EmailCodeForm } from "./methods/email-code-form"
import {
  EmailCodeMethod,
  LoginCardProps,
  LoginMethod,
  OIDCMethod,
  PasswordMethod,
  SAMLMethod,
} from "./types"
import { Button } from "@/components/ui/button"
import { OidcButton } from "@/components/auth/methods/oidc-button"
import { SAMLButton } from "@/components/auth/methods/saml-button"
import { Container } from "@/components/auth/container"
import { Logo } from "@/components/icons"

function isOAuth(m: LoginMethod): m is OIDCMethod {
  return m.type === "oidc"
}

function isSAML(m: LoginMethod): m is SAMLMethod {
  return m.type === "saml"
}

function isForm(m: LoginMethod): m is PasswordMethod | EmailCodeMethod {
  return m.type === "password" || m.type === "email-code"
}

function OIDCSection({ methods }: { methods: OIDCMethod[] }) {
  return (
    <div className="flex flex-col gap-2">
      {methods.map((m) => (
        <OidcButton key={m.provider.id} method={m} />
      ))}
    </div>
  )
}

function SAMLSection({ methods }: { methods: SAMLMethod[] }) {
  return (
    <div className="flex flex-col gap-2">
      {methods.map((m, i) => (
        <SAMLButton key={i} method={m} />
      ))}
    </div>
  )
}

function FormSection({
  methods,
}: {
  methods: (PasswordMethod | EmailCodeMethod)[]
}) {
  if (methods.length === 0) return null

  if (methods.length === 1) {
    const [m] = methods
    if (m.type === "password") return <PasswordForm method={m} mode="login" />
    if (m.type === "email-code") return <EmailCodeForm method={m} />
  }

  return (
    <Tabs defaultValue={methods[0].type}>
      <TabsList className="w-full">
        {methods.map((m) => (
          <TabsTrigger
            key={m.type}
            value={m.type}
            className="flex-1 cursor-pointer"
          >
            {m.name}
          </TabsTrigger>
        ))}
      </TabsList>
      {methods.map((m) => (
        <TabsContent key={m.type} value={m.type} className="mt-4">
          {m.type === "password" && <PasswordForm method={m} mode="login" />}
          {m.type === "email-code" && <EmailCodeForm method={m} />}
        </TabsContent>
      ))}
    </Tabs>
  )
}

export function LoginCard({
  methods,
  logo = <Logo />,
  title = "Welcome back",
  description = "Sign in to your account to continue",
  onNavigateToRegister,
  className,
}: LoginCardProps) {
  const oauthMethods = methods.filter(isOAuth)
  const samlMethods = methods.filter(isSAML)
  const formMethods = methods.filter(isForm)
  const hasButtons = oauthMethods.length > 0 || samlMethods.length > 0
  const showDivider = hasButtons && formMethods.length > 0

  return (
    <Container
      logo={logo}
      title={title}
      description={description}
      className={className}
    >
      <div className="flex flex-col gap-4">
        {oauthMethods.length > 0 && <OIDCSection methods={oauthMethods} />}
        {samlMethods.length > 0 && <SAMLSection methods={samlMethods} />}

        {showDivider && (
          <div className="flex items-center gap-3">
            <Separator className="flex-1" />
            <span className="text-xs text-muted-foreground">or</span>
            <Separator className="flex-1" />
          </div>
        )}

        {formMethods.length > 0 && <FormSection methods={formMethods} />}
      </div>

      {onNavigateToRegister && (
        <p className="mt-6 text-center text-sm text-muted-foreground">
          {"Don't have an account? "}
          <Button
            variant="link"
            onClick={onNavigateToRegister}
            className="p-0 font-medium"
          >
            Sign up
          </Button>
        </p>
      )}
    </Container>
  )
}
