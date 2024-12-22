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

package semantic

import (
	"regexp"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-version"

	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
	iacparser "github.com/tofuutils/tenv/v4/versionmanager/semantic/parser/iac"
	"github.com/tofuutils/tenv/v4/versionmanager/semantic/types"
)

const (
	LatestAllowedKey = "latest-allowed"
	LatestPreKey     = "latest-pre"
	LatestStableKey  = "latest-stable"
	LatestKey        = "latest"
	MinRequiredKey   = "min-required"

	LatestPrefix = "latest:"
	MinPrefix    = "min:"
)

func CmpVersion(v1Str string, v2Str string) int {
	v1, err1 := version.NewVersion(v1Str) //nolint
	v2, err2 := version.NewVersion(v2Str) //nolint

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

func ParsePredicate(behaviourOrConstraint string, displayName string, constraintInfo types.ConstraintInfo, iacExts []iacparser.ExtDescription, conf *config.Config) (types.PredicateInfo, error) {
	reverseOrder := true
	switch {
	case behaviourOrConstraint == MinRequiredKey:
		reverseOrder = false // start with older

		fallthrough // same predicate retrieving
	case behaviourOrConstraint == LatestAllowedKey:
		constraints, err := readIACfiles(constraintInfo, iacExts, conf)
		if err != nil {
			return types.PredicateInfo{}, err
		}
		if len(constraints) != 0 {
			return types.PredicateInfo{Predicate: predicateFromConstraint(constraints), ReverseOrder: reverseOrder}, nil
		}

		conf.Displayer.Display(loghelper.Concat("No ", displayName, " version requirement found in project files, fallback to ", LatestKey, " strategy"))

		fallthrough // fallback to latest
	case behaviourOrConstraint == LatestKey, behaviourOrConstraint == LatestStableKey:
		return types.PredicateInfo{Predicate: StableVersion, ReverseOrder: true}, nil
	case behaviourOrConstraint == LatestPreKey:
		return types.PredicateInfo{Predicate: alwaysTrue, ReverseOrder: true}, nil
	case strings.HasPrefix(behaviourOrConstraint, MinPrefix):
		reverseOrder = false // start with older

		fallthrough // same behaviour
	case strings.HasPrefix(behaviourOrConstraint, LatestPrefix):
		conf.Displayer.Display("Use of regexp is discouraged, try version constraint instead")

		re, err := regexp.Compile(behaviourOrConstraint[strings.Index(behaviourOrConstraint, ":")+1:])
		if err != nil {
			return types.PredicateInfo{}, err
		}

		return types.PredicateInfo{Predicate: re.MatchString, ReverseOrder: reverseOrder}, nil
	default:
		constraint, err := addDefaultConstraint(constraintInfo, conf, behaviourOrConstraint)
		if err != nil {
			return types.PredicateInfo{}, err
		}

		return types.PredicateInfo{Predicate: predicateFromConstraint(constraint), ReverseOrder: true}, nil
	}
}

func StableVersion(versionStr string) bool {
	v, err := version.NewVersion(versionStr)

	return err == nil && v.Prerelease() == ""
}

func addDefaultConstraint(constraintInfo types.ConstraintInfo, conf *config.Config, requireds ...string) (version.Constraints, error) {
	if defaultConstraint := constraintInfo.ReadDefaultConstraint(); defaultConstraint != "" {
		requireds = append(requireds, defaultConstraint)
	}
	conf.Displayer.Log(hclog.Debug, "Find", "constraints", requireds)

	var constraint version.Constraints
	for _, required := range requireds {
		temp, err := version.NewConstraint(required)
		if err != nil {
			return nil, err
		}
		constraint = append(constraint, temp...)
	}

	return constraint, nil
}

func alwaysTrue(string) bool {
	return true
}

func predicateFromConstraint(constraint version.Constraints) func(string) bool {
	return func(versionStr string) bool {
		v, err := version.NewVersion(versionStr)

		return err == nil && constraint.Check(v)
	}
}

func readIACfiles(constraintInfo types.ConstraintInfo, iacExts []iacparser.ExtDescription, conf *config.Config) (version.Constraints, error) {
	requireds, err := iacparser.GatherRequiredVersion(conf, iacExts)
	if err != nil {
		return nil, err
	}

	return addDefaultConstraint(constraintInfo, conf, requireds...)
}
