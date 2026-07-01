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
	"errors"
	"sync"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/muixstudio/clio/internal/domain/alert/entity"
)

// teamInputBuffer is the per-team channel buffer feeding a Dispatcher. It
// absorbs short bursts; the dispatcher drains it quickly since insertion is just
// a map write and notification happens on separate goroutines.
const teamInputBuffer = 64

// ErrManagerStopped is returned by Route and the mutation methods once the
// Manager has been stopped.
var ErrManagerStopped = errors.New("dispatch: manager stopped")

// RouteStore loads a team's routing tree. *gormstore.GormStore satisfies it. The
// context must already carry the organization scope. Mutations no longer go
// through the Manager — the HTTP handler writes them and publishes a reload — so
// the Manager only needs read access.
type RouteStore interface {
	LoadAlertRouteTree(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID) (*Route, error)
}

// Manager owns the dispatch pipeline for alert events and route commands. It
// subscribes to the bus through a Watermill router, routes incoming alerts to one
// Dispatcher per team, and rebuilds a team's Dispatcher when route commands are
// consumed.
//
// Alert routing and route mutation are decoupled: route writes are published as
// fire-and-forget commands, then applied here asynchronously. Applying a command
// updates persistence and atomically swaps in a fresh Dispatcher, so a route
// change is eventually consistent with the alert stream — matching Alertmanager's
// config-reload semantics. Mutations are serialized so concurrent admin edits
// stay linearizable; alert routing never blocks on them.
type Manager struct {
	store    RouteStore
	notifier Notifier
	logger   *zap.Logger
	router   *message.Router

	// mutate serializes route mutations (DB write + Dispatcher swap) so two
	// concurrent edits can't leave a stale tree loaded.
	mutate sync.Mutex

	// mtx guards the team registry. It is held only for map reads/writes, never
	// across I/O, so alert routing stays fast.
	mtx   sync.Mutex
	teams map[uuid.UUID]*teamDispatcher

	ctx    context.Context
	cancel context.CancelFunc
}

type RouteChanged struct {
	OrganizationID uuid.NullUUID
	TeamID         uuid.UUID
}

type AlertCreated struct {
	OrganizationID uuid.NullUUID
	TeamID         uuid.UUID
	Alert          *entity.NormalizedAlert
}

// teamDispatcher is one team's running Dispatcher and its input channel. done is
// closed when the dispatcher is retired (on reload or shutdown) so in-flight
// Route callers stop trying to send to it instead of blocking or panicking.
type teamDispatcher struct {
	disp *Dispatcher
	in   chan *entity.NormalizedAlert
	done chan struct{}
}

// NewManager builds a Manager and registers its bus handlers. Call Start before
// publishing alerts or route commands to the in-process bus.
func NewManager(subscriber message.Subscriber, store RouteStore, notifier Notifier, logger *zap.Logger) (*Manager, error) {
	if logger == nil {
		logger = zap.NewNop()
	}
	logger = logger.With(zap.String("component", "dispatch-manager"))

	r, err := message.NewRouter(message.RouterConfig{}, watermill.NopLogger{})
	if err != nil {
		return nil, err
	}
	r.AddMiddleware(middleware.Recoverer)

	m := &Manager{
		store:    store,
		notifier: notifier,
		logger:   logger,
		router:   r,
		teams:    map[uuid.UUID]*teamDispatcher{},
	}

	// Alert routing stays with the Manager; route mutation is decoupled into a
	// dedicated handler that calls back into the Manager's mutation methods.
	eventProcessor, _ := cqrs.NewEventProcessorWithConfig(
		r,
		cqrs.EventProcessorConfig{
			GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				return params.EventName, nil
			},
			SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return subscriber, nil
			},
			//Marshaler: marshaler,
			//Logger:    logger,
		},
	)

	eventProcessor.AddHandlers(
		cqrs.NewEventHandler(
			"OnRouteChange",
			func(ctx context.Context, event *RouteChanged) error {
				m.mutate.Lock()
				defer m.mutate.Unlock()

				return m.reload(ctx, event.OrganizationID, event.TeamID)
			},
		),
		cqrs.NewEventHandler(
			"OnAlertCreate",
			func(ctx context.Context, event *AlertCreated) error {
				m.mutate.Lock()
				defer m.mutate.Unlock()

				return m.handleAlert(ctx, event.OrganizationID, event.TeamID, event.Alert)
			},
		),
	)

	return m, nil
}

