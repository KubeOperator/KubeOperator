/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package allocator

// Interface manages the allocation of items out of a range. Interface
// should be threadsafe.
type Interface interface {
	Allocate(int) (bool, error)
	AllocateNext() (int, bool, error)
	Release(int) error
	ForEach(func(int))

	// For testing
	Has(int) bool

	// For testing
	Free() int
}

// Snapshottable is an Interface that can be snapshotted and restored. Snapshottable
// should be threadsafe.
type Snapshottable interface {
	Interface
	Snapshot() (string, []byte)
	Restore(string, []byte) error
}

type Factory func(max int, rangeSpec string) Interface
