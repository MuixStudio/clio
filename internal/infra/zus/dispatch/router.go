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
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/alertmanager/pkg/labels"
	"github.com/prometheus/common/model"

	dbmodel "github.com/muixstudio/clio/internal/persistence/gormstore/model"
)

// DefaultRouteOpts are the defaulting routing options which apply
// to the root route of a routing tree.
var DefaultRouteOpts = RouteOpts{
	GroupWait:         30 * time.Second,
	GroupInterval:     5 * time.Minute,
	RepeatInterval:    4 * time.Hour,
	GroupBy:           map[model.LabelName]struct{}{},
	GroupByAll:        false,
	MuteTimeIntervals: []string{},
}

// A Route is a node that contains definitions of how to handle alerts.
type Route struct {
	parent *Route

	// RouteID is the persistent identifier of this route (the DB primary
	// key). Unlike Idx, it is stable across rebuilds and can be used to
	// address a specific node in the database. (The ID method below returns
	// a path-based string identifier and is a different concept.)
	RouteID uuid.UUID

	// Name is the human-readable route name persisted with the route row.
	Name string

	// The configuration parameters for matches of this route.
	RouteOpts RouteOpts

	// Matchers an alert has to fulfill to match
	// this route.
	Matchers labels.Matchers

	// If true, an alert matches further routes on the same level.
	Continue bool

	// Children routes of this route.
	Routes []*Route

	// Idx contains the post-order index of this route within the tree.
	// It is recomputed on every BuildTree call and is NOT persisted.
	Idx int
}

// BuildTree rebuilds an in-memory routing tree from a flat list of route rows
// loaded from the database. Exactly one row must be a root (ParentID is NULL);
// every other row must reference an existing parent. Siblings are ordered by
// their Priority column (lower matches first), and each node is assigned a
// post-order Idx after the whole tree is constructed.
func BuildTree(rows []*dbmodel.AlertRoute) (*Route, error) {
	childrenOf := make(map[uuid.UUID][]*dbmodel.AlertRoute, len(rows))
	var root *dbmodel.AlertRoute
	for _, m := range rows {
		if !m.ParentID.Valid {
			if root != nil {
				return nil, fmt.Errorf("router: multiple root routes found (%s and %s)", root.ID, m.ID)
			}
			root = m
			continue
		}
		childrenOf[m.ParentID.UUID] = append(childrenOf[m.ParentID.UUID], m)
	}
	if root == nil {
		return nil, fmt.Errorf("router: no root route found")
	}
	for _, cs := range childrenOf {
		sort.SliceStable(cs, func(i, j int) bool { return cs[i].Priority < cs[j].Priority })
	}

	counter := 0
	return buildNode(root, nil, childrenOf, &counter)
}

// buildNode converts a single AlertRoute row (plus its descendants) into a
// *Route, inheriting unset routing options from the parent the same way
// Alertmanager does.
func buildNode(m *dbmodel.AlertRoute, parent *Route, childrenOf map[uuid.UUID][]*dbmodel.AlertRoute, counter *int) (*Route, error) {
	// Create default and overwrite with the row's settings, inheriting from
	// the parent where the row leaves a field unset.
	opts := DefaultRouteOpts
	if parent != nil {
		opts = parent.RouteOpts
	}

	if m.Receiver != "" {
		opts.Receiver = m.Receiver
	}

	if groupBy, err := decodeLabelNames(m.GroupBy); err != nil {
		return nil, fmt.Errorf("router: route %s: invalid group_by: %w", m.ID, err)
	} else if groupBy != nil {
		opts.GroupBy = map[model.LabelName]struct{}{}
		for _, ln := range groupBy {
			opts.GroupBy[ln] = struct{}{}
		}
		opts.GroupByAll = false
	} else if m.GroupByAll {
		opts.GroupByAll = true
	}

	if m.GroupWait != nil {
		opts.GroupWait = time.Duration(*m.GroupWait)
	}
	if m.GroupInterval != nil {
		opts.GroupInterval = time.Duration(*m.GroupInterval)
	}
	if m.RepeatInterval != nil {
		opts.RepeatInterval = time.Duration(*m.RepeatInterval)
	}

	mutes, err := decodeStrings(m.MuteTimeIntervals)
	if err != nil {
		return nil, fmt.Errorf("router: route %s: invalid mute_time_intervals: %w", m.ID, err)
	}
	opts.MuteTimeIntervals = mutes
	actives, err := decodeStrings(m.ActiveTimeIntervals)
	if err != nil {
		return nil, fmt.Errorf("router: route %s: invalid active_time_intervals: %w", m.ID, err)
	}
	opts.ActiveTimeIntervals = actives

	matchers, err := decodeMatchers(m.Matchers)
	if err != nil {
		return nil, fmt.Errorf("router: route %s: invalid matchers: %w", m.ID, err)
	}
	sort.Sort(matchers)

	route := &Route{
		parent:    parent,
		RouteID:   m.ID,
		Name:      m.Name,
		RouteOpts: opts,
		Matchers:  matchers,
		Continue:  m.Continue,
	}

	// Build child routes first so they get lower (post-order) indices.
	for _, cm := range childrenOf[m.ID] {
		child, err := buildNode(cm, route, childrenOf, counter)
		if err != nil {
			return nil, err
		}
		route.Routes = append(route.Routes, child)
	}

	// Assign this node's index after all of its children have been indexed.
	route.Idx = *counter
	*counter++

	return route, nil
}

