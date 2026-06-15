"use client"

import { Users } from "lucide-react"

import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Skeleton } from "@/components/ui/skeleton"
import { useZusTeam } from "@/hooks/use-zus-team"

export function TeamSelect({ className }: { className?: string }) {
  const { team, teams, isLoading, selectTeam } = useZusTeam()

  if (isLoading) {
    return <Skeleton className={`h-8 w-40 ${className ?? ""}`} />
  }

  if (teams.length === 0) {
    return null
  }

  return (
    <Select value={team?.id} onValueChange={selectTeam}>
      <SelectTrigger className={className}>
        <Users className="size-4 text-muted-foreground" />
        <SelectValue placeholder="Select a team" />
      </SelectTrigger>
      <SelectContent position="popper">
        <SelectGroup>
          {teams.map((t) => (
            <SelectItem key={t.id} value={t.id}>
              {t.name}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  )
}
