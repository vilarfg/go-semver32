// Copyright (c) 2020 Fernando G. Vilar
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package semver_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"gopkg.in/yaml.v3"

	semver "github.com/vilarfg/go-semver32"
)

func TestNumber(t *testing.T) {
	var (
		eMtb = semver.ErrorMajorTooBig("65536")
		emtb = semver.ErrorMinorTooBig("256")
		eptb = semver.ErrorPatchTooBig("256")
	)

	var tcs = []struct {
		M    uint16
		m, p byte
		bumpMError,
		bumpmError,
		bumppError error
		str, goStr string
	}{
		{0, 0, 0, nil, nil, nil, "0", "0.0.0"},
		{0, 0, 1, nil, nil, nil, "0.0.1", "0.0.1"},
		{0, 0, 255, nil, nil, eptb, "0.0.255", "0.0.255"},
		{0, 1, 0, nil, nil, nil, "0.1", "0.1.0"},
		{0, 1, 1, nil, nil, nil, "0.1.1", "0.1.1"},
		{0, 1, 255, nil, nil, eptb, "0.1.255", "0.1.255"},
		{0, 255, 0, nil, emtb, nil, "0.255", "0.255.0"},
		{0, 255, 1, nil, emtb, nil, "0.255.1", "0.255.1"},
		{0, 255, 255, nil, emtb, eptb, "0.255.255", "0.255.255"},
		{1, 0, 0, nil, nil, nil, "1", "1.0.0"},
		{1, 0, 1, nil, nil, nil, "1.0.1", "1.0.1"},
		{1, 0, 255, nil, nil, eptb, "1.0.255", "1.0.255"},
		{1, 1, 0, nil, nil, nil, "1.1", "1.1.0"},
		{1, 1, 1, nil, nil, nil, "1.1.1", "1.1.1"},
		{1, 1, 255, nil, nil, eptb, "1.1.255", "1.1.255"},
		{1, 255, 0, nil, emtb, nil, "1.255", "1.255.0"},
		{1, 255, 1, nil, emtb, nil, "1.255.1", "1.255.1"},
		{1, 255, 255, nil, emtb, eptb, "1.255.255", "1.255.255"},
		{65535, 0, 0, eMtb, nil, nil, "65535", "65535.0.0"},
		{65535, 0, 1, eMtb, nil, nil, "65535.0.1", "65535.0.1"},
		{65535, 0, 255, eMtb, nil, eptb, "65535.0.255", "65535.0.255"},
		{65535, 1, 0, eMtb, nil, nil, "65535.1", "65535.1.0"},
		{65535, 1, 1, eMtb, nil, nil, "65535.1.1", "65535.1.1"},
		{65535, 1, 255, eMtb, nil, eptb, "65535.1.255", "65535.1.255"},
		{65535, 255, 0, eMtb, emtb, nil, "65535.255", "65535.255.0"},
		{65535, 255, 1, eMtb, emtb, nil, "65535.255.1", "65535.255.1"},
		{65535, 255, 255, eMtb, emtb, eptb, "65535.255.255", "65535.255.255"},
	}

	for i, tc := range tcs {
		n := semver.NewNumber(tc.M, tc.m, tc.p)

		if pn, err := semver.ParseNumber(tc.str); err == nil {
			if pn != n {
				t.Errorf("tc[%d] Parsed Number mismatch expected: %s got: %s", i, tc.str, pn)
			}
		} else {
			t.Errorf("tc[%d] No Number parsing error expected, got: %s", i, err.Error())
		}

		if pn, err := semver.ParseNumber(tc.goStr); err == nil {
			if pn != n {
				t.Errorf("tc[%d] Parsed Number mismatch expected: %s got: %s", i, tc.str, pn)
			}
		} else {
			t.Errorf("tc[%d] No Number parsing error expected, got: %s", i, err.Error())
		}

		if M := n.Major(); M != tc.M {
			t.Errorf("tc[%d] Major component mismatch expected: %d got: %d", i, tc.M, M)
		}
		if m := n.Minor(); m != tc.m {
			t.Errorf("tc[%d] Minor component mismatch expected: %d got: %d", i, tc.m, m)
		}
		if p := n.Patch(); p != tc.p {
			t.Errorf("tc[%d] Patch component mismatch expected: %d got: %d", i, tc.p, p)
		}

		if sM := n.SetMajor(2).Major(); sM != 2 {
			t.Errorf("tc[%d] Set Major component mismatch expected: %d got: %d", i, 2, sM)
		}
		if sm := n.SetMinor(2).Minor(); sm != 2 {
			t.Errorf("tc[%d] Set Minor component mismatch expected: %d got: %d", i, 2, sm)
		}
		if sp := n.SetPatch(2).Patch(); sp != 2 {
			t.Errorf("tc[%d] Set Patch component mismatch expected: %d got: %d", i, 2, sp)
		}

		if nn, err := n.BumpMajor(); err == nil {
			if M := nn.Major(); M != tc.M+1 {
				t.Errorf("tc[%d] Bumped Major component mismatch expected: %d got: %d", i, tc.M+1, M)
			}
			if m := nn.Minor(); m != 0 {
				t.Errorf("tc[%d] Bumped Minor component mismatch expected: 0 got: %d", i, m)
			}
			if p := nn.Patch(); p != 0 {
				t.Errorf("tc[%d] Bumped Patch component mismatch expected: 0 got: %d", i, p)
			}
			if tc.bumpMError != nil {
				t.Errorf("tc[%d] Bumped Major error mismatch expected: %s got: nil", i, tc.bumpMError.Error())
			}
		} else if tc.bumpMError == nil {
			t.Errorf("tc[%d] Bumped Major error mismatch expected: nil got: %s", i, err.Error())
		} else {
			var te semver.ErrorMajorTooBig
			if ok := errors.As(err, &te); !ok {
				t.Errorf("tc[%d] Bumped Major error mismatch expected error type: %T got: %s", i, tc.bumpMError, err.Error())
			} else if tc.bumpMError != te {
				t.Errorf("tc[%d] Bumped Major error mismatch expected: %s got: %s", i, tc.bumpMError.Error(), err.Error())
			}
		}

		if nn, err := n.BumpMinor(); err == nil {
			if M := nn.Major(); M != tc.M {
				t.Errorf("tc[%d] Bumped Major component mismatch expected: %d got: %d", i, tc.M, M)
			}
			if m := nn.Minor(); m != tc.m+1 {
				t.Errorf("tc[%d] Bumped Minor component mismatch expected: %d got: %d", i, tc.m+1, m)
			}
			if p := nn.Patch(); p != 0 {
				t.Errorf("tc[%d] Bumped Patch component mismatch expected: 0 got: %d", i, p)
			}
			if tc.bumpmError != nil {
				t.Errorf("tc[%d] Bumped Major error mismatch expected: %s got: nil", i, tc.bumpmError.Error())
			}
		} else if tc.bumpmError == nil {
			t.Errorf("tc[%d] Bumped Major error mismatch expected: nil got: %s", i, err.Error())
		} else {
			var te semver.ErrorMinorTooBig
			if ok := errors.As(err, &te); !ok {
				t.Errorf("tc[%d] Bumped Minor error mismatch expected error type: %T got: %s", i, tc.bumpmError, err.Error())
			} else if tc.bumpmError != te {
				t.Errorf("tc[%d] Bumped Minor error mismatch expected: %s got: %s", i, tc.bumpmError.Error(), err.Error())
			}
		}

		if nn, err := n.BumpPatch(); err == nil {
			if M := nn.Major(); M != tc.M {
				t.Errorf("tc[%d] Bumped Major component mismatch expected: %d got: %d", i, tc.M, M)
			}
			if m := nn.Minor(); m != tc.m {
				t.Errorf("tc[%d] Bumped Minor component mismatch expected: %d got: %d", i, tc.m, m)
			}
			if p := nn.Patch(); p != tc.p+1 {
				t.Errorf("tc[%d] Bumped Patch component mismatch expected: %d got: %d", i, tc.p+1, p)
			}
			if tc.bumppError != nil {
				t.Errorf("tc[%d] Bumped Major error mismatch expected: %s got: nil", i, tc.bumppError.Error())
			}
		} else if tc.bumppError == nil {
			t.Errorf("tc[%d] Bumped Major error mismatch expected: nil got: %s", i, err.Error())
		} else {
			var te semver.ErrorPatchTooBig
			if ok := errors.As(err, &te); !ok {
				t.Errorf("tc[%d] Bumped Patch error mismatch expected error type: %T got: %s", i, tc.bumppError, err.Error())
			} else if tc.bumppError != te {
				t.Errorf("tc[%d] Bumped Patch error mismatch expected: %s got: %s", i, tc.bumppError.Error(), err.Error())
			}
		}

		if s := n.String(); s != tc.str {
			t.Errorf("tc[%d] string mismatch expected: %s got: %s", i, tc.str, s)
		}
		if s := n.GoString(); s != tc.goStr {
			t.Errorf("tc[%d] gostring mismatch expected: %s got: %s", i, tc.goStr, s)
		}

		if s, err := n.MarshalYAML(); err == nil {
			if got := s.(string); tc.str != got {
				t.Errorf("tc[%d] yaml mismatch expected %s got: %s", i, tc.str, got)
			} else {
				var un semver.Number

				if err := yaml.Unmarshal([]byte(got), &un); err != nil {
					t.Errorf("tc[%d] no yaml unmarshaling error expected got: %s", i, err.Error())
				} else if n != un {
					t.Errorf("tc[%d] yaml unmarshaling mismatch expected %s got: %s from %s", i, n, un, got)
				}
			}
		} else {
			t.Errorf("tc[%d] no yaml marshaling error expected got: %s", i, err.Error())
		}

		if b, err := n.MarshalJSON(); err == nil {
			if exp, got := strconv.Quote(tc.str), string(b); exp != got {
				t.Errorf("tc[%d] json mismatch expected %s got: %s", i, exp, got)
			} else {
				var un semver.Number

				if err := un.UnmarshalJSON(b); err != nil {
					t.Errorf("tc[%d] no json unmarshaling error expected got: %s", i, err.Error())
				} else if n != un {
					t.Errorf("tc[%d] json unmarshaling mismatch expected %s got: %s from %s", i, n, un, string(b))
				}
			}
		} else {
			t.Errorf("tc[%d] no json marshaling error expected got: %s", i, err.Error())
		}

		if b, err := n.MarshalText(); err == nil {
			if exp, got := tc.str, string(b); exp != got {
				t.Errorf("tc[%d] text mismatch expected %s got: %s", i, exp, got)
			} else {
				var un semver.Number

				if err := un.UnmarshalText(b); err != nil {
					t.Errorf("tc[%d] no text unmarshaling error expected got: %s", i, err.Error())
				} else if n != un {
					t.Errorf("tc[%d] text unmarshaling mismatch expected %s got: %s from %s", i, n, un, string(b))
				}

			}
		} else {
			t.Errorf("tc[%d] no text marshaling error expected got: %s", i, err.Error())
		}
	}
}

