# 告警系统设计文档

## 一、告警数据格式

### 1.1 Prometheus 标准告警格式

Prometheus 有两种形态：**触发时的 Alert 结构**和**Alertmanager 接收的 Webhook Payload**。

**Alertmanager Webhook Payload：**

```json
{
  "version": "4",
  "groupKey": "{}:{alertname=\"HighCPU\"}",
  "truncatedAlerts": 0,
  "status": "firing",
  "receiver": "slack-ops",
  "groupLabels": { "alertname": "HighCPU" },
  "commonLabels": {
    "alertname": "HighCPU",
    "env": "production",
    "severity": "critical"
  },
  "commonAnnotations": {
    "summary": "CPU usage above 90%"
  },
  "externalURL": "http://alertmanager:9093",
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "HighCPU",
        "instance": "node-01:9100",
        "severity": "critical"
      },
      "annotations": {
        "summary": "CPU usage above 90% on node-01",
        "description": "CPU is 94.3% for the last 5 minutes",
        "runbook_url": "https://wiki.internal/runbooks/high-cpu"
      },
      "startsAt": "2024-01-15T10:23:00.000Z",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://prometheus:9090/graph?...",
      "fingerprint": "a1b2c3d4e5f6a1b2"
    }
  ]
}
```

### 1.2 关键字段说明

| 字段 | 说明 |
|---|---|
| `labels` | 唯一标识这条告警，参与 fingerprint 计算，用于路由和抑制 |
| `annotations` | 人类可读描述，不参与 fingerprint，用于通知渲染 |
| `startsAt` | 告警首次触发时间，同一次触发多次推送此值不变 |
| `endsAt` | 告警恢复时间，firing 时为零值 `0001-01-01T00:00:00Z` |
| `fingerprint` | labels 集合的哈希值，标识"这类问题"，与时间无关 |

### 1.3 startsAt 和 endsAt 的作用

```
firing 阶段：
  startsAt = 首次触发时间（持续不变）
  endsAt   = 0001-01-01T00:00:00Z（零值，表示未恢复）

resolved 阶段：
  startsAt = 不变
  endsAt   = 实际恢复时间

计算告警持续时长：now - startsAt
计算 MTTR：endsAt - startsAt
```

### 1.4 fingerprint 计算逻辑

fingerprint 只对 **labels 集合**做哈希，和时间完全无关：

```go
func computeFingerprint(labels map[string]string) string {
    // 1. label key 排序，保证顺序一致
    keys := sorted(labels)

    // 2. 拼接 key=value
    for _, k := range keys {
        sb.WriteString(k + "=" + labels[k] + "\n")
    }

    // 3. FNV64a hash
    return fmt.Sprintf("%016x", fnv64a(sb.String()))
}
```

**结论：只要 labels 一模一样，任何时间算出来都是同一个值。**

```
今天:   { alertname="HighCPU", instance="node-01" }  →  fingerprint = 0xa1b2c3d4
两天后: { alertname="HighCPU", instance="node-01" }  →  fingerprint = 0xa1b2c3d4
```

fingerprint 变化的情况：

| 变化 | 结果 |
|---|---|
| `node-01` → `node-02` | fingerprint 变了（不同机器） |
| `severity: warning` → `critical` | fingerprint 变了（标签值变了） |
| 新增 `region="us-east"` | fingerprint 变了（标签集合变了） |

---

## 二、数据库表设计

### 2.1 三表模型

```
alert_definitions  →  这类告警是什么（定义层）
alerts             →  这次触发了（实例层）
alert_events       →  这次触发过程中发生了什么（事件层）
```

### 2.2 alert_definitions

每个 `(fingerprint, connector_id)` 唯一一条，描述告警的静态元数据。

