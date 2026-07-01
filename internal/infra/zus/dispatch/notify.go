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

package dispatch

import (
	"context"
	"time"

	"github.com/muixstudio/clio/internal/domain/alert/entity"
	"github.com/prometheus/common/model"
	"go.uber.org/zap"
)

// Notification is one flush of an aggregation group, ready to be delivered to a
// receiver. It carries everything a downstream notifier needs without exposing
// the dispatcher's internal grouping state.
type Notification struct {
	// Receiver is the route's receiver (e.g. a channel name/id). May be empty
	// if the matched route inherited no receiver.
	Receiver string

	// GroupKey uniquely identifies the aggregation group across flushes. It is
	// stable for the lifetime of the group and suitable as a dedup/idempotency
	// key in the notifier.
	GroupKey string

	// GroupLabels are the labels the alerts were grouped by.
	GroupLabels model.LabelSet

	// RepeatInterval is the route's repeat interval, propagated so the notifier
	// can honor it if it implements its own suppression on top of the
	// dispatcher's.
	RepeatInterval time.Duration

	// Alerts is the current batch for this group. Always non-empty.
	Alerts []*entity.NormalizedAlert
}

// Notifier delivers a grouped batch of alerts to its receiver. Implementations
// must be safe for concurrent use: the dispatcher calls Notify from one
// goroutine per aggregation group.
//
// Returning an error keeps the alerts in the group so the next flush retries
// them; returning nil lets the dispatcher drop resolved alerts from the group.
// will be renamed Assigner
type Notifier interface {
	Notify(ctx context.Context, n Notification) error
}

// NotifierFunc adapts a plain function to the Notifier interface.
type NotifierFunc func(ctx context.Context, n Notification) error

func (f NotifierFunc) Notify(ctx context.Context, n Notification) error { return f(ctx, n) }

// LogNotifier is a Notifier that logs each flushed group. It is a safe default
// for wiring up the pipeline before a real channel/courier notifier exists.
type LogNotifier struct{ logger *zap.Logger }

// NewLogNotifier returns a LogNotifier.
func NewLogNotifier(logger *zap.Logger) *LogNotifier {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &LogNotifier{logger: logger.With(zap.String("component", "log-notifier"))}
}

func (n *LogNotifier) Notify(_ context.Context, notification Notification) error {
	n.logger.Info("alert group flushed",
		zap.String("receiver", notification.Receiver),
		zap.String("group_key", notification.GroupKey),
		zap.Int("num_alerts", len(notification.Alerts)))
	return nil
}
