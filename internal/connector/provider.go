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
	"github.com/muixstudio/clio/internal/alert"
)

// ConnectorI is a configured instance of an alert source integration.
// Each concrete implementation handles one source type (e.g. Prometheus, Zabbix).
type ConnectorProvider interface {
	// Type returns the source type string, e.g. "prometheus".
	Type() string
	// Receive parses the raw HTTP payload and returns normalized Alerts.
	Normalize(raw []byte) ([]*alert.NormalizedAlert, error)
}

type ConnectorProviderProvider interface {
	ConnectorProviders() map[string]ConnectorProvider
}
