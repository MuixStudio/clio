"use client"

import { useCallback, useEffect, useMemo, useState } from "react"
import Link from "next/link"
import { useParams, useRouter } from "next/navigation"
import { ArrowLeft } from "lucide-react"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { useZusTeam } from "@/hooks/use-zus-team"
import { providers } from "@/models/connector-provider/constants"
import {
  deleteConnector,
  disableConnector,
  enableConnector,
  listConnectors,
} from "@/service/zus-connector"

import { ConnectorBasicInfoCard } from "./components/connector-basic-info-card"
import { ConnectorDangerCard } from "./components/connector-danger-card"
import { ConnectorStatsCard } from "./components/connector-stats-card"
import { ConnectorSummaryCard } from "./components/connector-summary-card"
import { ConnectorWebhookCard } from "./components/connector-webhook-card"
import { DeleteConnectorDialog } from "./components/delete-connector-dialog"
import type { ConnectorItem } from "./components/types"

type ConnectorResponse = { connectors?: ConnectorItem[] } | ConnectorItem[]

function fetchConnectors(teamId: string) {
  return (
    listConnectors as unknown as (
      teamId: string
    ) => Promise<{ data: ConnectorResponse }>
  )(teamId).then(({ data }) => (Array.isArray(data) ? data : (data.connectors ?? [])))
}

function getProvider(type: string) {
  return providers.find((provider) => provider.id === type)
}

function formatDate(value: string) {
  return new Date(value).toLocaleString("en-US", { hour12: false })
}

function copyText(value: string, message: string) {
  if (!value) return
  navigator.clipboard
    .writeText(value)
    .then(() => toast.success(message))
    .catch(() => toast.error("Copy failed, please copy manually"))
}

export default function Page() {
  const params = useParams<{ connetcor_id: string }>()
  const connectorId = params.connetcor_id
  const router = useRouter()
  const { team } = useZusTeam()
  const [connectorState, setConnectorState] = useState<{
    teamId: string
    connectors: ConnectorItem[]
  } | null>(null)
  const [isUpdating, setIsUpdating] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)

  const load = useCallback(() => {
    if (!team) return
    return fetchConnectors(team.id).then((connectors) => {
      setConnectorState({ teamId: team.id, connectors })
    })
  }, [team])

  useEffect(() => {
    let cancelled = false
    if (!team) return

    fetchConnectors(team.id).then((connectors) => {
      if (!cancelled) {
        setConnectorState({ teamId: team.id, connectors })
      }
    })

    return () => {
      cancelled = true
    }
  }, [team])

  const connectors =
    team && connectorState?.teamId === team.id ? connectorState.connectors : []
  const isLoading = Boolean(team && connectorState?.teamId !== team.id)
  const connector = connectors.find((item) => item.id === connectorId)
  const provider = connector ? getProvider(connector.type) : undefined

  const webhookUrl = useMemo(() => {
    if (!connector) return ""
    const url = `/zus/api/v1/webhook/${connector.type}/${connector.id}`
    return connector.token
      ? `${url}?token=${encodeURIComponent(connector.token)}`
      : url
  }, [connector])

  const toggleEnabled = async () => {
    if (!team || !connector) return
    setIsUpdating(true)
    try {
      if (connector.enabled) await disableConnector(team.id, connector.id)
      else await enableConnector(team.id, connector.id)
      toast.success(connector.enabled ? "Connector disabled" : "Connector enabled")
      await load()
    } catch (error) {
      console.error(error)
      toast.error("Operation failed")
    } finally {
      setIsUpdating(false)
    }
  }

  const remove = async () => {
    if (!team || !connector) return

    setIsUpdating(true)
    try {
      await deleteConnector(team.id, connector.id)
      toast.success("Connector deleted")
      setDeleteDialogOpen(false)
      router.push("/zus/connector")
    } catch (error) {
      console.error(error)
      toast.error("Delete failed")
    } finally {
      setIsUpdating(false)
    }
  }

  if (isLoading) {
    return (
      <div className="mx-auto flex w-full max-w-5xl flex-col gap-4 p-4">
        <Skeleton className="h-6 w-28" />
        <Skeleton className="h-24 w-full rounded-xl" />
        <Skeleton className="h-40 w-full rounded-xl" />
      </div>
    )
  }

  if (!connector || !provider) {
    return (
      <div className="mx-auto flex w-full max-w-5xl flex-col items-center gap-3 py-20 text-center text-sm text-muted-foreground">
        Connector not found. It may have been deleted.
        <Button asChild variant="outline">
          <Link href="/app/(main)/zus)">
            <ArrowLeft />
            Back to connectors
          </Link>
        </Button>
      </div>
    )
  }

  return (
    <>
      <div className="mx-auto flex w-full max-w-5xl flex-col gap-4 p-4">
        <Link
          href="/app/(main)/zus)"
          className="flex w-fit items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground"
        >
          <ArrowLeft className="size-4" />
          Connectors
        </Link>

        <ConnectorSummaryCard
          connector={connector}
          provider={provider}
          createdAtLabel={formatDate(connector.created_at)}
          isUpdating={isUpdating}
          onToggleEnabled={toggleEnabled}
          onRequestDelete={() => setDeleteDialogOpen(true)}
        />

        <ConnectorBasicInfoCard
          connector={connector}
          provider={provider}
          teamName={team?.name}
        />

        <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
          <ConnectorWebhookCard
            connector={connector}
            provider={provider}
            webhookUrl={webhookUrl}
            onCopyText={copyText}
          />
          <ConnectorStatsCard />
        </div>

        <ConnectorDangerCard
          enabled={connector.enabled}
          isUpdating={isUpdating}
          onToggleEnabled={toggleEnabled}
          onRequestDelete={() => setDeleteDialogOpen(true)}
        />
      </div>

      <DeleteConnectorDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        connectorName={connector.name}
        isUpdating={isUpdating}
        onConfirmDelete={remove}
      />
    </>
  )
}