```go
type AlertDefinition struct {
	ID             uuid.UUID      `gorm:"primaryKey;column:id;type:uuid;comment:primary key"`
	Fingerprint    string         `gorm:"column:fingerprint;uniqueIndex:idx_alert_definition_fingerprint;type:varchar(128);comment:alert fingerprint, hash of labels"`
	AlertName      string         `gorm:"column:alert_name;type:varchar(128);not null;comment:alert rule name"`
	Labels         datatypes.JSON `gorm:"column:labels;type:jsonb;not null;comment:full label set, source of fingerprint"`
	Severity       string         `gorm:"column:severity;type:varchar(32);comment:critical / warning / info"`
	ConnectorID    uuid.UUID      `gorm:"column:connector_id;uniqueIndex:idx_alert_definition_fingerprint;type:uuid;comment:FK to connectors"`
	OrganizationID uuid.NullUUID  `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	TeamID         uuid.UUID      `gorm:"column:team_id;type:uuid;default:null;comment:FK to teams"`
	FirstSeenAt    time.Time      `gorm:"column:first_seen_at;not null;comment:first time this fingerprint was seen"`
	LastSeenAt     time.Time      `gorm:"column:last_seen_at;not null;comment:last time this fingerprint was seen, updated on every trigger"`
}

func (AlertDefinition) TableName() string { return "alert_definition" }
```

### 2.3 alerts

每次独立触发一条记录，描述一次完整的告警生命周期。

```go
type Alert struct {
	ID                uuid.UUID      `gorm:"primaryKey;column:id;type:uuid;comment:primary key"`
	AlertDefinitionID uuid.UUID      `gorm:"column:alert_definition_id;type:uuid;not null;index:idx_alerts_alert_id;comment:FK to alert_definition"`
	Fingerprint       string         `gorm:"column:fingerprint;type:varchar(128);not null;uniqueIndex:udx_alerts_fingerprint_starts_at,option:NULLS NOT DISTINCT;index:idx_alerts_fingerprint;comment:FK to alert_definitions"`
	Status            string         `gorm:"column:status;type:varchar(32);not null;comment:firing / resolved"`
	Severity          string         `gorm:"column:severity;type:varchar(32);comment:critical / warning / info"`
	Labels            datatypes.JSON `gorm:"column:labels;type:jsonb;comment:labels"`
	Annotations       datatypes.JSON `gorm:"column:annotations;type:jsonb;comment:human-readable annotations"`
	StartsAt          time.Time      `gorm:"column:starts_at;not null;uniqueIndex:udx_alerts_fingerprint_starts_at,option:NULLS NOT DISTINCT;comment:time when this firing began"`
	EndsAt            *time.Time     `gorm:"column:ends_at;default:null;comment:time when resolved, NULL means still firing"`
	FlapCount         int            `gorm:"column:flap_count;default:0;comment:number of firing<->resolved transitions"`
	OrganizationID    uuid.NullUUID  `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	TeamID            uuid.UUID      `gorm:"column:team_id;type:uuid;default:null;index:idx_alerts_team_id;comment:FK to teams"`
	CreatedAt         time.Time      `gorm:"column:created_at;autoCreateTime;comment:record creation time"`
	UpdatedAt         time.Time      `gorm:"column:updated_at;autoUpdateTime;comment:record last update time"`
}

func (Alert) TableName() string { return "alerts" }
```

### 2.4 alert_events

每次状态变更一条记录，保留原始 payload 用于审计和排查。

```go
type AlertEventType string

const (
	AlertEventFiring     AlertEventType = "firing"
	AlertEventResolved   AlertEventType = "resolved"
	AlertEventFlapping   AlertEventType = "flapping"   // 短时间内反复横跳
	AlertEventSuppressed AlertEventType = "suppressed" // 被 silence / inhibit 压制
	AlertEventNotified   AlertEventType = "notified"   // 已发出通知（用于审计）
	AlertEventEscalated  AlertEventType = "escalated"  // 升级处理
)

type AlertEvent struct {
	ID             uuid.UUID      `gorm:"primaryKey;column:id;type:uuid;comment:primary key"`
	AlertID        uuid.UUID      `gorm:"column:alert_id;type:uuid;not null;index:idx_alert_events_alert_id;comment:FK to alerts"`
	OrganizationID uuid.NullUUID  `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	EventType      AlertEventType `gorm:"column:event_type;type:varchar(32);not null;index:idx_alert_events_event_type;comment:firing / resolved / flapping / suppressed / notified / escalated"`
	OccurredAt     time.Time      `gorm:"column:occurred_at;not null;index:idx_alert_events_occurred_at;comment:time this event occurred"`
	RawPayload     datatypes.JSON `gorm:"column:raw_payload;type:jsonb;comment:original push payload, for audit and incident replay"`
}

