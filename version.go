package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Version section struct
// The version will follow semver.
type Version struct {
	major   int
	minor   int
	patch   int
	message string
}

// Export as string
func (v Version) String() string {
	m := "\n"
	if v.message != "" {
		m = v.message
	}
	return fmt.Sprintf("[%d.%d.%d]\n%s", v.major, v.minor, v.patch, m)
}

// Format version string like v1.0.0
func (v Version) VersionString() string {
	return fmt.Sprintf("v%d.%d.%d", v.major, v.minor, v.patch)
}

// Bump patch version and create new version
func (v Version) BumpPatch(msg string) Version {
	return Version{
		major:   v.major,
		minor:   v.minor,
		patch:   v.patch + 1,
		message: msg,
	}
}

// Bump minor version and create new version
func (v Version) BumpMinor(msg string) Version {
	return Version{
		major:   v.major,
		minor:   v.minor + 1,
		message: msg,
	}
}

// Bump major version and create new version
func (v Version) BumpMajor(msg string) Version {
	return Version{
		major:   v.major + 1,
		message: msg,
	}
}

// Create version struct from parsed string
func createVersion(major, minor, patch, message string) Version {
	_major, _ := strconv.Atoi(major)
	_minor, _ := strconv.Atoi(minor)
	_patch, _ := strconv.Atoi(patch)
	return Version{
		major:   _major,
		minor:   _minor,
		patch:   _patch,
		message: strings.Trim(message, "\r\n"),
	}
}

// Version slice type to define sort interface
type Versions []Version

func (v Versions) Len() int {
	return len(v)
}
func (v Versions) Less(i, j int) bool {
	if v[i].major < v[j].major {
		return false
	} else if v[i].major > v[j].major {
		return true
	}

	if v[i].minor < v[j].minor {
		return false
	} else if v[i].minor > v[j].minor {
		return true
	}

	if v[i].patch < v[j].patch {
		return false
	} else if v[i].patch > v[j].patch {
		return true
	}

	return false
}
func (v Versions) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
