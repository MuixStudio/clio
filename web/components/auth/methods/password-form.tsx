"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { z } from "zod"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import type { PasswordMethod } from "../types"
import { Login } from "@/service/auth"

// ── Schemas ──────────────────────────────────────────────────────────────────

const loginSchema = z.object({
  identifier: z
    .string()
    .min(4, "At least 4 characters")
    .max(20, "At most 20 characters")
    .regex(/^[a-zA-Z0-9.@]+$/, "Letters, numbers, . and @ only"),
  password: z.string().min(1, "Password is required"),
})

const registerSchema = z
  .object({
    identifier: z
      .string()
      .min(4, "At least 4 characters")
      .max(20, "At most 20 characters")
      .regex(/^[a-zA-Z0-9.@]+$/, "Letters, numbers, . and @ only"),
    password: z.string().min(8, "Password must be at least 8 characters"),
    confirm: z.string(),
  })
  .refine((d) => d.password === d.confirm, {
    message: "Passwords do not match",
    path: ["confirm"],
  })

type LoginValues = z.infer<typeof loginSchema>
type RegisterValues = z.infer<typeof registerSchema>

// ── Login form ────────────────────────────────────────────────────────────────

function LoginFormInner({ method }: { method: PasswordMethod }) {
  const router = useRouter()
  const [serverError, setServerError] = useState<string | null>(null)
  const {
    register,
    handleSubmit,
    formState: { isSubmitting, errors },
  } = useForm<LoginValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: { identifier: "", password: "" },
  })

  async function onSubmit(v: LoginValues) {
    setServerError(null)
    try {
      await Login({
        method: "password",
        data: { identifier: v.identifier, password: v.password },
      })
      router.push("/")
    } catch {
      setServerError("Invalid identifier or password")
    }
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-3">
      <div className="flex flex-col gap-1.5">
        <Label htmlFor="auth-identifier">Username</Label>
        <Input
          id="auth-identifier"
          type="text"
          placeholder="username"
          autoComplete="username"
          aria-invalid={!!errors.identifier}
          {...register("identifier")}
        />
        {errors.identifier && (
          <p className="text-xs text-destructive">
            {errors.identifier.message}
          </p>
        )}
      </div>

      <div className="flex flex-col gap-1.5">
        <div className="flex items-center justify-between">
          <Label htmlFor="auth-password">Password</Label>
          {method.onForgotPassword && (
            <Button
              type="button"
              variant="link"
              size="sm"
              onClick={method.onForgotPassword}
              className="p-0 text-xs text-muted-foreground"
            >
              Forgot password?
            </Button>
          )}
        </div>
        <Input
          id="auth-password"
          type="password"
          placeholder="••••••••"
          autoComplete="current-password"
          aria-invalid={!!errors.password}
          {...register("password")}
        />
        {errors.password && (
          <p className="text-xs text-destructive">{errors.password.message}</p>
        )}
      </div>

      {serverError && <p className="text-xs text-destructive">{serverError}</p>}

      <Button type="submit" className="w-full" disabled={isSubmitting}>
        {isSubmitting ? "Signing in…" : "Sign in"}
      </Button>
    </form>
  )
}

// ── Register form ─────────────────────────────────────────────────────────────

function RegisterFormInner(_: { method: PasswordMethod }) {
  const {
    register,
    handleSubmit,
    formState: { isSubmitting, errors },
  } = useForm<RegisterValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: { identifier: "", password: "", confirm: "" },
  })

  return (
    <form onSubmit={handleSubmit(() => {})} className="flex flex-col gap-3">
      <div className="flex flex-col gap-1.5">
        <Label htmlFor="auth-identifier">Username</Label>
        <Input
          id="auth-identifier"
          type="text"
          placeholder="username"
          autoComplete="username"
          aria-invalid={!!errors.identifier}
          {...register("identifier")}
        />
        {errors.identifier && (
          <p className="text-xs text-destructive">
            {errors.identifier.message}
          </p>
        )}
      </div>

      <div className="flex flex-col gap-1.5">
        <Label htmlFor="auth-password">Password</Label>
        <Input
          id="auth-password"
          type="password"
          placeholder="••••••••"
          autoComplete="new-password"
          aria-invalid={!!errors.password}
          {...register("password")}
        />
        {errors.password && (
          <p className="text-xs text-destructive">{errors.password.message}</p>
        )}
      </div>

      <div className="flex flex-col gap-1.5">
        <Label htmlFor="auth-confirm">Confirm password</Label>
        <Input
          id="auth-confirm"
          type="password"
          placeholder="••••••••"
          autoComplete="new-password"
          aria-invalid={!!errors.confirm}
          {...register("confirm")}
        />
        {errors.confirm && (
          <p className="text-xs text-destructive">{errors.confirm.message}</p>
        )}
      </div>

      <Button type="submit" className="w-full" disabled={isSubmitting}>
        {isSubmitting ? "Creating account…" : "Create account"}
      </Button>
    </form>
  )
}

// ── Public export ─────────────────────────────────────────────────────────────

export function PasswordForm({
  method,
  mode,
}: {
  method: PasswordMethod
  mode: "login" | "register"
}) {
  if (mode === "register") return <RegisterFormInner method={method} />
  return <LoginFormInner method={method} />
}
