package gvm

import (
	"fmt"
	"sort"

	version "github.com/hashicorp/go-version"
)

type GoVersion struct {
	in      string
	version *version.Version
}

func ParseVersion(in string) (*GoVersion, error) {
	var v *version.Version

	if in != "tip" {
		var err error
		v, err = version.NewVersion(in)
		if err != nil {
			return nil, err
		}
	}

	return &GoVersion{in: in, version: v}, nil
}

func (v *GoVersion) String() string {
	if v.in == "tip" {
		return v.in
	}

	seg := v.version.Segments()
	if v.version.Prerelease() != "" {
		return fmt.Sprintf("%v.%v%v", seg[0], seg[1], v.version.Prerelease())
	}

	if len(seg) > 2 && seg[2] == 0 {
		return fmt.Sprintf("%v.%v", seg[0], seg[1])
	}
	return v.version.String()
}

func (v *GoVersion) LessThan(v2 *GoVersion) bool {
	if v.in == "tip" {
		return false
	}
	if v2.in == "tip" {
		return true
	}
	return v.version.LessThan(v2.version)
}

func (v *GoVersion) VendorSupport() (has, experimental bool) {
	if v.in == "tip" {
		return true, false
	}

	seg := v.version.Segments()
	if len(seg) < 2 {
		return false, false
	}

	return seg[1] >= 5, seg[1] == 5
}

func sortVersions(versions []*GoVersion) {
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].LessThan(versions[j])
	})
}
