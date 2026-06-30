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
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/muixstudio/clio/internal/domain/alert/entity"
	"github.com/prometheus/common/model"
	"go.uber.org/zap"
)

// flushFunc delivers a batch of alerts. It returns true if delivery succeeded,
// in which case resolved alerts may be evicted from the group.
type flushFunc func(ctx context.Context, alerts []*entity.NormalizedAlert) bool

// aggrGroup collects alerts that share a set of group labels under a single
// route and flushes them on a schedule: an initial GroupWait, then every
// GroupInterval. Within those flushes it suppresses repeats — a group is only
// delivered when its alert set changed or RepeatInterval has elapsed since the
// last successful delivery. This mirrors Alertmanager's dispatch + dedup
// behaviour closely enough to be safe in production without a full nflog.
type aggrGroup struct {
	labels   model.LabelSet
	opts     *RouteOpts
	groupKey string
	logger   *zap.Logger
	timeout  func(time.Duration) time.Duration

	ctx    context.Context
	cancel func()
	done   chan struct{}
	next   *time.Timer

	mtx        sync.Mutex
	alerts     map[string]*entity.NormalizedAlert // keyed by alert fingerprint
	hasFlushed bool
	lastFlush  time.Time
	lastDigest string
}

func newAggrGroup(ctx context.Context, labels model.LabelSet, r *Route, to func(time.Duration) time.Duration, logger *zap.Logger) *aggrGroup {
	if to == nil {
		to = func(d time.Duration) time.Duration { return d }
	}
	ag := &aggrGroup{
		labels:   labels,
		opts:     &r.RouteOpts,
		groupKey: fmt.Sprintf("%s:%s", r.Key(), labels),
		timeout:  to,
		alerts:   map[string]*entity.NormalizedAlert{},
		done:     make(chan struct{}),
	}
	ag.ctx, ag.cancel = context.WithCancel(ctx)
	ag.logger = logger.With(zap.String("aggr_group", ag.groupKey))
	// Wait GroupWait before the first flush.
	ag.next = time.NewTimer(ag.opts.GroupWait)
	return ag
}

// run drives the group's flush schedule until its context is cancelled.
func (ag *aggrGroup) run(nf flushFunc) {
	defer close(ag.done)
	defer ag.next.Stop()

	for {
		select {
		case <-ag.next.C:
			// Bound each flush so a slow notifier can't block the next one.
			ctx, cancel := context.WithTimeout(ag.ctx, ag.timeout(ag.opts.GroupInterval))

			ag.mtx.Lock()
			ag.next.Reset(ag.opts.GroupInterval)
			ag.hasFlushed = true
			ag.mtx.Unlock()

			ag.flush(ctx, nf)
			cancel()

		case <-ag.ctx.Done():
			return
		}
	}
}

func (ag *aggrGroup) stop() {
	ag.cancel()
	<-ag.done
}

// insert adds or replaces an alert in the group. If the group's initial wait is
// already over for this alert, it triggers an immediate flush.
func (ag *aggrGroup) insert(a *entity.NormalizedAlert) {
	ag.mtx.Lock()
	defer ag.mtx.Unlock()

	ag.alerts[a.Fingerprint] = a

	if !ag.hasFlushed && a.StartsAt.Add(ag.opts.GroupWait).Before(time.Now()) {
		ag.next.Reset(0)
	}
}

func (ag *aggrGroup) empty() bool {
	ag.mtx.Lock()
	defer ag.mtx.Unlock()
	return len(ag.alerts) == 0
}

// flush builds the current batch, decides whether it needs sending (changed set
// or repeat interval elapsed), delivers it, and on success evicts resolved
// alerts so the group can eventually be cleaned up.
func (ag *aggrGroup) flush(ctx context.Context, nf flushFunc) {
	ag.mtx.Lock()
	if len(ag.alerts) == 0 {
		ag.mtx.Unlock()
		return
	}

	batch := make([]*entity.NormalizedAlert, 0, len(ag.alerts))
	for _, a := range ag.alerts {
		batch = append(batch, a)
	}
	sort.Slice(batch, func(i, j int) bool { return batch[i].Fingerprint < batch[j].Fingerprint })

	digest := batchDigest(batch)
	now := time.Now()
	needsFlush := digest != ag.lastDigest || now.Sub(ag.lastFlush) >= ag.opts.RepeatInterval
	ag.mtx.Unlock()

	if !needsFlush {
		return
	}

	if !nf(ctx, batch) {
		// Delivery failed; keep everything and retry on the next tick.
		return
	}

	ag.mtx.Lock()
	ag.lastFlush = now
	ag.lastDigest = digest
	// Drop resolved alerts that haven't been replaced by a newer update since
	// we snapshotted the batch.
	for _, a := range batch {
		if a.Status != entity.StatusResolved {
			continue
		}
		if cur, ok := ag.alerts[a.Fingerprint]; ok && cur.UpdatedAt.Equal(a.UpdatedAt) {
			delete(ag.alerts, a.Fingerprint)
		}
	}
	ag.mtx.Unlock()
}

// batchDigest produces a stable fingerprint of a batch's identity + status so
// repeated identical groups can be suppressed between repeat intervals.
func batchDigest(batch []*entity.NormalizedAlert) string {
	var b strings.Builder
	for _, a := range batch {
		b.WriteString(a.Fingerprint)
		b.WriteByte('=')
		b.WriteString(string(a.Status))
		b.WriteByte(';')
	}
	return b.String()
}
