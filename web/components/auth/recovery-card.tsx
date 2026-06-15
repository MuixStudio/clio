"use client"

import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { z } from "zod"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { InputOTP, InputOTPGroup, InputOTPSlot } from "@/components/ui/input-otp"
import type { RecoveryCardProps } from "./types"
import { Container } from "@/components/auth/container"

type Step = "email" | "code" | "password"

const STEP_TITLE: Record<Step, string> = {
  email: "Forgot password?",
  code: "Check your email",
  password: "Set new password",
}

const emailSchema = z.object({
  email: z.email("Invalid email address"),
})

const codeSchema = z.object({
  code: z.string().length(6, "Please enter the full 6-digit code"),
})

const passwordSchema = z
  .object({
    password: z.string().min(8, "Password must be at least 8 characters"),
    confirm: z.string(),
  })
  .refine((d) => d.password === d.confirm, {
    message: "Passwords do not match",
    path: ["confirm"],
  })

type EmailValues = z.infer<typeof emailSchema>
type CodeValues = z.infer<typeof codeSchema>
type PasswordValues = z.infer<typeof passwordSchema>

export function RecoveryCard({
  method,
  logo,
  title,
  description,
  onNavigateToLogin,
  className,
}: RecoveryCardProps) {
  const [step, setStep] = useState<Step>("email")
  const [email, setEmail] = useState("")

  const stepTitle = title ?? STEP_TITLE[step]
  const stepDesc =
    description ??
    (step === "email"
      ? "Enter your email and we'll send you a verification code"
      : step === "code"
        ? `We sent a code to ${email}`
        : "Choose a strong password for your account")

  const emailForm = useForm<EmailValues>({
    resolver: zodResolver(emailSchema),
    defaultValues: { email: "" },
  })

  const codeForm = useForm<CodeValues>({
    resolver: zodResolver(codeSchema),
    defaultValues: { code: "" },
  })

  const passwordForm = useForm<PasswordValues>({
    resolver: zodResolver(passwordSchema),
    defaultValues: { password: "", confirm: "" },
  })

  async function handleSendCode(values: EmailValues) {
    await method.onSendCode(values.email)
    setEmail(values.email)
    setStep("code")
  }

  async function handleVerifyCode(values: CodeValues) {
    await method.onVerifyCode(email, values.code)
    setStep("password")
  }

  async function handleReset(values: PasswordValues) {
    await method.onReset(values.password)
  }

  function backToEmail() {
    codeForm.reset()
    setStep("email")
  }

  return (
    <Container logo={logo} title={stepTitle} description={stepDesc} className={className}>
      {step === "email" && (
        <form onSubmit={emailForm.handleSubmit(handleSendCode)} className="flex flex-col gap-3">
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="recovery-email">Email</Label>
            <Input
              id="recovery-email"
              type="email"
              placeholder="you@example.com"
              autoComplete="email"
              aria-invalid={!!emailForm.formState.errors.email}
              {...emailForm.register("email")}
            />
            {emailForm.formState.errors.email && (
              <p className="text-xs text-destructive">{emailForm.formState.errors.email.message}</p>
            )}
          </div>
          <Button type="submit" className="w-full" disabled={emailForm.formState.isSubmitting}>
            {emailForm.formState.isSubmitting ? "Sending…" : "Send verification code"}
          </Button>
        </form>
      )}

      {step === "code" && (
        <form onSubmit={codeForm.handleSubmit(handleVerifyCode)} className="flex flex-col gap-3">
          <div className="rounded-lg bg-muted/60 p-3 text-sm text-muted-foreground">
            <p className="font-semibold text-foreground">Reset your password</p>
            <p className="mt-1">
              Confirm your email address to reset your password. We&apos;ve sent a confirmation
              code to{" "}
              <span className="font-semibold text-foreground">{email}</span>.
              {" "}Check your inbox and enter the code here.
            </p>
          </div>

          <div className="flex flex-col items-center gap-1.5">
            <Label>Verification code</Label>
            <InputOTP
              maxLength={6}
              autoComplete="one-time-code"
              autoFocus
              value={codeForm.watch("code")}
              onChange={(val) => codeForm.setValue("code", val, { shouldValidate: true })}
            >
              <InputOTPGroup>
                <InputOTPSlot index={0} />
                <InputOTPSlot index={1} />
                <InputOTPSlot index={2} />
                <InputOTPSlot index={3} />
                <InputOTPSlot index={4} />
                <InputOTPSlot index={5} />
              </InputOTPGroup>
            </InputOTP>
            {codeForm.formState.errors.code && (
              <p className="text-xs text-destructive">{codeForm.formState.errors.code.message}</p>
            )}
          </div>

          <Button type="submit" className="w-full" disabled={codeForm.formState.isSubmitting}>
            {codeForm.formState.isSubmitting ? "Verifying…" : "Verify code"}
          </Button>

          <Button
            type="button"
            variant="link"
            onClick={backToEmail}
            className="inline h-auto w-auto p-0 text-xs font-normal text-muted-foreground underline underline-offset-4 hover:text-foreground"
          >
            Resend code to <span className="font-semibold">{email}</span>
          </Button>
        </form>
      )}

      {step === "password" && (
        <form onSubmit={passwordForm.handleSubmit(handleReset)} className="flex flex-col gap-3">
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="recovery-password">New password</Label>
            <Input
              id="recovery-password"
              type="password"
              placeholder="••••••••"
              autoComplete="new-password"
              autoFocus
              aria-invalid={!!passwordForm.formState.errors.password}
              {...passwordForm.register("password")}
            />
            {passwordForm.formState.errors.password && (
              <p className="text-xs text-destructive">{passwordForm.formState.errors.password.message}</p>
            )}
          </div>

          <div className="flex flex-col gap-1.5">
            <Label htmlFor="recovery-confirm">Confirm new password</Label>
            <Input
              id="recovery-confirm"
              type="password"
              placeholder="••••••••"
              autoComplete="new-password"
              aria-invalid={!!passwordForm.formState.errors.confirm}
              {...passwordForm.register("confirm")}
            />
            {passwordForm.formState.errors.confirm && (
              <p className="text-xs text-destructive">{passwordForm.formState.errors.confirm.message}</p>
            )}
          </div>

          <Button type="submit" className="w-full" disabled={passwordForm.formState.isSubmitting}>
            {passwordForm.formState.isSubmitting ? "Resetting…" : "Reset password"}
          </Button>
        </form>
      )}

      {onNavigateToLogin && (
        <p className="text-center text-sm text-muted-foreground">
          {"Remember your password? "}
          <Button type="button" variant="link" onClick={onNavigateToLogin} className="p-0 font-medium">
            Sign in
          </Button>
        </p>
      )}
    </Container>
  )
}