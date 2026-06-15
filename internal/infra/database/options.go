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

package database

type (
	Option func(*options)

	options struct {
		limit  int
		offset int
		orders []string
	}
	sortOrder string
)

var DESC sortOrder = "desc"
var ACS sortOrder = "acs"

func NewOptions() *options {
	return &options{
		limit:  10,
		offset: -1,
	}
}

func WithLimit(limit int) Option {
	return func(opt *options) {
		opt.limit = limit
	}
}

func WithOffset(offset int) Option {
	return func(opt *options) {
		opt.offset = offset
	}
}

func WithOrder(field string, r sortOrder) Option {
	return func(opt *options) {
		opt.orders = append(opt.orders, field+" "+string(r))
	}
}
