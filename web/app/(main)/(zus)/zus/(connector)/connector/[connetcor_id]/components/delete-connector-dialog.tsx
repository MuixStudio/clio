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

export function DeleteConnectorDialog({
  open,
  onOpenChange,
  connectorName,
  isUpdating,
  onConfirmDelete,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  connectorName: string
  isUpdating: boolean
  onConfirmDelete: () => void | Promise<void>
}) {
  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Delete connector?</AlertDialogTitle>
          <AlertDialogDescription>
            Are you sure you want to delete the connector &quot;{connectorName}&quot;? This
            action cannot be undone. Historical alerts will be preserved.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={isUpdating}>Cancel</AlertDialogCancel>
          <AlertDialogAction
            disabled={isUpdating}
            className={buttonVariants({ variant: "destructive" })}
            onClick={() => void onConfirmDelete()}
          >
            Delete
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
