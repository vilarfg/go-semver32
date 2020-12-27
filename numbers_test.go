// Copyright (c) 2020 Fernando G. Vilar
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package semver_test

import (
	"sort"
	"testing"

	semver "github.com/vilarfg/go-semver32"
)

func TestNumbers(t *testing.T) {
	numbers := semver.Numbers{
		semver.NewNumber(0, 1, 0),
		semver.NewNumber(0, 0, 0),
		semver.NewNumber(1, 1, 0),
	}

	if sort.IsSorted(numbers) {
		t.Error("expected numbers NOT to be sorted")
	}

	sort.Sort(numbers)

	if !sort.IsSorted(numbers) {
		t.Error("expected numbers to be sorted")
	}
}
