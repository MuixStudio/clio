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

// Package bus is the in-process message transport seam between publishers and
// subscribers. Its single responsibility is moving messages — it carries no
// commands, request/reply, or business logic. The backing pub/sub (today an
// in-process Watermill GoChannel) is an implementation detail; swapping in a
// durable broker (Kafka, NATS, Redis, SQL) later only changes New, because the
// rest of the system depends only on the Publisher and Subscriber roles.
package bus

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

// Bus exposes the only two roles the rest of the system depends on: a Publisher
// to emit messages and a Subscriber to consume them.
type Bus struct {
	Pub message.Publisher
	Sub message.Subscriber

	gc *gochannel.GoChannel
}

// New creates an in-process bus. buffer is the per-subscriber output buffer; a
// small positive value smooths bursts without unbounded memory. Pass 0 for an
// unbuffered channel.
func New(buffer int64) *Bus {
	gc := gochannel.NewGoChannel(
		gochannel.Config{OutputChannelBuffer: buffer},
		watermill.NopLogger{},
	)
	// A GoChannel is both a Publisher and a Subscriber.
	return &Bus{Pub: gc, Sub: gc, gc: gc}
}

// Close shuts down the underlying pub/sub.
func (b *Bus) Close() error { return b.gc.Close() }
