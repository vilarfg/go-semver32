// Copyright (c) 2020 Fernando G. Vilar
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package semver_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"gopkg.in/yaml.v3"

	semver "github.com/vilarfg/go-semver32"
)

func TestErrorStrings(t *testing.T) {
	var produceError = func(s string) error {
		_, err := semver.ParseNumber(s)
		return err
	}

	var tcs = []struct {
		err error
		msg string
	}{
		{produceError(""), "semver: number representation is empty"},
		{produceError("0-1"), "semver: invalid character '-' in: \"0-1\""},
		{produceError("65536.0.0"), "semver: major component is too big: \"65536.0.0\""},
		{produceError("0.256.0"), "semver: minor component is too big: \"0.256.0\""},
		{produceError("0.0.256"), "semver: patch component is too big: \"0.0.256\""},
	}

	var ice semver.ErrorInvalidCharacter
	for i, tc := range tcs {
		if tc.err.Error() != tc.msg {
			t.Errorf("tc[%d] Parsed Number error message mismatch expected: %s got: %s", i, tc.msg, tc.err.Error())
		}

		if errors.As(tc.err, &ice) && ice.Character() != '-' {
			t.Errorf("tc[%d] Parsed Number invalid character error mismatch expected: %c got: %c", i, 'c', ice.Character())
		}
	}

}

func TestError(t *testing.T) {
	var tcs = []struct {
		input string
		err   error
	}{
		{"", semver.ErrorEmpty{}},
		{"0-1", semver.ErrorInvalidCharacter{}},
		{"0.0.256", semver.ErrorPatchTooBig("256")},
		{"0.0.257", semver.ErrorPatchTooBig("257")},
		{"0.256.0", semver.ErrorMinorTooBig("256")},
		{"0.256.256", semver.ErrorMinorTooBig("256")},
		{"0.257.0", semver.ErrorMinorTooBig("257")},
		{"0.257.256", semver.ErrorMinorTooBig("257")},
		{"65536.0.0", semver.ErrorMajorTooBig("65536")},
		{"65536.256.0", semver.ErrorMajorTooBig("65536")},
		{"65536.0.256", semver.ErrorMajorTooBig("65536")},
		{"65536.256.256", semver.ErrorMajorTooBig("65536")},
		{"65537.0.0", semver.ErrorMajorTooBig("65537")},
		{"65537.256.0", semver.ErrorMajorTooBig("65537")},
		{"65537.0.256", semver.ErrorMajorTooBig("65537")},
		{"65537.256.256", semver.ErrorMajorTooBig("65537")},
	}

	for i, tc := range tcs {
		if n, err := semver.ParseNumber(tc.input); n != 0 {
			t.Errorf("tc[%d] Parsed Number mismatch expected: 0 got: %s", i, n)
		} else if e := err.(*semver.Error).Unwrap(); fmt.Sprintf("%T", e) != fmt.Sprintf("%T", tc.err) {
			t.Errorf("tc[%d] Parsed Number error mismatch expected: %T got: %T", i, tc.err, e)
		}
	}
}

func TestUnmarshaling(t *testing.T) {
	var tcs = []struct {
		s   string
		err error
	}{
		{"", semver.ErrorEmpty{}},
		{"1-0", semver.ErrorInvalidCharacter{}},
		{"0.0.256", semver.ErrorPatchTooBig("0.0.256")},
		{"0.256.0", semver.ErrorMinorTooBig("0.256.0")},
		{"65536.0.0", semver.ErrorMajorTooBig("65536.0.0")},
	}

	for i, tc := range tcs {
		var n semver.Number
		if err := yaml.Unmarshal([]byte(tc.s), &n); err == nil {
			if len(tc.s) > 0 {
				t.Errorf("tc[%d] YAML unmarshaling error expected, got: nil", i)
			}
		} else {
			if e := err.(*semver.Error).Unwrap(); fmt.Sprintf("%T", e) != fmt.Sprintf("%T", tc.err) {
				t.Errorf("tc[%d] YAML unmarshaling error mismatch expected: %T got: %T", i, tc.err, e)
			} else {
				if ice, ok := e.(semver.ErrorInvalidCharacter); ok {
					if ice.Error() != `invalid character '-' in: "1-0"` {
						t.Errorf("tc[%d] YAML unmarshaling error message mismatch expected: %s got: %s", i, tc.err, e)
					}
				} else if e.Error() != tc.err.Error() {
					t.Errorf("tc[%d] YAML unmarshaling error message mismatch expected: %s got: %s", i, tc.err, e)
				}
			}
		}
	}

	for i, tc := range tcs {
		var n semver.Number
		if err := json.Unmarshal([]byte(`"`+tc.s+`"`), &n); err == nil {
			t.Errorf("tc[%d] JSON unmarshaling error expected, got: nil", i)
		} else {
			if e := err.(*semver.Error).Unwrap(); fmt.Sprintf("%T", e) != fmt.Sprintf("%T", tc.err) {
				t.Errorf("tc[%d] JSON unmarshaling error mismatch expected: %T got: %T", i, tc.err, e)
			} else {
				if ice, ok := e.(semver.ErrorInvalidCharacter); ok {
					if ice.Error() != `invalid character '-' in: "1-0"` {
						t.Errorf("tc[%d] JSON unmarshaling error message mismatch expected: %s got: %s", i, tc.err, e)
					}
				} else if e.Error() != tc.err.Error() {
					t.Errorf("tc[%d] JSON unmarshaling error message mismatch expected: %s got: %s", i, tc.err, e)
				}
			}
		}
	}

	for i, tc := range tcs {
		var n semver.Number
		if err := n.UnmarshalText([]byte(tc.s)); err == nil {
			t.Errorf("tc[%d] JSON unmarshaling error expected, got: nil", i)
		} else {
			if e := err.(*semver.Error).Unwrap(); fmt.Sprintf("%T", e) != fmt.Sprintf("%T", tc.err) {
				t.Errorf("tc[%d] JSON unmarshaling error mismatch expected: %T got: %T", i, tc.err, e)
			} else {
				if ice, ok := e.(semver.ErrorInvalidCharacter); ok {
					if ice.Error() != `invalid character '-' in: "1-0"` {
						t.Errorf("tc[%d] JSON unmarshaling error message mismatch expected: %s got: %s", i, tc.err, e)
					}
				} else if e.Error() != tc.err.Error() {
					t.Errorf("tc[%d] JSON unmarshaling error message mismatch expected: %s got: %s", i, tc.err, e)
				}
			}
		}
	}
}
