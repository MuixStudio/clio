"use client"

import { Suspense, type ComponentProps, useMemo, useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import {
  ArrowLeft,
  ArrowRight,
  Check,
  Copy,
  Info,
  KeyRound,
  Plus,
  RefreshCw,
  Send,
  ShieldCheck,
  X,
} from "lucide-react"
import { toast } from "sonner"

import { type Provider } from "@/models/connector-provider/types"
import { providers } from "@/models/connector-provider/constants"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { Separator } from "@/components/ui/separator"
import { cn } from "@/lib/utils"
import {
  createConnector,
  enableConnector,
  type Connector,
} from "@/service/zus-connector"
import { useZusTeam } from "@/hooks/use-zus-team"
import {
  HeaderBar,
  HeaderBarLeft,
  HeaderBarTitle,
  HeaderBarDescription,
  HeaderBarActions,
} from "@/app/(main)/components/header-bar/header-bar"

type Step = 1 | 2 | 3 | 4

type LabelRow = {
  key: string
  value: string
}

const STEP_ITEMS: Array<{ value: Step; label: string }> = [
  { value: 1, label: "选择提供商" },
  { value: 2, label: "基础配置" },
  { value: 3, label: "接入凭证 & 文档" },
  { value: 4, label: "验证连接" },
]

const selectableProviders = providers.filter(
  (provider) => provider.id !== "webhook"
)

function getDefaultLabels(provider?: Provider): LabelRow[] {
  const entries = Object.entries(provider?.defaults?.labels ?? {})
  if (entries.length === 0) {
    return [
      { key: "env", value: "prod" },
      { key: "region", value: "cn-east-1" },
    ]
  }
  return entries.map(([key, value]) => ({ key, value }))
}

function WizardStepper({ step }: { step: Step }) {
  return (
    <div className="mx-auto flex w-full max-w-5xl items-center justify-center px-3 py-6">
      {STEP_ITEMS.map((item, index) => {
        const isCompleted = item.value < step
        const isActive = item.value === step

        return (
          <div key={item.value} className="flex min-w-0 items-center">
            <div className="flex items-center gap-2.5">
              <div
                className={cn(
                  "flex size-8 items-center justify-center rounded-full border text-sm font-medium transition-colors",
                  isCompleted &&
                    "border-primary bg-primary text-primary-foreground",
                  isActive &&
                    "border-primary bg-background text-foreground ring-4 ring-muted",
                  !isCompleted &&
                    !isActive &&
                    "border-muted-foreground/30 text-muted-foreground"
                )}
              >
                {isCompleted ? <Check className="size-4" /> : item.value}
              </div>
              <span
                className={cn(
                  "hidden text-sm whitespace-nowrap sm:inline",
                  isActive || isCompleted
                    ? "font-medium text-foreground"
                    : "text-muted-foreground"
                )}
              >
                {item.label}
              </span>
            </div>
            {index < STEP_ITEMS.length - 1 && (
              <div
                className={cn(
                  "mx-3 h-px w-8 bg-border sm:w-14 lg:w-24",
                  isCompleted && "bg-primary"
                )}
              />
            )}
          </div>
        )
      })}
    </div>
  )
}

function ProviderIcon({
  provider,
  className,
}: {
  provider: Provider
  className?: string
}) {
  const Icon = provider.icon
  return <Icon className={cn("size-5", className)} />
}

function ProviderSelectionStep({
  selectedProviderId,
  onSelect,
  onCancel,
  onNext,
}: {
  selectedProviderId: string | null
  onSelect: (provider: Provider) => void
  onCancel: () => void
  onNext: () => void
}) {
  return (
    <Card className="mx-auto w-full max-w-5xl">
      <CardHeader>
        <CardTitle className="text-xl">选择接入源提供商</CardTitle>
        <CardDescription>
          选择要接入的告警源类型，下一步填写实例配置。
        </CardDescription>
      </CardHeader>
      <CardContent className="p-5">
        <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
          {selectableProviders.map((provider) => {
            const selected = provider.id === selectedProviderId
            return (
              <button
                key={provider.id}
                type="button"
                onClick={() => onSelect(provider)}
                className={cn(
                  "flex items-center gap-4 rounded-xl border-2 bg-background p-4 text-left hover:border-primary/40 hover:bg-muted/40",
                  selected &&
                    "border-primary bg-primary/5 hover:border-primary hover:bg-primary/5"
                )}
              >
                <span className="flex size-11 shrink-0 items-center justify-center rounded-lg bg-muted text-foreground">
                  <ProviderIcon provider={provider} />
                </span>
                <span className="min-w-0 flex-1">
                  <span className="block leading-none font-medium">
                    {provider.name}
                  </span>
                  <span className="mt-1.5 block text-sm text-muted-foreground">
                    {provider.subtitle ?? provider.category}
                  </span>
                </span>
                {selected && <Check className="size-4 text-primary" />}
              </button>
            )
          })}
        </div>
      </CardContent>
      <CardFooter className="justify-between">
        <Button type="button" variant="ghost" onClick={onCancel}>
          取消
        </Button>
        <Button type="button" disabled={!selectedProviderId} onClick={onNext}>
          下一步
          <ArrowRight />
        </Button>
      </CardFooter>
    </Card>
  )
}

function ConfigStep({
  provider,
  instanceName,
  labels,
  isCreating,
  onInstanceNameChange,
  onLabelsChange,
  onBack,
  onChangeProvider,
  onCreate,
}: {
  provider: Provider
  instanceName: string
  labels: LabelRow[]
  isCreating: boolean
  onInstanceNameChange: (value: string) => void
  onLabelsChange: (value: LabelRow[]) => void
  onBack: () => void
  onChangeProvider: () => void
  onCreate: ComponentProps<"form">["onSubmit"]
}) {
  const updateLabel = (index: number, field: keyof LabelRow, value: string) => {
    onLabelsChange(
      labels.map((label, currentIndex) =>
        currentIndex === index ? { ...label, [field]: value } : label
      )
    )
  }

  return (
    <Card className="mx-auto w-full max-w-5xl">
      <CardHeader>
        <div className="flex items-start gap-3">
          <div className="flex size-11 shrink-0 items-center justify-center rounded-lg bg-muted text-foreground">
            <ProviderIcon provider={provider} />
          </div>
          <div>
            <CardTitle className="text-xl">接入 {provider.name}</CardTitle>
            <CardDescription>
              推送方式 {provider.pushMethod ?? "Webhook"}
            </CardDescription>
          </div>
        </div>
        <CardAction>
          <Button type="button" variant="outline" onClick={onChangeProvider}>
            <ArrowLeft />
            更换提供商
          </Button>
        </CardAction>
      </CardHeader>
      <form onSubmit={onCreate}>
        <CardContent className="p-5">
          <FieldGroup>
            <Field>
              <FieldLabel htmlFor="connector-name">实例名称 *</FieldLabel>
              <Input
                id="connector-name"
                value={instanceName}
                onChange={(event) => onInstanceNameChange(event.target.value)}
                placeholder={`${provider.name}-生产环境`}
                required
              />
            </Field>

            <Field>
              <div className="flex items-center gap-2">
                <FieldLabel>默认标签</FieldLabel>
                <FieldDescription>（注入每条归一告警）</FieldDescription>
              </div>
              <div className="flex flex-col gap-2">
                {labels.map((label, index) => (
                  <div key={index} className="flex items-center gap-2">
                    <Input
                      value={label.key}
                      className="w-1/4"
                      onChange={(event) =>
                        updateLabel(index, "key", event.target.value)
                      }
                      placeholder="key"
                    />
                    <span className="text-muted-foreground">=</span>
                    <Input
                      value={label.value}
                      className="w-1/4"
                      onChange={(event) =>
                        updateLabel(index, "value", event.target.value)
                      }
                      placeholder="value"
                    />
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      disabled={labels.length === 1}
                      onClick={() =>
                        onLabelsChange(labels.filter((_, i) => i !== index))
                      }
                      aria-label="删除标签"
                    >
                      <X />
                    </Button>
                  </div>
                ))}
              </div>
              <Button
                type="button"
                variant="outline"
                className="w-fit"
                onClick={() =>
                  onLabelsChange([...labels, { key: "", value: "" }])
                }
              >
                <Plus />
                添加标签
              </Button>
            </Field>

            <div className="flex items-start gap-2 rounded-lg border bg-muted/50 p-3 text-sm text-muted-foreground">
              <Info className="mt-0.5 size-4 shrink-0" />
              <span>
                此步骤还没有 Webhook 地址 —— 创建后由平台生成专属地址（含鉴权
                token）。
              </span>
            </div>
          </FieldGroup>
        </CardContent>
        <CardFooter className="justify-between">
          <Button type="button" variant="outline" onClick={onBack}>
            <ArrowLeft />
            上一步
          </Button>
          <Button type="submit" disabled={isCreating || !instanceName.trim() }>
            {isCreating ? "创建中..." : "创建接入源"}
            <ArrowRight />
          </Button>
        </CardFooter>
      </form>
    </Card>
  )
}

function CredentialsStep({
  provider,
  connector,
  webhookUrl,
  onBack,
  onVerify,
}: {
  provider: Provider
  connector: Connector
  webhookUrl: string
  onBack: () => void
  onVerify: () => void
}) {
  const DocsContent = provider.docs?.Content

  return (
    <div className="mx-auto flex w-full max-w-6xl flex-col gap-4">
      <div className="flex items-start gap-3 rounded-xl border border-green-500/20 bg-green-500/10 p-4 text-sm">
        <div className="flex size-8 shrink-0 items-center justify-center rounded-full bg-green-600 text-white">
          <Check className="size-4" />
        </div>
        <div>
          <div className="font-medium">实例「{connector.name}」已创建</div>
          <div className="mt-1 text-muted-foreground">
            下面是它的专属推送地址，配置到 {provider.name} 即可开始接收。
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>① Webhook 推送地址</CardTitle>
            <CardDescription>
              复制后配置到 {provider.name} 的 Webhook 媒介
            </CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col gap-4 p-4">
            <div className="flex items-center gap-2 rounded-lg border bg-background p-2">
              <code className="min-w-0 flex-1 truncate text-xs">
                {webhookUrl}
              </code>
              <Button
                type="button"
                size="sm"
                onClick={() => copyText(webhookUrl, "已复制 Webhook 地址")}
              >
                <Copy />
                复制
              </Button>
            </div>
            <div className="flex flex-wrap gap-2">
              <Button
                type="button"
                variant="outline"
                disabled={!connector.token}
                onClick={() => copyText(connector.token ?? "", "已复制 token")}
              >
                <KeyRound />
                复制 token
              </Button>
              <Button
                type="button"
                variant="outline"
                disabled
                title="暂无 token 轮换 API"
              >
                <RefreshCw />
                轮换
              </Button>
            </div>
            <div className="flex items-start gap-2 rounded-lg bg-muted/50 p-3 text-xs text-muted-foreground">
              <ShieldCheck className="mt-0.5 size-4 shrink-0" />
              <span>对方主动 POST 推送，无需对平台开放入站端口。</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>② 配置文档 · {provider.name}</CardTitle>
            <CardDescription>
              {provider.docs?.title ?? "按步骤完成推送配置"}
            </CardDescription>
            <CardAction>
              <Badge variant="outline">MDX</Badge>
            </CardAction>
          </CardHeader>
          <CardContent className="p-4">
            {DocsContent ? (
              <DocsContent />
            ) : (
              <div className="rounded-lg border border-dashed p-4 text-sm text-muted-foreground">
                该提供商暂未配置文档。
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      <div className="flex items-center justify-between pt-2">
        <Button type="button" variant="outline" onClick={onBack}>
          <ArrowLeft />
          返回修改
        </Button>
        <Button type="button" onClick={onVerify}>
          我已配置，去验证
          <ArrowRight />
        </Button>
      </div>
    </div>
  )
}

function ValidationStep({
  provider,
  connector,
  labels,
  isFinishing,
  onBack,
  onFinish,
}: {
  provider: Provider
  connector: Connector
  labels: LabelRow[]
  isFinishing: boolean
  onBack: () => void
  onFinish: () => void
}) {
  const labelRecord = Object.fromEntries(
    labels
      .filter((label) => label.key.trim())
      .map((label) => [label.key, label.value])
  )
  const env = labelRecord.env ?? "prod"
  const region = labelRecord.region ?? "cn-east-1"

  return (
    <div className="mx-auto flex w-full max-w-6xl flex-col gap-4">
      <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>等待首个事件</CardTitle>
            <CardDescription>触发一条告警以验证链路</CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col gap-4 p-4">
            <div className="space-y-3">
              {[
                "等待来自告警源的首个推送",
                "鉴权 token 将在接收时校验",
                "字段映射成功后会写入归一结果",
              ].map((item) => (
                <div key={item} className="flex items-center gap-2 text-sm">
                  <span className="flex size-5 items-center justify-center rounded-full bg-green-500/15 text-green-600">
                    <Check className="size-3.5" />
                  </span>
                  {item}
                </div>
              ))}
            </div>
            <Separator />
            <div className="flex flex-wrap gap-2">
              <Button
                type="button"
                variant="outline"
                disabled
                title="暂无发送测试事件 API"
              >
                <Send />
                发送测试事件
              </Button>
              <Button
                type="button"
                variant="outline"
                disabled
                title="暂无原始 payload 查询 API"
              >
                查看原始 payload
              </Button>
            </div>
            <div className="flex items-start gap-2 rounded-lg bg-muted/50 p-3 text-xs text-muted-foreground">
              <Info className="mt-0.5 size-4 shrink-0" />
              <span>若长时间无接收，请检查地址 / token / 防火墙出站。</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>归一结果预览</CardTitle>
            <CardDescription>原始报文不入库，仅保留归一后结果</CardDescription>
          </CardHeader>
          <CardContent className="p-4">
            <div className="rounded-xl border bg-muted/30 p-4">
              <div className="flex items-center gap-2">
                <Badge className="bg-orange-500/15 text-orange-600">
                  <span className="size-1.5 rounded-full bg-orange-500" />
                  P1
                </Badge>
                <span className="font-medium">CPU 使用率过高</span>
              </div>
              <div className="mt-4 grid gap-2 text-sm text-muted-foreground">
                <PreviewRow label="host" value="db-prod-03" />
                <PreviewRow label="service" value="storage" />
                <PreviewRow label="env" value={env} />
                <PreviewRow label="region" value={region} />
                <PreviewRow
                  label="source"
                  value={connector.name || provider.name}
                />
                <div className="pt-2 text-foreground">
                  原始级别 High → 映射 P1 ✓
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="flex items-center justify-between pt-2">
        <Button type="button" variant="outline" onClick={onBack}>
          <ArrowLeft />
          上一步
        </Button>
        <Button type="button" disabled={isFinishing} onClick={onFinish}>
          {isFinishing ? "启用中..." : "完成 · 启用接入源"}
        </Button>
      </div>
    </div>
  )
}

function PreviewRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between gap-4 rounded-lg bg-background px-3 py-2">
      <span>{label}</span>
      <span className="font-mono text-xs text-foreground">{value}</span>
    </div>
  )
}

function copyText(value: string, message: string) {
  if (!value) return
  navigator.clipboard
    .writeText(value)
    .then(() => toast.success(message))
    .catch(() => toast.error("复制失败，请手动复制"))
}

function labelsToObject(labels: LabelRow[]) {
  return Object.fromEntries(
    labels
      .map((label) => [label.key.trim(), label.value.trim()])
      .filter(([key]) => key)
  )
}

function ConnectPageContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { team } = useZusTeam()
  const initialType = searchParams.get("type")
  const initialProvider = selectableProviders.find(
    (provider) => provider.id === initialType
  )

  const [step, setStep] = useState<Step>(initialProvider ? 2 : 1)
  const [selectedProviderId, setSelectedProviderId] = useState<string | null>(
    initialProvider?.id ?? null
  )
  const [instanceName, setInstanceName] = useState(
    initialProvider?.defaults?.instanceName ?? ""
  )
  const [labels, setLabels] = useState<LabelRow[]>(() =>
    getDefaultLabels(initialProvider)
  )
  const [connector, setConnector] = useState<Connector | null>(null)
  const [isCreating, setIsCreating] = useState(false)
  const [isFinishing, setIsFinishing] = useState(false)

  const selectedProvider = useMemo(
    () =>
      selectableProviders.find(
        (provider) => provider.id === selectedProviderId
      ),
    [selectedProviderId]
  )

  const webhookUrl = useMemo(() => {
    if (!connector || !selectedProvider) return ""
    const url = `/api/v1/webhook/${selectedProvider.id}/${connector.id}`
    return connector.token
      ? `${url}?token=${encodeURIComponent(connector.token)}`
      : url
  }, [connector, selectedProvider])

  const applyProviderDefaults = (provider: Provider) => {
    setSelectedProviderId(provider.id)
    setInstanceName(provider.defaults?.instanceName ?? "")
    setLabels(getDefaultLabels(provider))
  }

  const handleCreate: ComponentProps<"form">["onSubmit"] = async (event) => {
    event.preventDefault()
    if (!selectedProvider || !instanceName.trim() || !team) return

    setIsCreating(true)
    try {
      const response = await createConnector(
        instanceName.trim(),
        selectedProvider.id,
        labelsToObject(labels),
        team.id
      )
      setConnector(response.data)
      setStep(3)
    } catch (error) {
      console.error(error)
    } finally {
      setIsCreating(false)
    }
  }

  const handleFinish = async () => {
    if (!connector || !selectedProvider || !team) return

    setIsFinishing(true)
    try {
      await enableConnector(team.id, connector.id)
      router.push(`/zus/connector?type=${selectedProvider.id}`)
    } catch (error) {
      console.error(error)
    } finally {
      setIsFinishing(false)
    }
  }

  return (
    <div className="h-full w-full bg-background px-4 py-6">
      <HeaderBar>
        <HeaderBarLeft>
          <HeaderBarTitle>接入提供商</HeaderBarTitle>
          <HeaderBarDescription>
            <p className="text-xs text-muted-foreground">
              Select an integration source to create a new connector.
            </p>
          </HeaderBarDescription>
        </HeaderBarLeft>
        <HeaderBarActions>
          {/* <HeaderBarItem>
            <Button variant="outline">导出</Button>
          </HeaderBarItem>
          <HeaderBarItem>
            <Button>新建</Button>
          </HeaderBarItem> */}
        </HeaderBarActions>
      </HeaderBar>

      <WizardStepper step={step} />

      {step === 1 && (
        <ProviderSelectionStep
          selectedProviderId={selectedProviderId}
          onSelect={applyProviderDefaults}
          onCancel={() => router.push("/zus/connector/provider")}
          onNext={() => selectedProvider && setStep(2)}
        />
      )}

      {step === 2 && selectedProvider && (
        <ConfigStep
          provider={selectedProvider}
          instanceName={instanceName}
          labels={labels}
          isCreating={isCreating}
          onInstanceNameChange={setInstanceName}
          onLabelsChange={setLabels}
          onBack={() => setStep(1)}
          onChangeProvider={() => setStep(1)}
          onCreate={handleCreate}
        />
      )}

      {step === 3 && selectedProvider && connector && (
        <CredentialsStep
          provider={selectedProvider}
          connector={connector}
          webhookUrl={webhookUrl}
          onBack={() => setStep(2)}
          onVerify={() => setStep(4)}
        />
      )}

      {step === 4 && selectedProvider && connector && (
        <ValidationStep
          provider={selectedProvider}
          connector={connector}
          labels={labels}
          isFinishing={isFinishing}
          onBack={() => setStep(3)}
          onFinish={handleFinish}
        />
      )}
    </div>
  )
}

export default function Page() {
  return (
    <Suspense fallback={null}>
      <ConnectPageContent />
    </Suspense>
  )
}
