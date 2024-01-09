/*
 *
 * Copyright 2024 gotofuenv authors.
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

package tofuversion

import "github.com/Masterminds/semver/v3"

func cmpVersion(a string, b string) int {
	v1, err1 := semver.NewVersion(a)
	v2, err2 := semver.NewVersion(b)

	if hasErr1, hasErr2 := err1 != nil, err2 != nil; hasErr1 {
		if hasErr2 {
			return 0
		}
		return -1
	} else if hasErr2 {
		return 1
	}

	return v1.Compare(v2)
}

func parsePredicate(requestedVersion string) (func(string) bool, bool, error) {
	predicate := alwaysTrue
	reverseOrder := true
	switch requestedVersion {
	case "min-required":
		reverseOrder = false // start with older
		fallthrough          // same predicate retrieving
	case "latest-allowed":
		// TODO predicate from HCL parsing
	case "latest":
		// nothing todo (alwaysTrue and reverseOrder will work)
	default:
		constraint, err := semver.NewConstraint(requestedVersion)
		if err != nil {
			return nil, false, err
		}

		predicate = func(version string) bool {
			v, err := semver.NewVersion(version)
			if err != nil {
				return false
			}

			ok, _ := constraint.Validate(v)
			return ok
		}
	}
	return predicate, reverseOrder, nil
}

func alwaysTrue(string) bool {
	return true
}