func (AlertEvent) TableName() string { return "alert_events" }
```

### 2.5 索引汇总

| 索引名 | 表 | 字段 | 类型 | 作用 |
|---|---|---|---|---|
| `idx_alert_definition_fingerprint` | alert_definition | fingerprint, connector_id | 联合唯一 | 同一 connector 下 fingerprint 唯一 |
| `udx_alerts_fingerprint_starts_at` | alerts | fingerprint, starts_at | 联合唯一 | 防重复入库（NULLS NOT DISTINCT） |
| `idx_alerts_alert_id` | alerts | alert_definition_id | 普通 | 关联查询 alert_definition |
| `idx_alerts_fingerprint` | alerts | fingerprint | 普通 | 查某类告警所有历史 |
| `idx_alerts_organization_id` | alerts | organization_id | 普通 | 多租户隔离 |
| `idx_alerts_team_id` | alerts | team_id | 普通 | 按团队查询 |
| `idx_alert_events_alert_id` | alert_events | alert_id | 普通 | 查某次触发的所有事件 |
| `idx_alert_events_event_type` | alert_events | event_type | 普通 | 按事件类型筛选 |
| `idx_alert_events_occurred_at` | alert_events | occurred_at | 普通 | 时间范围查询 |

`NULLS NOT DISTINCT` 说明：PostgreSQL 默认将 NULL 视为不同值，加此选项后 NULL 也参与唯一约束判断，防止 `ends_at=NULL` 的记录被重复插入。

### 2.6 为什么不加物理外键

告警数据是 append-heavy、需要定期清理的时序型数据，物理外键有两个问题：

1. **清理困难**：必须按顺序删除（先 events → 再 alerts → 再 definitions），大批量时锁表严重
2. **写入开销**：每次写入都要校验外键，高并发下有性能影响

业界做法（Prometheus、Loki、各大 APM 系统）：**不加物理外键，靠应用层保证一致性**。

### 2.7 数据清理策略

分级保留，不同数据保留不同时长：

```
alert_events       →  保留 30 天（量最大，先清）
alerts             →  保留 90 天
alert_definitions  →  永久保留（每类告警一条，量极小）
```

---

## 三、告警标准化

### 3.1 为什么需要标准化

不同来源的告警格式差异大：

| 来源 | status 字段 | severity 字段 | 时间格式 | fingerprint |
|---|---|---|---|---|
| Prometheus | `firing/resolved` | `labels.severity` | RFC3339 | 原生提供 |
| Zabbix | `PROBLEM/RESOLVED/OK` | `Disaster/High/Average/Warning/Info` | Unix timestamp | 无，需自行计算 |

标准化层将所有来源统一为内部格式，后续存储、路由、去重、通知完全不感知来源。

### 3.2 标准化后的统一结构

```go
type NormalizedAlert struct {
    Fingerprint    string            // 统一计算
    Source         string            // prometheus / zabbix / custom
    Status         string            // 统一为 firing / resolved
    Severity       string            // 统一为 critical / warning / info
    Labels         map[string]string // 参与 fingerprint 计算，用于路由
    Annotations    map[string]string // 原始 annotations 全量透传
    StartsAt       time.Time
    EndsAt         *time.Time        // nil = 仍在 firing
    GeneratorURL   string
    OrganizationID uuid.NullUUID
    TeamID         uuid.NullUUID
    RawPayload     []byte            // 原始推送数据，存入 alert_events
}
```

### 3.3 Prometheus 标准化

```
原始字段              →  标准字段
─────────────────────────────────────────
labels.severity       →  Severity（统一大小写映射）
status                →  Status（直接映射）
startsAt              →  StartsAt
endsAt（零值处理）    →  EndsAt（nil 或实际时间）
fingerprint（原生）   →  Fingerprint（优先用，为空则自行计算）
annotations           →  Annotations（直接透传原始 map）
```

**endsAt 零值处理：** Prometheus firing 时 endsAt = `0001-01-01T00:00:00Z`，标准化时转为 `nil`。

**severity 映射：**

```
critical / crit / p1  →  critical
warning / warn / p2   →  warning
info / information / p3 → info
其他                  →  info（兜底）
```

### 3.4 Zabbix 标准化

```
原始字段                    →  标准字段
───────────────────────────────────────────────────
status(PROBLEM)             →  firing
status(RESOLVED/OK)         →  resolved
severity(Disaster/High)     →  critical
severity(Average/Warning)   →  warning
severity(Info/其他)         →  info
clock(Unix timestamp)       →  StartsAt
trigger_name                →  Labels.alertname
host                        →  Labels.instance
event_name / trigger_name   →  Annotations.summary
trigger_description         →  Annotations.description
computeFingerprint(labels)  →  Fingerprint（自行计算）
```

### 3.5 annotations 处理原则

- `summary`：**必须有值**，按优先级降级提取（`summary → message → title → alertname`）
- `description`：允许为空
- `Raw`：原始 annotations 全量保留，防止自定义字段丢失

---

## 四、告警抖动（Flapping）处理

### 4.1 什么是抖动

告警在 firing 和 resolved 之间反复横跳：

```
t0: firing   → 推送 #1
t1: resolved → 推送 #2
t2: firing   → 推送 #3  （抖动）
t3: resolved → 推送 #4  （抖动）
t4: firing   → 推送 #5
```

### 4.2 Prometheus vs Zabbix 的抖动识别差异

| | Prometheus | Zabbix |
|---|---|---|
| 识别依据 | `startsAt` 是否相同 | 时间窗口（FlapWindow） |
| 准确度 | 精确 | 近似 |
| 原理 | 同一次触发 startsAt 不变 | 每次都是新 startsAt，靠时间推断 |

### 4.3 抖动窗口

```go
const DefaultFlapWindow = 5 * time.Minute
```

距上次 resolved 不足 5 分钟再次 firing → 判定为抖动，复活上一条 alert 而非新建。

---

## 五、存储逻辑（Save）

### 5.1 整体流程

```
告警进来
    ↓
