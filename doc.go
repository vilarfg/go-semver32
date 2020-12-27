// Copyright (c) 2020 Fernando G. Vilar
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package semver offers a way to represent SemVer numbers in 32 bits.
//
// These SemVer numbers are NOT spec compliant as defined on semver.org;
// as they do not hold prerelease or build informantion and the maximum values
// for the major, minor and patch components are 65,535, 255 and 255
// respectively.
package semver
