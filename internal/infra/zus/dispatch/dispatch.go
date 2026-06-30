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

// Package dispatch implements a minimal, production-ready alert dispatcher
// modeled on Prometheus Alertmanager. It matches incoming alerts against a
// routing tree, sorts them into aggregation groups, and flushes each group to a
// Notifier on the route's GroupWait/GroupInterval/RepeatInterval schedule.
package dispatch

import (
	"context"
	"sync"
	"time"

	"github.com/muixstudio/clio/internal/domain/alert/entity"
	"github.com/prometheus/common/model"
	"go.uber.org/zap"
)

// cleanupInterval is how often the dispatcher reaps empty aggregation groups.
const cleanupInterval = 30 * time.Second

// Dispatcher sorts incoming alerts into aggregation groups per matched route
// and drives their notification schedule. A Dispatcher is bound to a single
// routing tree; rebuild the tree and create a new Dispatcher on config reload.
type Dispatcher struct {
	route    *Route
	notifier Notifier
	logger   *zap.Logger
	timeout  func(time.Duration) time.Duration

	mtx    sync.Mutex
	groups map[*Route]map[model.Fingerprint]*aggrGroup

	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

// Option customizes a Dispatcher.
type Option func(*Dispatcher)

// WithFlushTimeout overrides how a GroupInterval is turned into the per-flush
// context deadline. Mainly useful in tests; defaults to the identity function.
func WithFlushTimeout(f func(time.Duration) time.Duration) Option {
	return func(d *Dispatcher) { d.timeout = f }
}

// New returns a Dispatcher that routes alerts through route and delivers
// flushed groups to notifier.
func New(route *Route, notifier Notifier, logger *zap.Logger, opts ...Option) *Dispatcher {
	if logger == nil {
		logger = zap.NewNop()
	}
	d := &Dispatcher{
		route:    route,
		notifier: notifier,
		logger:   logger.With(zap.String("component", "dispatcher")),
		timeout:  func(dur time.Duration) time.Duration { return dur },
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// Run consumes alerts from in until the channel is closed or Stop is called,
// blocking for the dispatcher's lifetime. It is safe to call Run at most once.
func (d *Dispatcher) Run(in <-chan *entity.NormalizedAlert) {
	d.mtx.Lock()
	d.groups = map[*Route]map[model.Fingerprint]*aggrGroup{}
	d.ctx, d.cancel = context.WithCancel(context.Background())
	d.done = make(chan struct{})
	d.mtx.Unlock()

	defer close(d.done)

	cleanup := time.NewTicker(cleanupInterval)
	defer cleanup.Stop()

	for {
		select {
		case a, ok := <-in:
			if !ok {
				return
			}
			d.dispatch(a)

		case <-cleanup.C:
			d.reapEmptyGroups()

		case <-d.ctx.Done():
			return
		}
	}
}

// Stop cancels all in-flight notifications and group schedules, then waits for
// Run to return. It is safe to call on a nil or never-started Dispatcher.
func (d *Dispatcher) Stop() {
	if d == nil {
		return
	}
	d.mtx.Lock()
	if d.cancel == nil {
		d.mtx.Unlock()
		return
	}
	d.cancel()
	d.cancel = nil
	done := d.done
	d.mtx.Unlock()

	<-done
}

// dispatch matches an alert against the routing tree and inserts it into the
// aggregation group of every matched route.
func (d *Dispatcher) dispatch(a *entity.NormalizedAlert) {
	lset := labelSet(a.Labels)
	for _, r := range d.route.Match(lset) {
		d.processAlert(a, r)
	}
}

func (d *Dispatcher) processAlert(a *entity.NormalizedAlert, route *Route) {
	gl := groupLabels(a, route)
	fp := gl.Fingerprint()

	d.mtx.Lock()
	routeGroups, ok := d.groups[route]
	if !ok {
		routeGroups = map[model.Fingerprint]*aggrGroup{}
		d.groups[route] = routeGroups
	}

	ag, ok := routeGroups[fp]
	if ok {
		d.mtx.Unlock()
		ag.insert(a)
		return
	}

	ag = newAggrGroup(d.ctx, gl, route, d.timeout, d.logger)
	routeGroups[fp] = ag
	d.mtx.Unlock()

	// Insert the first alert before starting run() so the initial flush sees it.
	ag.insert(a)

	go ag.run(func(ctx context.Context, alerts []*entity.NormalizedAlert) bool {
		n := Notification{
			Receiver:       route.RouteOpts.Receiver,
			GroupKey:       ag.groupKey,
			GroupLabels:    ag.labels,
			RepeatInterval: route.RouteOpts.RepeatInterval,
			Alerts:         alerts,
		}
		if err := d.notifier.Notify(ctx, n); err != nil {
			if ctx.Err() == nil {
				d.logger.Error("notify failed",
					zap.String("group_key", ag.groupKey),
					zap.Int("num_alerts", len(alerts)),
					zap.Error(err))
			}
			return false
		}
		return true
	})
}

// reapEmptyGroups stops and removes aggregation groups that no longer hold any
// alerts, releasing their goroutines and timers.
func (d *Dispatcher) reapEmptyGroups() {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	for route, groups := range d.groups {
		for fp, ag := range groups {
			if ag.empty() {
				ag.stop()
				delete(groups, fp)
			}
		}
		if len(groups) == 0 {
			delete(d.groups, route)
		}
	}
}

// groupLabels extracts the subset of an alert's labels used to key its
// aggregation group, per the route's GroupBy / GroupByAll options.
func groupLabels(a *entity.NormalizedAlert, route *Route) model.LabelSet {
	gl := model.LabelSet{}
	for k, v := range a.Labels {
		ln := model.LabelName(k)
		if _, ok := route.RouteOpts.GroupBy[ln]; ok || route.RouteOpts.GroupByAll {
			gl[ln] = model.LabelValue(v)
		}
	}
	return gl
}

// labelSet converts a domain label map into a prometheus LabelSet for matching.
func labelSet(m map[string]string) model.LabelSet {
	ls := make(model.LabelSet, len(m))
	for k, v := range m {
		ls[model.LabelName(k)] = model.LabelValue(v)
	}
	return ls
}