// Match does a depth-first left-to-right search through the route tree
// and returns the matching routing nodes.
func (r *Route) Match(lset model.LabelSet) []*Route {
	if !r.Matchers.Matches(lset) {
		return nil
	}

	var all []*Route

	for _, cr := range r.Routes {
		matches := cr.Match(lset)

		all = append(all, matches...)

		if matches != nil && !cr.Continue {
			break
		}
	}

	// If no child nodes were matches, the current node itself is a match.
	if len(all) == 0 {
		all = append(all, r)
	}

	return all
}

// Key returns a key for the route. It does not uniquely identify the route in general.
func (r *Route) Key() string {
	b := strings.Builder{}

	if r.parent != nil {
		b.WriteString(r.parent.Key())
		b.WriteRune('/')
	}
	b.WriteString(r.Matchers.String())
	return b.String()
}

// ID returns a unique identifier for the route.
func (r *Route) ID() string {
	b := strings.Builder{}

	if r.parent != nil {
		b.WriteString(r.parent.ID())
		b.WriteRune('/')
	}

	b.WriteString(r.Matchers.String())

	if r.parent != nil {
		for i := range r.parent.Routes {
			if r == r.parent.Routes[i] {
				b.WriteRune('/')
				b.WriteString(strconv.Itoa(i))
				break
			}
		}
	}

	return b.String()
}

// Walk traverses the route tree in depth-first order.
func (r *Route) Walk(visit func(*Route)) {
	visit(r)
	for i := range r.Routes {
		r.Routes[i].Walk(visit)
	}
}

// RouteOpts holds various routing options necessary for processing alerts
// that match a given route.
type RouteOpts struct {
	// The identifier of the associated notification configuration.
	Receiver string

	// What labels to group alerts by for notifications.
	GroupBy map[model.LabelName]struct{}

	// Use all alert labels to group.
	GroupByAll bool

	// How long to wait to group matching alerts before sending
	// a notification.
	GroupWait      time.Duration
	GroupInterval  time.Duration
	RepeatInterval time.Duration

	// A list of time intervals for which the route is muted.
	MuteTimeIntervals []string

	// A list of time intervals for which the route is active.
	ActiveTimeIntervals []string
}

func (ro *RouteOpts) String() string {
	var groupBy []model.LabelName
	for ln := range ro.GroupBy {
		groupBy = append(groupBy, ln)
	}
	return fmt.Sprintf("<RouteOpts send_to:%q group_by:%q group_by_all:%t timers:%q|%q>",
		ro.Receiver, groupBy, ro.GroupByAll, ro.GroupWait, ro.GroupInterval)
}

// MarshalJSON returns a JSON representation of the routing options.
func (ro *RouteOpts) MarshalJSON() ([]byte, error) {
	v := struct {
		Receiver       string           `json:"receiver"`
		GroupBy        model.LabelNames `json:"groupBy"`
		GroupByAll     bool             `json:"groupByAll"`
		GroupWait      time.Duration    `json:"groupWait"`
		GroupInterval  time.Duration    `json:"groupInterval"`
		RepeatInterval time.Duration    `json:"repeatInterval"`
	}{
		Receiver:       ro.Receiver,
		GroupByAll:     ro.GroupByAll,
		GroupWait:      ro.GroupWait,
		GroupInterval:  ro.GroupInterval,
		RepeatInterval: ro.RepeatInterval,
	}
	for ln := range ro.GroupBy {
		v.GroupBy = append(v.GroupBy, ln)
	}

	return json.Marshal(&v)
}
