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

package channel

import (
	channelDomain "github.com/muixstudio/clio/internal/channel"
	"github.com/muixstudio/clio/internal/logger"
)

type dependencies interface {
	channelDomain.ChannelPersisterProvider
	channelDomain.ChannelMemberPersisterProvider

	logger.Logger
}

// ChannelHandler handles channel management endpoints.
type ChannelHandler struct {
	d dependencies
}

// NewChannelHandler returns a ChannelHandler.
func NewChannelHandler(d dependencies) *ChannelHandler {
	return &ChannelHandler{d: d}
}
