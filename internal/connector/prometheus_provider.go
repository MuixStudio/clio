/*
 * Copyright 2026 MuixStudio
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package connector

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/alert"
)

// ============================================================
// Prometheus 原始数据结构
// ============================================================

type PrometheusWebhookPayload struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	TruncatedAlerts   int               `json:"truncatedAlerts"`
	Status            string            `json:"status"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []PrometheusAlert `json:"alerts"`
}

type PrometheusAlert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

// ============================================================
// PrometheusNormalizer
// ============================================================

type PrometheusNormalizer struct{}

func NewPrometheusNormalizer() *PrometheusNormalizer {
	return &PrometheusNormalizer{}
}

func (n *PrometheusNormalizer) Normalize(raw []byte) ([]*alert.NormalizedAlert, error) {
	var payload PrometheusWebhookPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal prometheus payload: %w", err)
	}

	if len(payload.Alerts) == 0 {
		return nil, fmt.Errorf("prometheus payload contains no alerts")
	}

	results := make([]*alert.NormalizedAlert, 0, len(payload.Alerts))
	for _, alert := range payload.Alerts {
		normalized, err := n.normalizeAlert(alert)
		if err != nil {
			return nil, fmt.Errorf("failed to normalize alert %s: %w", alert.Fingerprint, err)
		}
		results = append(results, normalized)
	}

	return results, nil
}

func (n *PrometheusNormalizer) Type() string {
	return "prometheus"
}

func (n *PrometheusNormalizer) normalizeAlert(pa PrometheusAlert) (*alert.NormalizedAlert, error) {
	status, err := mapPrometheusStatus(pa.Status)
	if err != nil {
		return nil, err
	}

	fingerprint := pa.Fingerprint
	if fingerprint == "" {
		fingerprint = computeFingerprint(pa.Labels)
	}

	alertRaw, err := json.Marshal(pa)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal alert raw: %w", err)
	}

	// annotations 直接透传原始 map
	// 在需要 summary / description 的地方（通知渲染等）再按优先级提取
	annotations := make(map[string]string, len(pa.Annotations))
	for k, v := range pa.Annotations {
		annotations[k] = v
	}

	return &alert.NormalizedAlert{
		ID:          uuid.New(),
		Fingerprint: fingerprint,
		//ConnectorId:    uuid.NullUUID{},
		Status:      status,
		Severity:    extractSeverity(pa.Labels),
		Labels:      pa.Labels,
		Annotations: annotations,
		StartsAt:    pa.StartsAt,
		EndsAt:      extractEndsAt(status, pa.EndsAt),
		//TeamID:         n.TeamID,
		RawPayload: alertRaw,
	}, nil
}

func mapPrometheusStatus(status string) (alert.Status, error) {
	switch status {
	case "firing":
		return alert.StatusFiring, nil
	case "resolved":
		return alert.StatusResolved, nil
	default:
		return "", fmt.Errorf("unknown prometheus status: %s", status)
	}
}

func extractEndsAt(status alert.Status, endsAt time.Time) *time.Time {
	if status != "resolved" {
		return nil
	}
	if endsAt.IsZero() || endsAt.Year() == 1 {
		return nil
	}
	return &endsAt
}

func extractSeverity(labels map[string]string) alert.Severity {
	raw, ok := labels["severity"]
	if !ok || raw == "" {
		return "info"
	}
	switch strings.ToLower(raw) {
	case "critical", "crit", "p1":
		return alert.SeverityCritical
	case "warning", "warn", "p2":
		return alert.SeverityHigh
	case "info", "information", "p3":
		return alert.SeverityMedium
	default:
		return alert.SeverityLow
	}
}

func computeFingerprint(labels map[string]string) string {
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(labels[k])
		sb.WriteByte('\n')
	}

	h := fnv.New64a()
	h.Write([]byte(sb.String()))
	return fmt.Sprintf("%016x", h.Sum64())
}
