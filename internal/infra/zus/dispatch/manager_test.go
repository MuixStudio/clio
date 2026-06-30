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
	"testing"
	"time"

	"github.com/google/uuid"
	alertEntity "github.com/muixstudio/clio/internal/domain/alert/entity"
	"github.com/muixstudio/clio/internal/infra/bus"
)

func TestManagerConsumesAlertMessages(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	b := bus.New(8)
	defer b.Close()

	store := newManagerTestStore()
	mgr := startManagerTest(t, ctx, b, store)
	defer mgr.Stop()

	teamID := uuid.New()
	publisher := NewAlertPublisher(b.Pub)
	if err := publisher.PublishAlerts(ctx, []*alertEntity.NormalizedAlert{{
		TeamID:      teamID,
		Fingerprint: "fp-1",
		Status:      alertEntity.StatusFiring,
		Labels:      map[string]string{"alertname": "HighLatency"},
	}}); err != nil {
		t.Fatalf("PublishAlerts() error = %v", err)
	}

	if got := waitUUID(t, ctx, store.loadCh, "route tree load"); got != teamID {
		t.Fatalf("loaded team = %s, want %s", got, teamID)
	}
}

func TestManagerConsumesReloadCommands(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	b := bus.New(8)
	defer b.Close()

	store := newManagerTestStore()
	mgr := startManagerTest(t, ctx, b, store)
	defer mgr.Stop()

	teamID := uuid.New()
	publisher := NewReloadPublisher(b.Pub)
	if err := publisher.ReloadAlertRoutes(ctx, uuid.NullUUID{}, teamID); err != nil {
		t.Fatalf("ReloadAlertRoutes() error = %v", err)
	}

	if got := waitUUID(t, ctx, store.loadCh, "route tree reload"); got != teamID {
		t.Fatalf("reloaded team = %s, want %s", got, teamID)
	}
}

func startManagerTest(t *testing.T, ctx context.Context, b *bus.Bus, store *managerTestStore) *Manager {
	t.Helper()
	mgr, err := NewManager(b.Sub, store, NotifierFunc(func(context.Context, Notification) error { return nil }), nil)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}
	if err := mgr.Start(ctx); err != nil {
		t.Fatalf("Manager.Start() error = %v", err)
	}
	return mgr
}

// managerTestStore is a RouteStore that records each tree load. The Manager only
// reads now — CRUD lives in the HTTP handler — so a loader is all it needs.
type managerTestStore struct {
	tree   *Route
	loadCh chan uuid.UUID
}

func newManagerTestStore() *managerTestStore {
	return &managerTestStore{
		tree: &Route{
			RouteID:   uuid.New(),
			RouteOpts: DefaultRouteOpts,
		},
		loadCh: make(chan uuid.UUID, 4),
	}
}

func (s *managerTestStore) LoadAlertRouteTree(_ context.Context, _ uuid.NullUUID, teamID uuid.UUID) (*Route, error) {
	s.loadCh <- teamID
	return s.tree, nil
}

func waitUUID(t *testing.T, ctx context.Context, ch <-chan uuid.UUID, what string) uuid.UUID {
	t.Helper()
	select {
	case id := <-ch:
		return id
	case <-ctx.Done():
		t.Fatalf("timed out waiting for %s: %v", what, ctx.Err())
		return uuid.Nil
	}
}