func ExampleNewNumber() {
	n := semver.NewNumber(0, 1, 0)
	fmt.Printf("%d => %s", n, n)
	// Output: 256 => 0.1
}

func ExampleNumber_SetMajor() {
	n := semver.NewNumber(0, 1, 0)
	n = n.SetMajor(1)

	fmt.Printf("%d is %s", n, n)
	// Output: 65792 is 1.1
}

func ExampleNumber_SetMinor() {
	n := semver.NewNumber(0, 1, 0)
	n = n.SetMinor(2)

	fmt.Printf("%d is %s", n, n)
	// Output: 512 is 0.2
}

func ExampleNumber_SetPatch() {
	n := semver.NewNumber(0, 1, 0)
	n = n.SetPatch(1)

	fmt.Printf("%d is %s", n, n)
	// Output: 257 is 0.1.1
}

func ExampleNumber_BumpMajor() {
	n := semver.NewNumber(0, 1, 1)
	n, err := n.BumpMajor()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%d is %s", n, n)
	}
	// Output: 65536 is 1
}

func ExampleNumber_BumpMinor() {
	n := semver.NewNumber(0, 1, 1)
	n, err := n.BumpMinor()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%d is %s", n, n)
	}
	// Output: 512 is 0.2
}

func ExampleNumber_BumpPatch() {
	n := semver.NewNumber(0, 1, 1)
	n, err := n.BumpPatch()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%d is %s", n, n)
	}
	// Output: 258 is 0.1.2
}

func ExampleNumber_Major() {
	n := semver.NewNumber(1, 2, 3)

	fmt.Println(n.Major())
	// Output: 1
}

func ExampleNumber_Minor() {
	n := semver.NewNumber(1, 2, 3)

	fmt.Println(n.Minor())
	// Output: 2
}

func ExampleNumber_Patch() {
	n := semver.NewNumber(1, 2, 3)

	fmt.Println(n.Patch())
	// Output: 3
}

func ExampleNumber_String() {
	fmt.Printf("%s", semver.Number(768))
	// Output: 0.3
}

func ExampleNumber_GoString() {
	fmt.Printf("%#v", semver.Number(768))
	// Output: 0.3.0
}

func ExampleParseNumber() {
	n, err := semver.ParseNumber("1.2")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%d is %s", n, n)
	}
	// Output 258 => 1.2
}
