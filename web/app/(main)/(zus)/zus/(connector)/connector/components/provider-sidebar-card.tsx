import Link from "next/link"

import { Badge } from "@/components/ui/badge"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { providers } from "@/models/connector-provider/constants"

export function ProviderSidebarCard() {
  return (
    <Card className="h-fit">
      <CardHeader>
        <CardTitle>Provider</CardTitle>
        <CardDescription>Available alert source types</CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col gap-2 p-1.5">
        {providers.map((provider) => {
          const Icon = provider.icon
          return (
            <Link
              key={provider.id}
              href={`/zus/connector/new?type=${provider.id}`}
              className="group rounded-lg border p-2 hover:border-primary hover:bg-muted/40"
            >
              <div className="flex items-start gap-3">
                <span className="flex size-9 shrink-0 items-center justify-center rounded-md bg-muted text-foreground group-hover:bg-background">
                  <Icon className="size-5" />
                </span>
                <span className="min-w-0 flex-1">
                  <span className="flex items-center gap-2">
                    <span className="truncate text-sm font-medium">{provider.name}</span>
                    <Badge variant="outline" className="shrink-0 text-[10px]">
                      {provider.pushMethod}
                    </Badge>
                  </span>
                  <span className="mt-1 block text-xs text-muted-foreground">
                    {provider.subtitle} · {provider.category}
                  </span>
                  <span className="mt-2 line-clamp-2 block text-xs text-muted-foreground">
                    {provider.description}
                  </span>
                </span>
              </div>
            </Link>
          )
        })}
      </CardContent>
    </Card>
  )
}
