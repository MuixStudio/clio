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

package alert

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/infra/zus/dispatch"
	"github.com/prometheus/alertmanager/pkg/labels"
	"github.com/prometheus/common/model"
	"go.uber.org/zap"

	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/infra/response"
)

type routeMatcherRequest struct {
	Type  string `json:"type" validate:"required,oneof='=' '!=' '=~' '!~'"`
	Name  string `json:"name" validate:"required"`
	Value string `json:"value"`
}

type createRouteRequest struct {
	ParentID *uuid.UUID            `json:"parent_id"`
	Name     string                `json:"name" validate:"required,min=1,max=128"`
	Matchers []routeMatcherRequest `json:"matchers"`
	Continue bool                  `json:"continue"`

	Receiver            string   `json:"receiver"`
	GroupBy             []string `json:"group_by"`
	GroupByAll          bool     `json:"group_by_all"`
	GroupWait           *string  `json:"group_wait"`
	GroupInterval       *string  `json:"group_interval"`
	RepeatInterval      *string  `json:"repeat_interval"`
	MuteTimeIntervals   []string `json:"mute_time_intervals"`
	ActiveTimeIntervals []string `json:"active_time_intervals"`
}

type updateRouteRequest struct {
	Name     *string                `json:"name" validate:"omitempty,min=1,max=128"`
	Matchers *[]routeMatcherRequest `json:"matchers"`
	Continue *bool                  `json:"continue"`

	Receiver            *string   `json:"receiver"`
	GroupBy             *[]string `json:"group_by"`
	GroupByAll          *bool     `json:"group_by_all"`
	GroupWait           *string   `json:"group_wait"`
	GroupInterval       *string   `json:"group_interval"`
	RepeatInterval      *string   `json:"repeat_interval"`
	MuteTimeIntervals   *[]string `json:"mute_time_intervals"`
	ActiveTimeIntervals *[]string `json:"active_time_intervals"`
}

type routeItem struct {
	ID                  uuid.UUID      `json:"id"`
	Name                string         `json:"name"`
	Idx                 int            `json:"idx"`
	Matchers            []routeMatcher `json:"matchers"`
	Continue            bool           `json:"continue"`
	Receiver            string         `json:"receiver"`
	GroupBy             []string       `json:"group_by"`
	GroupByAll          bool           `json:"group_by_all"`
	GroupWait           string         `json:"group_wait"`
	GroupInterval       string         `json:"group_interval"`
	RepeatInterval      string         `json:"repeat_interval"`
	MuteTimeIntervals   []string       `json:"mute_time_intervals"`
	ActiveTimeIntervals []string       `json:"active_time_intervals"`
	Routes              []routeItem    `json:"routes"`
}

type routeMatcher struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type createRouteResponse struct {
	ID uuid.UUID `json:"id"`
}

var routeMatchTypeFromString = map[string]labels.MatchType{
	labels.MatchEqual.String():     labels.MatchEqual,
	labels.MatchNotEqual.String():  labels.MatchNotEqual,
	labels.MatchRegexp.String():    labels.MatchRegexp,
	labels.MatchNotRegexp.String(): labels.MatchNotRegexp,
}

func parseUUIDParam(c *gin.Context, param string, label string) (uuid.UUID, bool) {
	raw := c.Param(param)
	id, err := uuid.Parse(raw)
	if err != nil {
		response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid argument "+label+": "+raw))
		return uuid.Nil, false
	}
	return id, true
}

func parseRouteMatchers(req []routeMatcherRequest) (labels.Matchers, error) {
	matchers := make(labels.Matchers, 0, len(req))
	for _, rm := range req {
		mt, ok := routeMatchTypeFromString[rm.Type]
		if !ok {
			return nil, errors.BadRequest("INVALID_MATCHER", "invalid matcher type: "+rm.Type)
		}
		m, err := labels.NewMatcher(mt, rm.Name, rm.Value)
		if err != nil {
			return nil, errors.BadRequest("INVALID_MATCHER", "invalid matcher").WithCause(err)
		}
		matchers = append(matchers, m)
	}
	return matchers, nil
}

func parseDurationPtr(value *string, field string) (*time.Duration, error) {
	if value == nil {
		return nil, nil
	}
	d, err := time.ParseDuration(*value)
	if err != nil {
		return nil, errors.BadRequest("INVALID_ARGUMENT", "invalid duration "+field+": "+*value).WithCause(err)
	}
	return &d, nil
}

