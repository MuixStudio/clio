import Link from "next/link"
import { BookOpen, Copy, Info, KeyRound, RefreshCw } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

import type { ConnectorItem, ConnectorProvider } from "./types"

export function ConnectorWebhookCard({
  connector,
  provider,
  webhookUrl,
  onCopyText,
}: {
  connector: ConnectorItem
  provider: ConnectorProvider
  webhookUrl: string
  onCopyText: (value: string, message: string) => void
}) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Webhook Credentials</CardTitle>
        <CardDescription>
          Copy and configure this as the {provider.name} webhook destination
        </CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col gap-3 p-4">
        <div className="flex items-center gap-2 rounded-lg border bg-muted/40 p-2">
          <code className="min-w-0 flex-1 truncate text-xs">{webhookUrl}</code>
          <Button
            type="button"
            size="sm"
            onClick={() => onCopyText(webhookUrl, "Webhook URL copied")}
          >
            <Copy />
            Copy
          </Button>
        </div>
        <div className="flex flex-wrap gap-2">
          <Button
            type="button"
            variant="outline"
            disabled={!connector.token}
            onClick={() => onCopyText(connector.token ?? "", "Token copied")}
          >
            <KeyRound />
            Copy Token
          </Button>
          <Button type="button" variant="outline" disabled title="Token rotation API is not available yet">
            <RefreshCw />
            Rotate Token
          </Button>
          {provider.docs && (
            <Button asChild variant="outline">
              <Link href={`/app/(main)/(zus)/zus/(connector)/connector/new?type=${provider.id}`}>
                <BookOpen />
                Documentation
              </Link>
            </Button>
          )}
        </div>
        <div className="flex items-start gap-2 rounded-lg bg-muted/50 p-3 text-xs text-muted-foreground">
          <Info className="mt-0.5 size-4 shrink-0" />
          <span>
            The provider pushes via outbound POST requests, so no inbound port needs to be
            opened on this platform. If nothing has been received for a while, check the URL,
            token, and outbound firewall rules.
          </span>
        </div>
      </CardContent>
    </Card>
  )
}