① upsertDefinition（alert_definitions UPSERT）
    ↓
② status?
   ├── firing   → ③ handleFiring
   └── resolved → ④ handleResolved
```

### 5.2 第一步：upsertDefinition

```
查 alert_definitions WHERE (fingerprint, connector_id)

├── 不存在 → INSERT 完整记录
│             first_seen_at = StartsAt
│             last_seen_at  = StartsAt
│
└── 存在   → UPDATE last_seen_at = StartsAt
              （其他字段不动）
```

### 5.3 第三步：handleFiring

```
查 alerts WHERE (fingerprint, starts_at, org_id)

├── 找到，status=firing
│     重复推送 → 直接忽略，return nil
│
├── 找到，status=resolved
│     reviveAlert：
│       UPDATE status=firing, ends_at=NULL, flap_count+1
│       INSERT alert_events(flapping)
│
└── 找不到
      isFlapping？
      查最近一条 ends_at IS NOT NULL 的记录
      │
      ├── 距今 < 5分钟（抖动）
      │     handleFlapping：
      │       UPDATE 上一条 status=firing, ends_at=NULL, flap_count+1
      │       INSERT alert_events(flapping)
      │
      └── 距今 >= 5分钟 或 无历史（全新触发）
            createNewAlert：
              INSERT alerts(status=firing, ends_at=NULL, flap_count=0)
              INSERT alert_events(firing)
