"use client"

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { PasswordForm } from "./methods/password-form"
import { EmailCodeForm } from "./methods/email-code-form"
import type { RegisterCardProps, RegisterMethod } from "./types"
import { Button } from "@/components/ui/button"
import { Container } from "@/components/auth/container"

const TAB_LABELS: Record<RegisterMethod["type"], string> = {
  "email-code": "Email",
  password: "Password",
}

function FormSection({ methods }: { methods: RegisterMethod[] }) {
  if (methods.length === 0) return null

  if (methods.length === 1) {
    const [m] = methods
    if (m.type === "email-code") return <EmailCodeForm method={m} />
    if (m.type === "password")
      return <PasswordForm mode="register" method={m} />
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
            {TAB_LABELS[m.type]}
          </TabsTrigger>
        ))}
      </TabsList>
      {methods.map((m) => (
        <TabsContent key={m.type} value={m.type} className="mt-4">
          {m.type === "email-code" && <EmailCodeForm method={m} />}
          {m.type === "password" && <PasswordForm mode="register" method={m} />}
        </TabsContent>
      ))}
    </Tabs>
  )
}

export function RegisterCard({
  methods,
  logo,
  title = "Create an account",
  description = "Sign up to get started today",
  onNavigateToLogin,
  className,
}: RegisterCardProps) {
  return (
    <Container
      logo={logo}
      title={title}
      description={description}
      className={className}
    >
      <FormSection methods={methods} />

      {onNavigateToLogin && (
        <p className="text-center text-sm text-muted-foreground">
          {"Already have an account? "}
          <Button
            variant="link"
            onClick={onNavigateToLogin}
            className="p-0 font-medium"
          >
            Sign in
          </Button>
        </p>
      )}
    </Container>
  )
}
