// Copyright (c) 2020 Fernando G. Vilar
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package semver

import "fmt"

// Error wraps all errors produced within this package
// so that they can be easily identified by its consumers.
type Error struct{ e error }

// Unwrap returns the error Error is wrapping.
func (e Error) Unwrap() error { return e.e }

// Error satisfies the error interface.
func (e Error) Error() string { return "semver: " + e.e.Error() }

// ErrorEmpty is an error to signal that the representation of the Number
// doesn't contain enough information for it to be parsed.
type ErrorEmpty struct{}

// Error satisfies the error interface.
func (e ErrorEmpty) Error() string { return "number representation is empty" }

// ErrorInvalidCharacter is an error to signal that the
// representation of the Number contains an invalid character.
type ErrorInvalidCharacter struct {
	s string
	c byte
}

// Error satisfies the error interface.
func (e ErrorInvalidCharacter) Error() string {
	return fmt.Sprintf("invalid character '%c' in: \"%s\"", e.c, e.s)
}

// Character returns the invalid character that
// caused the ErrorInvalidCharacter error.
func (e ErrorInvalidCharacter) Character() byte { return e.c }

// ErrorMajorTooBig is an error to signal that the major component of the
// Number is out of bounds (too big).
type ErrorMajorTooBig string

// Error satisfies the error interface.
func (e ErrorMajorTooBig) Error() string {
	return "major component is too big: \"" + string(e) + "\""
}

// ErrorMinorTooBig is an error to signal that the minor component of the
// Number is out of bounds (too big).
type ErrorMinorTooBig string

// Error satisfies the error interface.
func (e ErrorMinorTooBig) Error() string {
	return "minor component is too big: \"" + string(e) + "\""
}

// ErrorPatchTooBig is an error to signal that the patch component of the
// Number is out of bounds (too big).
type ErrorPatchTooBig string

// Error satisfies the error interface.
func (e ErrorPatchTooBig) Error() string {
	return "patch component is too big: \"" + string(e) + "\""
}

var (
	errorEmpty       = &Error{ErrorEmpty{}}
	errorMajorTooBig = &Error{ErrorMajorTooBig("65536")}
	errorMinorTooBig = &Error{ErrorMinorTooBig("256")}
	errorPatchTooBig = &Error{ErrorPatchTooBig("256")}
)
