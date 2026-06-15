"use client"

import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { z } from "zod"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import type { OIDCMethod } from "../types"

const schema = z.object({
  slug: z.string().min(1, "Organization is required").regex(
    /^[a-z0-9-]+$/,
    "Only lowercase letters, numbers and hyphens"
  ),
})

type Values = z.infer<typeof schema>

export function SSOForm({ method }: { method: OIDCMethod }) {
  const { register, handleSubmit, formState: { isSubmitting, errors } } = useForm<Values>({
    resolver: zodResolver(schema),
    defaultValues: { slug: "" },
  })

  async function onSubmit(values: Values) {
    // await method.onAction(values.slug)
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-3">
      <div className="flex flex-col gap-1.5">
        <Label htmlFor="auth-sso">Organization</Label>
        <Input
          id="auth-sso"
          placeholder="your-organization"
          autoComplete="organization"
          autoCapitalize="none"
          autoCorrect="off"
          spellCheck={false}
          aria-invalid={!!errors.slug}
          {...register("slug")}
        />
        {errors.slug
          ? <p className="text-xs text-destructive">{errors.slug.message}</p>
          : <p className="text-xs text-muted-foreground">Enter your organization slug to continue with SSO</p>
        }
      </div>

      <Button type="submit" className="w-full" disabled={isSubmitting}>
        {isSubmitting ? "Redirecting…" : "Continue with SSO"}
      </Button>
    </form>
  )
}
