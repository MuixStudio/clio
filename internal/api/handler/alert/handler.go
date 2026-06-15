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
	"github.com/muixstudio/clio/internal/alert"
	"github.com/muixstudio/clio/internal/logger"
)

type dependencies interface {
	alert.AlertPersisterProvider

	logger.Logger
}

// AlertHandler handles alert query endpoints.
type AlertHandler struct {
	d dependencies
}

// NewAlertHandler returns an AlertHandler.
func NewAlertHandler(d dependencies) *AlertHandler {
	return &AlertHandler{d: d}
}
