// Copyright (c) 2020 Fernando G. Vilar
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package semver

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Number represents a SemVer number capable of storing
// major components from 0 to 65,535,
// minor components from 0 to 255 and
// patch components from 0 to 255.
// No prerelease or build information is included.
type Number uint32

// SetMajor returns a new Number with
// the major component set to the specified value.
// The minor and patch components remain the same.
func (n Number) SetMajor(major uint16) Number { return n&invMajorMask | Number(major)<<16 }

// SetMinor returns a new Number with
// the minor component set to the specified value.
// The major and patch components remain the same.
func (n Number) SetMinor(minor byte) Number { return n&invMinorMask | Number(minor)<<8 }

// SetPatch returns a new Number with
// the patch component set to the specified value.
// The major and minor components remain the same.
func (n Number) SetPatch(patch byte) Number { return n&invPatchMask | Number(patch) }

// BumpMajor returns a new Number with the major component increased by 1.
// The minor and patch components are set to 0.
//
// It will produce an error if the resulting value were to be out of bounds.
func (n Number) BumpMajor() (Number, error) {
	if n.Major() == 65535 {
		return 0, errorMajorTooBig
	}
	return Number(n.Major()+1) << 16, nil
}

// BumpMinor returns a new Number with the minor component increased by 1.
// The major component remains unaffected, but the patch components is set to 0.
//
// It will produce an error if the resulting value were to be out of bounds.
func (n Number) BumpMinor() (Number, error) {
	if n.Minor() == 255 {
		return 0, errorMinorTooBig
	}
	return n&majorMask | Number(n.Minor()+1)<<8, nil
}

// BumpPatch returns a new Number with the patch component increased by 1.
// The major and minor components remain unaffected.
//
// It will produce an error if the resulting value were to be out of bounds.
func (n Number) BumpPatch() (Number, error) {
	if n.Patch() == 255 {
		return 0, errorPatchTooBig
	}
	return n&invPatchMask | Number(n.Patch()+1), nil
}

// Major returns the major component of the Number.
func (n Number) Major() uint16 { return uint16(n >> 16 & invMajorMask) }

// Minor returns the major component of the Number.
func (n Number) Minor() byte { return byte(n >> 8 & patchMask) }

// Patch returns the major component of the Number.
func (n Number) Patch() byte { return byte(n & patchMask) }

// String satisfies the fmt.Stringer interface.
func (n Number) String() string {
	var b strings.Builder
	write(&b, n)
	return b.String()
}

// GoString satisfies the fmt.GoStringer interface.
func (n Number) GoString() string {
	return fmt.Sprintf("%d.%d.%d", n>>16&invMajorMask, n>>8&patchMask, n&patchMask)
}

// MarshalYAML satisfies the gopkg.in/yaml.v3.Marshaler interface.
func (n Number) MarshalYAML() (interface{}, error) {
	return n.String(), nil
}

// MarshalJSON satisfies the encoding/json.Marshaler interface.
func (n Number) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte('"')
	write(&b, n)
	b.WriteByte('"')
	return b.Bytes(), nil
}

// MarshalText satisfies the encoding.TextMarshaler interface.
func (n Number) MarshalText() ([]byte, error) {
	var b bytes.Buffer
	write(&b, n)
	return b.Bytes(), nil
}

// UnmarshalYAML satisfies the gopkg.in/yaml.v3.Unmarshaler interface.
func (n *Number) UnmarshalYAML(value *yaml.Node) error {
	nn, err := ParseNumber(value.Value)
	if err != nil {
		return err
	}
	*n = nn
	return nil
}

// UnmarshalJSON satisfies the encoding/json.Unmarshaler interface.
func (n *Number) UnmarshalJSON(data []byte) error {
	if len(data) > 2 {
		return n.UnmarshalText(data[1 : len(data)-1])
	}
	return &Error{ErrorEmpty{}}
}

// UnmarshalText satisfies the encoding.TextUnmarshaler interface.
func (n *Number) UnmarshalText(text []byte) error {
	var l = len(text)

	if l == 0 {
		return errorEmpty
	}

	var M, m, p Number

	for i, partIndex := 0, 0; partIndex < 3 && i < len(text); i++ {
		if c := text[i]; c == '.' {
			partIndex++
		} else if c >= '0' && c <= '9' {
			switch partIndex {
			case 0:
				if M = M*10 + Number(c-'0'); M > 65535 {
					return &Error{ErrorMajorTooBig(text)}
				}
			case 1:
				if m = m*10 + Number(c-'0'); m > 255 {
					return &Error{ErrorMinorTooBig(text)}
				}
			case 2:
				if p = p*10 + Number(c-'0'); p > 255 {
					return &Error{ErrorPatchTooBig(text)}
				}
			}
		} else {
			return &Error{ErrorInvalidCharacter{string(text), c}}
		}
	}

	*n = M<<16 | m<<8 | p
	return nil
}

// ParseNumber takes a string, parses it and
// returns a Number if parsing was successful.
func ParseNumber(s string) (Number, error) {
	var l = len(s)

	if l == 0 {
		return 0, errorEmpty
	}

	var M, m, p Number

	for i, partIndex := 0, 0; partIndex < 3 && i < l; i++ {
		if c := s[i]; c == '.' {
			partIndex++
		} else if c >= '0' && c <= '9' {
			switch partIndex {
			case 0:
				if M = M*10 + Number(c-'0'); M > 65535 {
					return 0, &Error{ErrorMajorTooBig(s)}
				}
			case 1:
				if m = m*10 + Number(c-'0'); m > 255 {
					return 0, &Error{ErrorMinorTooBig(s)}
				}
			case 2:
				if p = p*10 + Number(c-'0'); p > 255 {
					return 0, &Error{ErrorPatchTooBig(s)}
				}
			}
		} else {
			return 0, &Error{ErrorInvalidCharacter{s, c}}
		}
	}

	return M<<16 | m<<8 | p, nil
}

// NewNumber creates a Number out of its major, minor and patch components.
func NewNumber(major uint16, minor, patch byte) Number {
	return Number(major)<<16 | Number(minor)<<8 | Number(patch)
}

const majorMask Number = 0b11111111111111110000000000000000
const minorMask Number = 0b00000000000000001111111100000000
const patchMask Number = 0b00000000000000000000000011111111
const invMajorMask Number = minorMask | patchMask
const invMinorMask Number = majorMask | patchMask
const invPatchMask Number = majorMask | minorMask

func write(b interface {
	WriteString(string) (n int, err error)
	WriteByte(byte) error
	Write([]byte) (n int, err error)
}, n Number) {
	p := n & patchMask

	b.WriteString(strconv.Itoa(int(n >> 16 & invMajorMask)))
	if m := n >> 8 & patchMask; m > 0 {
		b.WriteByte('.')
		b.WriteString(strconv.Itoa(int(m)))
		if p > 0 {
			b.WriteByte('.')
			b.WriteString(strconv.Itoa(int(p)))
		}
	} else if p > 0 {
		b.Write(zeroMinor)
		b.WriteString(strconv.Itoa(int(p)))
	}
}

var zeroMinor = []byte{'.', '0', '.'}
