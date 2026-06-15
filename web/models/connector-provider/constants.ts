import { BarChart3, Cloud, CloudCog } from "lucide-react"

import { PrometheusIcon, ZabbixIcon } from "@/components/icons"

import {
  AliyunMonitorDocs,
  AwsCloudWatchDocs,
  GrafanaDocs,
  PrometheusDocs,
  ZabbixDocs,
} from "./docs"
import { type Providers } from "./types"

export const providers: Providers = [
  {
    id: "prometheus",
    name: "Prometheus",
    description:
      "Receive alerts from Prometheus Alertmanager via webhook. Supports firing and resolved events with label-based routing.",
    icon: PrometheusIcon,
    category: "Metrics",
    subtitle: "指标",
    pushMethod: "Webhook",
    tags: ["Alertmanager", "Webhook"],
    defaults: {
      instanceName: "prometheus-生产环境",
      teamId: "infra",
      labels: {
        env: "prod",
        region: "cn-east-1",
      },
    },
    docs: {
      title: "Prometheus Alertmanager",
      Content: PrometheusDocs,
    },
  },
  {
    id: "zabbix",
    name: "Zabbix",
    description:
      "Ingest problem and recovery events from Zabbix via media-type webhook. Maps trigger severity to unified alert levels.",
    icon: ZabbixIcon,
    category: "Infrastructure",
    subtitle: "主机/网络",
    pushMethod: "Webhook",
    tags: ["Webhook", "Trigger"],
    defaults: {
      instanceName: "zabbix-生产环境",
      teamId: "infra",
      labels: {
        env: "prod",
        region: "cn-east-1",
      },
    },
    docs: {
      title: "Zabbix Webhook",
      Content: ZabbixDocs,
    },
  },
  {
    id: "grafana",
    name: "Grafana",
    description:
      "Receive alert notifications from Grafana contact points and normalize labels, severity, and status.",
    icon: BarChart3,
    category: "Observability",
    subtitle: "可视化",
    pushMethod: "Webhook",
    tags: ["Alerting", "Webhook"],
    defaults: {
      instanceName: "grafana-生产环境",
      teamId: "infra",
      labels: {
        env: "prod",
        region: "cn-east-1",
      },
    },
    docs: {
      title: "Grafana Alerting",
      Content: GrafanaDocs,
    },
  },
  {
    id: "aliyun-monitor",
    name: "阿里云监控",
    description:
      "接收阿里云云监控告警回调，归一化云资源、地域、等级与告警状态。",
    icon: CloudCog,
    category: "Cloud",
    subtitle: "云监控",
    pushMethod: "Webhook",
    tags: ["CloudMonitor", "Webhook"],
    defaults: {
      instanceName: "aliyun-生产环境",
      teamId: "infra",
      labels: {
        env: "prod",
        region: "cn-east-1",
      },
    },
    docs: {
      title: "阿里云监控",
      Content: AliyunMonitorDocs,
    },
  },
  {
    id: "aws-cloudwatch",
    name: "AWS CloudWatch",
    description:
      "Ingest CloudWatch alarm notifications through webhook-compatible delivery and normalize cloud dimensions.",
    icon: Cloud,
    category: "Cloud",
    subtitle: "云监控",
    pushMethod: "Webhook",
    tags: ["CloudWatch", "Webhook"],
    defaults: {
      instanceName: "cloudwatch-生产环境",
      teamId: "infra",
      labels: {
        env: "prod",
        region: "us-east-1",
      },
    },
    docs: {
      title: "AWS CloudWatch",
      Content: AwsCloudWatchDocs,
    },
  },
]
