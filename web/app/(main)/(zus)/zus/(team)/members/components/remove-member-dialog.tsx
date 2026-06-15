import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { buttonVariants } from "@/components/ui/button"
import type { TeamMember } from "@/service/zus-team"

export function RemoveMemberDialog({
  member,
  open,
  onOpenChange,
  isRemoving,
  onConfirmRemove,
}: {
  member: TeamMember | null
  open: boolean
  onOpenChange: (open: boolean) => void
  isRemoving: boolean
  onConfirmRemove: () => void | Promise<void>
}) {
  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Remove member?</AlertDialogTitle>
          <AlertDialogDescription>
            Are you sure you want to remove &quot;{member?.email}&quot; from
            this team? This action cannot be undone.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={isRemoving}>Cancel</AlertDialogCancel>
          <AlertDialogAction
            disabled={isRemoving}
            className={buttonVariants({ variant: "destructive" })}
            onClick={() => void onConfirmRemove()}
          >
            Remove
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}