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

	"github.com/prometheus/alertmanager/pkg/labels"
	"github.com/prometheus/common/model"
)

// jsonMatcher is the on-disk representation of a single label matcher.
// Type is one of "=", "!=", "=~", "!~" (labels.MatchType.String()).
type jsonMatcher struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

var matchTypeFromString = map[string]labels.MatchType{
	labels.MatchEqual.String():     labels.MatchEqual,
	labels.MatchNotEqual.String():  labels.MatchNotEqual,
	labels.MatchRegexp.String():    labels.MatchRegexp,
	labels.MatchNotRegexp.String(): labels.MatchNotRegexp,
}

// decodeMatchers turns the JSON matcher array stored on a route row into a
// labels.Matchers. A nil/empty payload yields an empty (always-matching)
// matcher set, which is the expected behaviour for the root route.
func decodeMatchers(raw []byte) (labels.Matchers, error) {
	if len(raw) == 0 {
		return labels.Matchers{}, nil
	}
	var jms []jsonMatcher
	if err := json.Unmarshal(raw, &jms); err != nil {
		return nil, err
	}
	matchers := make(labels.Matchers, 0, len(jms))
	for _, jm := range jms {
		mt, ok := matchTypeFromString[jm.Type]
		if !ok {
			return nil, fmt.Errorf("unknown matcher type %q", jm.Type)
		}
		m, err := labels.NewMatcher(mt, jm.Name, jm.Value)
		if err != nil {
			return nil, err
		}
		matchers = append(matchers, m)
	}
	return matchers, nil
}

// EncodeMatchers serializes a labels.Matchers into the JSON form stored in the
// alert_routes.matchers column. It is exported so callers building rows (e.g.
// the persistence layer) produce a representation decodeMatchers understands.
func EncodeMatchers(matchers labels.Matchers) ([]byte, error) {
	jms := make([]jsonMatcher, 0, len(matchers))
	for _, m := range matchers {
		jms = append(jms, jsonMatcher{Type: m.Type.String(), Name: m.Name, Value: m.Value})
	}
	return json.Marshal(jms)
}

// decodeLabelNames decodes a JSON array of label names. A nil/empty payload
// returns nil so the caller can treat it as "inherit from parent".
func decodeLabelNames(raw []byte) ([]model.LabelName, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var names []model.LabelName
	if err := json.Unmarshal(raw, &names); err != nil {
		return nil, err
	}
	if len(names) == 0 {
		return nil, nil
	}
	return names, nil
}

// decodeStrings decodes a JSON array of strings (e.g. mute/active time interval
// names). A nil/empty payload returns nil.
func decodeStrings(raw []byte) ([]string, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var out []string
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out, nil
}