// Start binds the Manager's lifetime to ctx, starts its bus router, and returns
// once the router has subscribed. Dispatchers created afterwards are torn down
// when ctx is cancelled or Stop is called.
func (m *Manager) Start(ctx context.Context) error {
	m.ctx, m.cancel = context.WithCancel(ctx)
	errc := make(chan error, 1)
	go func() { errc <- m.router.Run(m.ctx) }()

	select {
	case <-m.router.Running():
		return nil
	case err := <-errc:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop tears down the bus router and all team dispatchers. It is idempotent.
func (m *Manager) Stop() {
	if m == nil || m.cancel == nil {
		return
	}
	_ = m.router.Close()
	m.cancel()

	m.mtx.Lock()
	teams := m.teams
	m.teams = map[uuid.UUID]*teamDispatcher{}
	m.mtx.Unlock()

	for _, td := range teams {
		m.retire(td)
	}
}

// handleAlert decodes an alert message, rebuilds its org scope, and routes it.
// The alert is already persisted upstream, so a routing failure is logged and
// acked (return nil) rather than retried — the DB remains the source of truth.
func (m *Manager) handleAlert(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID, a *entity.NormalizedAlert) error {
	if err := m.Route(ctx, orgID, teamID, a); err != nil {
		m.logger.Error("dropping alert; routing failed",
			zap.String("team_id", teamID.String()),
			zap.String("fingerprint", a.Fingerprint), zap.Error(err))
	}
	return nil
}

// Route hands one alert to its team's Dispatcher, creating the Dispatcher (and
// loading its routing tree) on first use. It returns once the alert is queued;
// grouping and notification happen asynchronously inside the Dispatcher.
func (m *Manager) Route(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID, a *entity.NormalizedAlert) error {
	if m.ctx == nil {
		return ErrManagerStopped
	}
	td, err := m.dispatcherFor(ctx, orgID, teamID)
	if err != nil {
		return err
	}
	select {
	case td.in <- a:
		return nil
	case <-td.done:
		// Dispatcher was retired by a concurrent reload between lookup and send.
		// The alert is dropped; the DB remains the source of truth.
		return nil
	case <-m.ctx.Done():
		return ErrManagerStopped
	}
}

// reload loads the latest tree and atomically swaps in a fresh Dispatcher. The
// caller must hold m.mutate.
func (m *Manager) reload(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID) error {
	tree, err := m.store.LoadAlertRouteTree(ctx, orgID, teamID)
	if err != nil {
		return err
	}

	td := m.startDispatcher(teamID, tree)

	m.mtx.Lock()
	old := m.teams[teamID]
	m.teams[teamID] = td
	m.mtx.Unlock()

	m.retire(old)
	return nil
}

// dispatcherFor returns the team's Dispatcher, creating and starting it (and
// loading its routing tree) on first use.
func (m *Manager) dispatcherFor(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID) (*teamDispatcher, error) {
	m.mtx.Lock()
	if td, ok := m.teams[teamID]; ok {
		m.mtx.Unlock()
		return td, nil
	}
	m.mtx.Unlock()

	// Load outside the registry lock so a slow tree load doesn't stall routing
	// for other teams.
	tree, err := m.store.LoadAlertRouteTree(ctx, orgID, teamID)
	if err != nil {
		return nil, err
	}
	td := m.startDispatcher(teamID, tree)

	m.mtx.Lock()
	defer m.mtx.Unlock()
	// Another goroutine may have created it while we were loading; if so, keep
	// the established one and retire ours.
	if existing, ok := m.teams[teamID]; ok {
		m.retire(td)
		return existing, nil
	}
	m.teams[teamID] = td
	return td, nil
}

// startDispatcher builds and runs a Dispatcher for tree without registering it.
func (m *Manager) startDispatcher(teamID uuid.UUID, tree *Route) *teamDispatcher {
	in := make(chan *entity.NormalizedAlert, teamInputBuffer)
	disp := New(tree, m.notifier, m.logger.With(zap.String("team_id", teamID.String())))
	go disp.Run(in)
	return &teamDispatcher{disp: disp, in: in, done: make(chan struct{})}
}

// retire stops a (possibly nil) team dispatcher and signals in-flight Route
// callers to stop sending to it. The input channel is deliberately left unclosed
// — senders select on done, and the Dispatcher exits via Stop — so a racing send
// can never hit a closed channel.
func (m *Manager) retire(td *teamDispatcher) {
	if td == nil {
		return
	}
	close(td.done)
	td.disp.Stop()
}
