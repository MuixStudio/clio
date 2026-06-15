import { Button } from "@/components/ui/button"
import { useRouter } from "next/navigation"

type Props = {
  accountAgreementUrl?: string
}

export function PrivacyNotice({ accountAgreementUrl }: Props) {
  const router = useRouter()
  const  url = accountAgreementUrl || ""

  return (
    <p className="mt-4 text-center text-xs text-muted-foreground">
      By continuing, you accept the{" "}
      <Button
        type="button"
        variant="link"
        size="sm"
        onClick={() => router.push(url)}
        className="h-auto p-0 text-xs text-muted-foreground underline underline-offset-4 hover:text-foreground"
      >
        Account Agreement
      </Button>
    </p>
  )
}
