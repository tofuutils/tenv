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
	"fmt"
	"slices"

	"github.com/hashicorp/go-version"
	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/pkg/apimsg"
	terragruntparser "github.com/tofuutils/tenv/versionmanager/semantic/parser/terragrunt"
	tfparser "github.com/tofuutils/tenv/versionmanager/semantic/parser/tf"
	tgswitchparser "github.com/tofuutils/tenv/versionmanager/semantic/parser/tgswitch"
)

const (
	LatestAllowedKey = "latest-allowed"
	LatestStableKey  = "latest-stable"
	LatestKey        = "latest"
	MinRequiredKey   = "min-required"
)

var (
	TfPredicateReaders = []func(*config.Config) (func(string) bool, bool, error){readTfVersionFromTerragruntFile, readTfFiles}                   //nolint
	TgPredicateReaders = []func(*config.Config) (func(string) bool, bool, error){readTgVersionFromTgSwitchFile, readTgVersionFromTerragruntFile} //nolint
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

func LatestVersionFromList(versions []string) (string, error) {
	versionLen := len(versions)
	if versionLen == 0 {
		return "", apimsg.ErrReturn
	}

	slices.SortFunc(versions, CmpVersion)

	return versions[versionLen-1], nil
}

// the boolean returned as second value indicates to reverse order for filtering.
func ParsePredicate(requestedVersion string, displayName string, predicateReaders []func(*config.Config) (func(string) bool, bool, error), conf *config.Config) (func(string) bool, bool, error) {
	reverseOrder := true
	switch requestedVersion {
	case MinRequiredKey:
		reverseOrder = false // start with older

		fallthrough // same predicate retrieving
	case LatestAllowedKey:
		for _, reader := range predicateReaders {
			predicate, found, err := reader(conf)
			if err != nil {
				return nil, false, err
			}
			if found {
				return predicate, reverseOrder, nil
			}
		}

		if conf.Verbose {
			fmt.Println("No", displayName, "version requirement found in files, fallback to latest-stable") //nolint
		}

		return StableVersion, true, nil // erase min-required case
	case LatestStableKey:
		return StableVersion, true, nil
	case LatestKey:
		return alwaysTrue, true, nil
	default:
		constraint, err := version.NewConstraint(requestedVersion)
		if err != nil {
			return nil, false, err
		}

		return predicateFromConstraint(constraint), true, nil
	}
}

func StableVersion(versionStr string) bool {
	v, err := version.NewVersion(versionStr)

	return err == nil && v.Prerelease() == ""
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

// the boolean returned as second value indicates if a predicate was found.
func readPredicate(constraintRetriever func(*config.Config) (string, error), conf *config.Config) (func(string) bool, bool, error) {
	constraintStr, err := constraintRetriever(conf)
	if err != nil || constraintStr == "" {
		return nil, false, err
	}

	constraint, err := version.NewConstraint(constraintStr)
	if err != nil {
		return nil, false, err
	}

	return predicateFromConstraint(constraint), true, nil
}

// the boolean returned as second value indicates if a predicate was found.
func readTfFiles(conf *config.Config) (func(string) bool, bool, error) {
	requireds, err := tfparser.GatherRequiredVersion(conf.Verbose)
	if err != nil {
		return nil, false, err
	}

	var constraint version.Constraints
	for _, required := range requireds {
		temp, err := version.NewConstraint(required)
		if err != nil {
			return nil, false, err
		}
		constraint = append(constraint, temp...)
	}
	if len(constraint) == 0 {
		return nil, false, nil
	}

	return predicateFromConstraint(constraint), true, nil
}

// the boolean returned as second value indicates if a predicate was found.
func readTfVersionFromTerragruntFile(conf *config.Config) (func(string) bool, bool, error) {
	return readPredicate(terragruntparser.RetrieveTerraformVersionConstraint, conf)
}

// the boolean returned as second value indicates if a predicate was found.
func readTgVersionFromTerragruntFile(conf *config.Config) (func(string) bool, bool, error) {
	return readPredicate(terragruntparser.RetrieveTerraguntVersionConstraint, conf)
}

// the boolean returned as second value indicates if a predicate was found.
func readTgVersionFromTgSwitchFile(conf *config.Config) (func(string) bool, bool, error) {
	return readPredicate(tgswitchparser.RetrieveTerraguntVersion, conf)
}