```

### 5.4 第四步：handleResolved

```
查 alerts WHERE fingerprint AND ends_at IS NULL AND org_id

├── 找到 → UPDATE status=resolved, ends_at=实际结束时间
│           INSERT alert_events(resolved)
│
└── 找不到 → 静默忽略（遗留推送）
```

**ends_at 取值优先级：**

```
1. alert.EndsAt 不为 nil → 用 Prometheus/Zabbix 给的实际时间
2. alert.EndsAt 为 nil   → 用 time.Now() 兜底
   注意：不能用 StartsAt 兜底，StartsAt 是开始时间
```

### 5.5 三张表变化汇总

| 场景 | alert_definitions | alerts | alert_events |
|---|---|---|---|
| 第一次 firing | INSERT | INSERT | INSERT(firing) |
| 重复 firing | UPDATE last_seen_at | 不变 | 不变 |
| resolved | UPDATE last_seen_at | UPDATE ends_at | INSERT(resolved) |
| 抖动（startsAt 相同）| UPDATE last_seen_at | UPDATE flap+1 | INSERT(flapping) |
| 抖动（startsAt 不同）| UPDATE last_seen_at | UPDATE flap+1 | INSERT(flapping) |
| 再次触发（新 startsAt）| UPDATE last_seen_at | INSERT | INSERT(firing) |

---

## 六、常用查询

```sql
-- 当前活跃告警
SELECT * FROM alerts WHERE ends_at IS NULL;

-- 某类告警的所有触发历史
SELECT * FROM alerts
WHERE fingerprint = 'abc123'
ORDER BY starts_at DESC;

-- 某类告警触发了多少次
SELECT COUNT(*) FROM alerts WHERE fingerprint = 'abc123';

-- 平均恢复时长（MTTR）
SELECT AVG(ends_at - starts_at)
FROM alerts
WHERE fingerprint = 'abc123' AND ends_at IS NOT NULL;

-- 最近7天抖动最严重的告警
SELECT fingerprint, SUM(flap_count)
FROM alerts
WHERE starts_at > NOW() - INTERVAL '7 days'
GROUP BY fingerprint
ORDER BY SUM(flap_count) DESC;

-- 某次触发的完整过程
SELECT * FROM alert_events
WHERE alert_id = 'uuid-A'
ORDER BY occurred_at;

-- 某类告警的完整画像
SELECT
    ad.alert_name,
    ad.labels,
    COUNT(a.id)                      AS total_fires,
    MAX(a.starts_at)                 AS last_fired,
    AVG(a.ends_at - a.starts_at)     AS avg_duration
FROM alert_definition ad
JOIN alerts a ON a.fingerprint = ad.fingerprint
WHERE ad.fingerprint = 'abc123'
GROUP BY ad.alert_name, ad.labels;

-- 某个组织的当前活跃告警
SELECT * FROM alerts
WHERE organization_id = 'org-uuid' AND ends_at IS NULL;

-- 某个团队的告警历史
SELECT * FROM alerts
WHERE team_id = 'team-uuid'
ORDER BY starts_at DESC;

-- 某个组织的告警事件流水
SELECT ae.* FROM alert_events ae
WHERE ae.organization_id = 'org-uuid'
ORDER BY ae.occurred_at DESC;
```

---

## 七、扩展新告警来源

只需要实现一个新的 Normalizer，三张表和后续所有逻辑完全不用动：

```go
type Normalizer interface {
    Normalize(raw []byte) ([]*NormalizedAlert, error)
}

// 新增来源只需实现这个接口
type GrafanaNormalizer struct{}
func (n *GrafanaNormalizer) Normalize(raw []byte) ([]*NormalizedAlert, error) { ... }

type CloudWatchNormalizer struct{}
func (n *CloudWatchNormalizer) Normalize(raw []byte) ([]*NormalizedAlert, error) { ... }
```

HTTP 路由层按来源分发：

```
/alerts/prometheus  → PrometheusNormalizer
/alerts/zabbix      → ZabbixNormalizer
/alerts/grafana     → GrafanaNormalizer
```