"use client"

import { useCallback, useEffect, useMemo, useState } from "react"
import {
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
  type SortingState,
} from "@tanstack/react-table"
import { toast } from "sonner"

import {
  HeaderBar,
  HeaderBarActions,
  HeaderBarDescription,
  HeaderBarLeft,
  HeaderBarTitle,
} from "@/app/(main)/components/header-bar/header-bar"
import { DataTable } from "@/components/data-table/data-table"
import { DataTableSkeleton } from "@/components/data-table/data-table-skeleton"
import { DataTableViewOptions } from "@/components/data-table/data-table-view-options"
import { Input } from "@/components/ui/input"
import { useZusTeam } from "@/hooks/use-zus-team"
import { listMembers, removeMember, type TeamMember } from "@/service/zus-team"

import { AddMemberDialog } from "./components/add-member-dialog"
import { getMemberColumns } from "./components/columns"
import { RemoveMemberDialog } from "./components/remove-member-dialog"

export default function Page() {
  const { team, isLoading: isTeamLoading } = useZusTeam()
  const [members, setMembers] = useState<TeamMember[]>([])
  const [loading, setLoading] = useState(true)
  const [sorting, setSorting] = useState<SortingState>([])
  const [globalFilter, setGlobalFilter] = useState("")
  const [removeTarget, setRemoveTarget] = useState<TeamMember | null>(null)
  const [isRemoving, setIsRemoving] = useState(false)

  const fetchMembers = useCallback(async () => {
    if (!team) {
      setMembers([])
      setLoading(false)
      return
    }
    setLoading(true)
    try {
      const { data } = await listMembers(team.id)
      setMembers(data.members)
    } catch {
      toast.error("Failed to load members")
    } finally {
      setLoading(false)
    }
  }, [team])

  useEffect(() => {
    void fetchMembers()
  }, [fetchMembers])

  const handleRemove = async () => {
    if (!team || !removeTarget) return
    setIsRemoving(true)
    try {
      await removeMember(team.id, removeTarget.id)
      toast.success("Member removed")
      setRemoveTarget(null)
      await fetchMembers()
    } catch {
      toast.error("Failed to remove member")
    } finally {
      setIsRemoving(false)
    }
  }

  const columns = useMemo(
    () => getMemberColumns({ onRemove: setRemoveTarget }),
    []
  )

  const table = useReactTable({
    data: members,
    columns,
    state: { sorting, globalFilter },
    onSortingChange: setSorting,
    onGlobalFilterChange: setGlobalFilter,
    getRowId: (row) => row.id,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
  })

  return (
    <>
      {/*header*/}
      <HeaderBar>
        <HeaderBarLeft>
          <HeaderBarTitle className="flex items-center gap-2">
            Members
          </HeaderBarTitle>
          <HeaderBarDescription>
            Manage your team&apos;s members. Add or remove members.
          </HeaderBarDescription>
        </HeaderBarLeft>
        {team && (
          <HeaderBarActions>
            <AddMemberDialog teamId={team.id} onAdded={fetchMembers} />
          </HeaderBarActions>
        )}
      </HeaderBar>

      {/*members*/}
      <div className="flex flex-col gap-4 p-4">
        {isTeamLoading || loading ? (
          <DataTableSkeleton columnCount={4} />
        ) : !team ? (
          <p className="text-sm text-muted-foreground">
            Please select a team first
          </p>
        ) : (
          <DataTable table={table}>
            <div className="flex items-center justify-between gap-2 p-1">
              <Input
                placeholder="Search members..."
                value={globalFilter}
                onChange={(e) => setGlobalFilter(e.target.value)}
                className="h-8 w-64"
              />
              <DataTableViewOptions table={table} />
            </div>
          </DataTable>
        )}
      </div>

      <RemoveMemberDialog
        member={removeTarget}
        open={!!removeTarget}
        onOpenChange={(open) => {
          if (!open) setRemoveTarget(null)
        }}
        isRemoving={isRemoving}
        onConfirmRemove={handleRemove}
      />
    </>
  )
}
