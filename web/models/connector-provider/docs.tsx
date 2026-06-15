function DocSection({
  title,
  description,
  code,
}: {
  title: string
  description: string
  code?: string
}) {
  return (
    <section className="rounded-lg border p-3">
      <h4 className="font-medium">{title}</h4>
      <p className="mt-1 text-sm text-muted-foreground">{description}</p>
      {code && (
        <pre className="mt-3 overflow-x-auto rounded-lg bg-foreground p-3 text-xs text-background">
          <code>{code}</code>
        </pre>
      )}
    </section>
  )
}

export function PrometheusDocs() {
  return (
    <div className="space-y-3">
      <DocSection
        title="打开 Alertmanager 配置"
        description="在 alertmanager.yml 中找到 receivers 配置段。"
      />
      <DocSection
        title="粘贴 Webhook URL"
        description="新增 webhook_configs，并将平台生成的 URL 填入 url 字段。"
        code={`receivers:\n  - name: clio\n    webhook_configs:\n      - url: <WEBHOOK_URL>`}
      />
      <DocSection
        title="绑定路由并重载"
        description="将 receiver 绑定到 route 后，重载 Alertmanager 配置即可生效。"
      />
    </div>
  )
}

export function ZabbixDocs() {
  return (
    <div className="space-y-3">
      <DocSection
        title="新建 Webhook 媒介类型"
        description="管理 → 报警媒介类型 → 新建，类型选择 Webhook。"
      />
      <DocSection
        title="粘贴 URL 与消息模板"
        description="将平台生成的 URL 填入脚本参数，并按需要配置告警消息模板。"
        code={`URL=<WEBHOOK_URL>\nTOKEN=<TOKEN>`}
      />
      <DocSection
        title="绑定到动作 / 接收人"
        description="在「动作」里把该媒介指给接收人即可生效。"
      />
    </div>
  )
}

export function GrafanaDocs() {
  return (
    <div className="space-y-3">
      <DocSection
        title="创建 Contact point"
        description="Alerting → Contact points → New contact point，类型选择 Webhook。"
      />
      <DocSection
        title="配置 Webhook URL"
        description="将平台生成的推送地址粘贴到 URL，并保存 Contact point。"
      />
      <DocSection
        title="绑定 Notification policy"
        description="在 Notification policies 中将告警路由到该 Contact point。"
      />
    </div>
  )
}

export function AliyunMonitorDocs() {
  return (
    <div className="space-y-3">
      <DocSection
        title="新建 Webhook 通知渠道"
        description="进入云监控告警联系人或通知渠道配置，新增 Webhook 类型渠道。"
      />
      <DocSection
        title="粘贴 URL 与消息模板"
        description="将平台生成的 Webhook 地址配置到回调 URL，并使用默认告警模板。"
        code={`POST <WEBHOOK_URL>\nContent-Type: application/json`}
      />
      <DocSection
        title="绑定告警规则"
        description="在告警规则中选择该通知渠道，保存后触发一条告警进行验证。"
      />
    </div>
  )
}

export function AwsCloudWatchDocs() {
  return (
    <div className="space-y-3">
      <DocSection
        title="准备告警通知目标"
        description="在 CloudWatch Alarm 的通知目标中接入可转发 Webhook 的通道。"
      />
      <DocSection
        title="配置推送地址"
        description="将平台生成的 Webhook 地址配置到通知通道的目标 URL。"
      />
      <DocSection
        title="触发测试告警"
        description="使用 CloudWatch Alarm 测试事件验证推送、鉴权与字段映射。"
      />
    </div>
  )
}
