/*
 *
 * Copyright 2024 tofuutils authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package versionfinder

import (
	"regexp"
	"strings"
)

const versionRegexpRaw string = `(v?[0-9]+(\.[0-9]+){0,2}(-[0-9A-Za-z\-.]+)?|alpha\-?[0-9]+)`

var (
	versionRegexp      = regexp.MustCompilePOSIX(versionRegexpRaw)             //nolint
	exactVersionRegexp = regexp.MustCompilePOSIX("^" + versionRegexpRaw + "$") //nolint
)

// Find returns a version without starting 'v'.
func Find(versionStr string) string {
	versionStr = versionRegexp.FindString(versionStr)
	if versionStr != "" && versionStr[0] == 'v' {
		return versionStr[1:]
	}

	return versionStr
}

func IsValid(versionStr string) bool {
	return exactVersionRegexp.MatchString(versionStr)
}

// Clean cleans the version string. IsValid(versionStr) must be true.
func Clean(versionStr string) string {
	if strings.HasPrefix(versionStr, "alpha") {
		return versionStr
	}

	before, after, found := strings.Cut(versionStr, "-")
	parts := strings.SplitN(before, ".", 3)
	major, minor, fixes := parts[0], "0", "0"
	if major[0] == 'v' {
		major = major[1:]
	}

	switch len(parts) {
	case 3:
		fixes = parts[2]

		fallthrough
	case 2:
		minor = parts[1]
	}

	var builder strings.Builder
	builder.WriteString(major)
	builder.WriteByte('.')
	builder.WriteString(minor)
	builder.WriteByte('.')
	builder.WriteString(fixes)
	if found {
		builder.WriteByte('-')
		builder.WriteString(after)
	}

	return builder.String()
}
