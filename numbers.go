// Copyright (c) 2020 Fernando G. Vilar
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package semver

// Numbers is a slice of Number.
type Numbers []Number

// Len is implemented so that Numbers satisfies sort.Interface
func (ns Numbers) Len() int { return len(ns) }

// Swap is implemented so that Numbers satisfies sort.Interface
func (ns Numbers) Swap(i, j int) { ns[i], ns[j] = ns[j], ns[i] }

// Less is implemented so that Numbers satisfies sort.Interface
func (ns Numbers) Less(i, j int) bool { return ns[i] < ns[j] }
