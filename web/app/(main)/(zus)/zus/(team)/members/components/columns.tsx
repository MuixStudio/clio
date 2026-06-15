"use client"

import type { ColumnDef } from "@tanstack/react-table"
import { format } from "date-fns"
import { Copy, MoreHorizontal, UserMinus } from "lucide-react"
import { toast } from "sonner"

import { DataTableColumnHeader } from "@/components/data-table/data-table-column-header"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import type { TeamMember } from "@/service/zus-team"

export function getMemberColumns({
  onRemove,
}: {
  onRemove: (member: TeamMember) => void
}): ColumnDef<TeamMember>[] {
  return [
    {
      accessorKey: "email",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} label="Member" />
      ),
      cell: ({ row }) => {
        const email = row.original.email
        return (
          <div className="flex items-center gap-2">
            <Avatar size="sm">
              <AvatarFallback>{email.slice(0, 1).toUpperCase()}</AvatarFallback>
            </Avatar>
            <span className="font-medium">{email}</span>
          </div>
        )
      },
      meta: { label: "Member" },
    },
    {
      accessorKey: "user_id",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} label="User ID" />
      ),
      cell: ({ row }) => {
        const userId = row.original.user_id
        return (
          <div className="flex items-center gap-1">
            <span className="font-mono text-xs text-muted-foreground">
              {userId}
            </span>
            <Button
              variant="ghost"
              size="icon"
              className="size-6"
              onClick={() => {
                void navigator.clipboard.writeText(userId)
                toast.success("User ID copied")
              }}
            >
              <Copy />
            </Button>
          </div>
        )
      },
      meta: { label: "User ID" },
    },
    {
      accessorKey: "joined_at",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} label="Joined At" />
      ),
      cell: ({ row }) =>
        format(new Date(row.original.joined_at), "yyyy-MM-dd HH:mm"),
      meta: { label: "Joined At" },
    },
    {
      id: "actions",
      cell: ({ row }) => (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon" className="size-8">
              <MoreHorizontal />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem
              variant="destructive"
              onClick={() => onRemove(row.original)}
            >
              <UserMinus />
              Remove member
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      ),
      enableSorting: false,
      enableHiding: false,
      size: 48,
    },
  ]
}