func toMatcherItems(matchers labels.Matchers) []routeMatcher {
	items := make([]routeMatcher, 0, len(matchers))
	for _, m := range matchers {
		items = append(items, routeMatcher{Type: m.Type.String(), Name: m.Name, Value: m.Value})
	}
	return items
}

func toGroupByNames(groupBy map[model.LabelName]struct{}) []string {
	items := make([]string, 0, len(groupBy))
	for name := range groupBy {
		items = append(items, string(name))
	}
	return items
}

func toRouteItem(route *dispatch.Route) routeItem {
	item := routeItem{
		ID:                  route.RouteID,
		Name:                route.Name,
		Idx:                 route.Idx,
		Matchers:            toMatcherItems(route.Matchers),
		Continue:            route.Continue,
		Receiver:            route.RouteOpts.Receiver,
		GroupBy:             toGroupByNames(route.RouteOpts.GroupBy),
		GroupByAll:          route.RouteOpts.GroupByAll,
		GroupWait:           route.RouteOpts.GroupWait.String(),
		GroupInterval:       route.RouteOpts.GroupInterval.String(),
		RepeatInterval:      route.RouteOpts.RepeatInterval.String(),
		MuteTimeIntervals:   route.RouteOpts.MuteTimeIntervals,
		ActiveTimeIntervals: route.RouteOpts.ActiveTimeIntervals,
		Routes:              make([]routeItem, 0, len(route.Routes)),
	}
	for _, child := range route.Routes {
		item.Routes = append(item.Routes, toRouteItem(child))
	}
	return item
}

// GetRouteTree handles GET /team/:team_id/alert-routes.
func (h *AlertHandler) GetRouteTree() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()
		teamID, ok := parseUUIDParam(c, "team_id", "team id")
		if !ok {
			return
		}

		tree, err := h.d.AlertRouteLoader().LoadAlertRouteTree(c.Request.Context(), uuid.NullUUID{}, teamID)
		if err != nil {
			log.Error("get alert route tree: failed", zap.Stringer("team_id", teamID), zap.Error(err))
			response.Fail(c, err)
			return
		}
		response.SuccessWithData(c, toRouteItem(tree))
	}
}

// CreateRoute handles POST /team/:team_id/alert-route.
func (h *AlertHandler) CreateRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()
		teamID, ok := parseUUIDParam(c, "team_id", "team id")
		if !ok {
			return
		}

		var req createRouteRequest
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			log.Warn("create alert route: invalid request", zap.Stringer("team_id", teamID), zap.Error(err))
			response.Fail(c, err)
			return
		}

		matchers, err := parseRouteMatchers(req.Matchers)
		if err != nil {
			response.Fail(c, err)
			return
		}
		groupWait, err := parseDurationPtr(req.GroupWait, "group_wait")
		if err != nil {
			response.Fail(c, err)
			return
		}
		groupInterval, err := parseDurationPtr(req.GroupInterval, "group_interval")
		if err != nil {
			response.Fail(c, err)
			return
		}
		repeatInterval, err := parseDurationPtr(req.RepeatInterval, "repeat_interval")
		if err != nil {
			response.Fail(c, err)
			return
		}

		in := dispatch.AddRouteInput{
			Name:                req.Name,
			Matchers:            matchers,
			Continue:            req.Continue,
			Receiver:            req.Receiver,
			GroupBy:             req.GroupBy,
			GroupByAll:          req.GroupByAll,
			GroupWait:           groupWait,
			GroupInterval:       groupInterval,
			RepeatInterval:      repeatInterval,
			MuteTimeIntervals:   req.MuteTimeIntervals,
			ActiveTimeIntervals: req.ActiveTimeIntervals,
		}
		if req.ParentID != nil {
			in.ParentID = uuid.NullUUID{UUID: *req.ParentID, Valid: true}
		}

		id, err := h.d.AlertRouteWriter().AddAlertRoute(c.Request.Context(), uuid.NullUUID{}, teamID, in)
		if err != nil {
			log.Error("create alert route: failed", zap.Stringer("team_id", teamID), zap.Error(err))
			response.Fail(c, err)
			return
		}
		h.reloadRoutes(c, teamID)
		log.Info("alert route created", zap.Stringer("team_id", teamID), zap.Stringer("route_id", id))
		response.SuccessWithData(c, createRouteResponse{ID: id})
	}
}

