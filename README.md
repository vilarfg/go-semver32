# go-semver32

[![Sourcegraph](https://sourcegraph.com/github.com/vilarfg/go-semver32/-/badge.svg)](https://sourcegraph.com/github.com/vilarfg/go-semver32?badge)
[![Build Status](https://travis-ci.com/vilarfg/go-semver32.svg?branch=main)](https://travis-ci.com/vilarfg/go-semver32)
[![codecov](https://codecov.io/gh/vilarfg/go-semver32/branch/main/graph/badge.svg?token=3AFQQXD0QA)](https://codecov.io/gh/vilarfg/go-semver32)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/vilarfg/go-semver32)](https://pkg.go.dev/github.com/vilarfg/go-semver32)
[![rcard](https://goreportcard.com/badge/github.com/vilarfg/go-semver32?v=4)](https://goreportcard.com/report/github.com/vilarfg/go-semver32)
[![License](https://img.shields.io/github/license/vilarfg/go-semver32)](https://raw.githubusercontent.com/vilarfg/go-semver32/master/LICENSE)

Package semver offers a way to represent SemVer numbers in 32 bits.

These SemVer numbers are NOT spec compliant as defined on semver.org;
as they do not hold prerelease or build informantion and the maximum values
for the major, minor and patch components are 65,535, 255 and 255
respectively.

Use this package:

- if you don't need to store prerelease and/or build info
- if you are certain
  - the Major component will never exceed 65,535
  - the Minor and Patch components will never exceed 255

## Installation

```sh
go get -u github.com/vilarfg/go-semver32
```

## Usage

```go
package name

import "github.com/vilarfg/go-semver32"

func Func() {
    n, err := semver.ParseNumber("0.1.1")
    if err != nil {
        // handle error, though no error will be produced 
        // for that specific string.
    }
    fmt.Printf("%d is v%s", n, n) // => 257 is v0.1.1
}
```

## License

[MIT](https://github.com/vilarfg/go-semver32/blob/master/LICENSE). Copyright © 2020 Fernando G. Vilar
