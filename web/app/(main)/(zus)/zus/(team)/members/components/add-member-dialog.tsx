"use client"

import { useState } from "react"
import { UserPlus } from "lucide-react"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { addMember } from "@/service/zus-team"

export function AddMemberDialog({
  teamId,
  onAdded,
}: {
  teamId: string
  onAdded: () => void | Promise<void>
}) {
  const [open, setOpen] = useState(false)
  const [userId, setUserId] = useState("")
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async () => {
    const value = userId.trim()
    if (!value) return

    setSubmitting(true)
    try {
      await addMember(teamId, value)
      toast.success("Member added")
      setUserId("")
      setOpen(false)
      await onAdded()
    } catch {
      toast.error("Failed to add member")
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <Dialog
      open={open}
      onOpenChange={(next) => {
        setOpen(next)
        if (!next) setUserId("")
      }}
    >
      <DialogTrigger asChild>
        <Button size="sm">
          <UserPlus />
          Add Member
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add Member</DialogTitle>
          <DialogDescription>
            Enter the user ID to add to this team.
          </DialogDescription>
        </DialogHeader>
        <FieldGroup>
          <Field>
            <FieldLabel htmlFor="member-user-id">User ID</FieldLabel>
            <Input
              id="member-user-id"
              placeholder="e.g. 3fa85f64-5717-4562-b3fc-2c963f66afa6"
              value={userId}
              onChange={(e) => setUserId(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") void handleSubmit()
              }}
            />
            <FieldDescription>
              The user&apos;s unique identifier (UUID). Ask a team admin if
              you don&apos;t have it.
            </FieldDescription>
          </Field>
        </FieldGroup>
        <DialogFooter>
          <DialogClose asChild>
            <Button variant="outline" disabled={submitting}>
              Cancel
            </Button>
          </DialogClose>
          <Button
            onClick={() => void handleSubmit()}
            disabled={submitting || !userId.trim()}
          >
            Add
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}