// UpdateRoute handles PATCH /team/:team_id/alert-route/:route_id.
func (h *AlertHandler) UpdateRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()
		teamID, ok := parseUUIDParam(c, "team_id", "team id")
		if !ok {
			return
		}
		routeID, ok := parseUUIDParam(c, "route_id", "route id")
		if !ok {
			return
		}

		var req updateRouteRequest
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			log.Warn("update alert route: invalid request", zap.Stringer("team_id", teamID), zap.Stringer("route_id", routeID), zap.Error(err))
			response.Fail(c, err)
			return
		}

		in := dispatch.UpdateRouteInput{
			RouteID:             routeID,
			Name:                req.Name,
			Continue:            req.Continue,
			Receiver:            req.Receiver,
			GroupBy:             req.GroupBy,
			GroupByAll:          req.GroupByAll,
			MuteTimeIntervals:   req.MuteTimeIntervals,
			ActiveTimeIntervals: req.ActiveTimeIntervals,
		}
		if req.Matchers != nil {
			matchers, err := parseRouteMatchers(*req.Matchers)
			if err != nil {
				response.Fail(c, err)
				return
			}
			in.Matchers = &matchers
		}
		var err error
		if in.GroupWait, err = parseDurationPtr(req.GroupWait, "group_wait"); err != nil {
			response.Fail(c, err)
			return
		}
		if in.GroupInterval, err = parseDurationPtr(req.GroupInterval, "group_interval"); err != nil {
			response.Fail(c, err)
			return
		}
		if in.RepeatInterval, err = parseDurationPtr(req.RepeatInterval, "repeat_interval"); err != nil {
			response.Fail(c, err)
			return
		}

		if err := h.d.AlertRouteWriter().UpdateAlertRoute(c.Request.Context(), uuid.NullUUID{}, teamID, in); err != nil {
			log.Error("update alert route: failed", zap.Stringer("team_id", teamID), zap.Stringer("route_id", routeID), zap.Error(err))
			response.Fail(c, err)
			return
		}
		h.reloadRoutes(c, teamID)
		log.Info("alert route updated", zap.Stringer("team_id", teamID), zap.Stringer("route_id", routeID))
		response.SuccessOK(c)
	}
}

// DeleteRoute handles DELETE /team/:team_id/alert-route/:route_id.
func (h *AlertHandler) DeleteRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()
		teamID, ok := parseUUIDParam(c, "team_id", "team id")
		if !ok {
			return
		}
		routeID, ok := parseUUIDParam(c, "route_id", "route id")
		if !ok {
			return
		}

		if err := h.d.AlertRouteWriter().DeleteAlertRoute(c.Request.Context(), uuid.NullUUID{}, teamID, routeID); err != nil {
			log.Error("delete alert route: failed", zap.Stringer("team_id", teamID), zap.Stringer("route_id", routeID), zap.Error(err))
			response.Fail(c, err)
			return
		}
		h.reloadRoutes(c, teamID)
		log.Info("alert route deleted", zap.Stringer("team_id", teamID), zap.Stringer("route_id", routeID))
		response.SuccessOK(c)
	}
}

// ReloadRoutes handles POST /team/:team_id/alert-routes/reload.
func (h *AlertHandler) ReloadRoutes() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()
		teamID, ok := parseUUIDParam(c, "team_id", "team id")
		if !ok {
			return
		}

		if err := h.d.AlertRouteReloader().ReloadAlertRoutes(c.Request.Context(), uuid.NullUUID{}, teamID); err != nil {
			log.Error("reload alert routes: failed", zap.Stringer("team_id", teamID), zap.Error(err))
			response.Fail(c, err)
			return
		}
		log.Info("alert route reload command published", zap.Stringer("team_id", teamID))
		response.SuccessOK(c)
	}
}

// reloadRoutes publishes a reload so the live Dispatcher picks up a just-applied
// CRUD change. The write is already persisted, so a publish failure is logged
// but does not fail the request — the next reload reconciles the dispatcher.
func (h *AlertHandler) reloadRoutes(c *gin.Context, teamID uuid.UUID) {
	if err := h.d.AlertRouteReloader().ReloadAlertRoutes(c.Request.Context(), uuid.NullUUID{}, teamID); err != nil {
		h.d.LoggerProvider().Error("publish alert route reload: failed",
			zap.Stringer("team_id", teamID), zap.Error(err))
	}
}
