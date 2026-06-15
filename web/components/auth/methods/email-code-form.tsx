"use client"

import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { z } from "zod"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { InputOTP, InputOTPGroup, InputOTPSlot } from "@/components/ui/input-otp"
import type { EmailCodeMethod } from "../types"

const emailSchema = z.object({
  email: z.email("Invalid email address"),
})

const codeSchema = z.object({
  code: z.string().length(6, "Please enter the full 6-digit code"),
})

type EmailValues = z.infer<typeof emailSchema>
type CodeValues = z.infer<typeof codeSchema>

type Step = "email" | "code"

export function EmailCodeForm({ method }: { method: EmailCodeMethod }) {
  const [step, setStep] = useState<Step>("email")
  const [email, setEmail] = useState("")

  const emailForm = useForm<EmailValues>({
    resolver: zodResolver(emailSchema),
    defaultValues: { email: "" },
  })

  const codeForm = useForm<CodeValues>({
    resolver: zodResolver(codeSchema),
    defaultValues: { code: "" },
  })

  async function handleSendCode(values: EmailValues) {
    await method.onSendCode(values.email)
    setEmail(values.email)
    setStep("code")
  }

  async function handleVerify(values: CodeValues) {
    await method.onVerify(email, values.code)
  }

  if (step === "code") {
    const { formState: { isSubmitting, errors } } = codeForm

    return (
      <form onSubmit={codeForm.handleSubmit(handleVerify)} className="flex flex-col gap-3">
        <p className="text-sm text-muted-foreground">
          We sent a code to <span className="font-medium text-foreground">{email}</span>
        </p>

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
          {errors.code && (
            <p className="text-xs text-destructive">{errors.code.message}</p>
          )}
        </div>

        <Button type="submit" className="w-full" disabled={isSubmitting}>
          {isSubmitting ? "Verifying…" : "Verify code"}
        </Button>

        <Button
          type="button"
          variant="link"
          onClick={() => { setStep("email"); codeForm.reset() }}
          className="inline h-auto p-0 text-xs font-normal text-muted-foreground underline underline-offset-4 hover:text-foreground"
        >
          Use a different email <span className="font-semibold">{email}</span>
        </Button>
      </form>
    )
  }

  const { formState: { isSubmitting, errors } } = emailForm

  return (
    <form onSubmit={emailForm.handleSubmit(handleSendCode)} className="flex flex-col gap-3">
      <div className="flex flex-col gap-1.5">
        <Label htmlFor="auth-email-code">Email</Label>
        <Input
          id="auth-email-code"
          type="email"
          placeholder="you@example.com"
          autoComplete="email"
          aria-invalid={!!errors.email}
          {...emailForm.register("email")}
        />
        {errors.email && (
          <p className="text-xs text-destructive">{errors.email.message}</p>
        )}
      </div>

      <Button type="submit" className="w-full" disabled={isSubmitting}>
        {isSubmitting ? "Sending…" : "Send code"}
      </Button>
    </form>
  )
}
