// Package version handles numeric major.minor.patch[.build][-suffix] versions.
// This is the project's numeric build convention, not SemVer prerelease ordering.
package version

import (
	"strconv"
	"strings"
)

type Version string

func (ver Version) String() string      { return string(ver) }
func (ver Version) TrimPrefix() Version { return Version(strings.TrimPrefix(string(ver), "v")) }
func (ver Version) TrimSuffix() Version { s, _, _ := strings.Cut(string(ver), "-"); return Version(s) }

func (ver Version) parts() []string {
	parts := strings.Split(ver.TrimPrefix().TrimSuffix().String(), ".")
	if len(parts) > 4 {
		parts = parts[:4]
	}
	for len(parts) < 3 {
		parts = append(parts, "0")
	}
	return parts
}
func (ver Version) PkgSuffix() string {
	parts := ver.parts()
	build := "0"
	if len(parts) == 4 {
		build = parts[3]
	}
	return strings.Join(parts[:3], ".") + "-" + build
}
func (ver Version) Tag() string {
	parts := ver.parts()
	tag := strings.Join(parts[:3], ".")
	if len(parts) == 4 {
		tag += "-" + parts[3]
	}
	return tag
}

// CompatibleVersion keeps the patch and optional build, padding short versions.
func (ver Version) CompatibleVersion() Version {
	prefix := ""
	if strings.HasPrefix(string(ver), "v") {
		prefix = "v"
	}
	return Version(prefix + strings.Join(ver.parts(), "."))
}

func (ver Version) numbers() ([4]int64, bool) {
	var result [4]int64
	parts := strings.Split(ver.TrimPrefix().TrimSuffix().String(), ".")
	if len(parts) > 4 {
		return result, false
	}
	for i, part := range parts {
		value, err := strconv.ParseInt(part, 10, 64)
		// The numeric encoding reserves three digits per release field and six for build.
		limit := int64(1000)
		if i == 3 {
			limit = 1000000
		}
		if err != nil || value < 0 || value >= limit {
			return result, false
		}
		result[i] = value
	}
	return result, true
}

// Number returns the numeric encoding, or zero for invalid/out-of-range versions.
func (ver Version) Number() int64 {
	n, ok := ver.numbers()
	if !ok {
		return 0
	}
	return n[0]*1000000000000 + n[1]*1000000000 + n[2]*1000000 + n[3]
}

// Compatibility requires identical release fields and a server build >= client build.
func (ver Version) Compatibility(other Version) bool {
	a, aok := ver.numbers()
	b, bok := other.numbers()
	return aok && bok && a[0] == b[0] && a[1] == b[1] && a[2] == b[2] && a[3] >= b[3]
}
func (ver Version) Equal(other Version) bool       { return ver == other }
func (ver Version) IsLowerThan(other Version) bool { return ver.Number() < other.Number() }
