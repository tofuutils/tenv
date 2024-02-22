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
	"github.com/hashicorp/go-version"
	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/pkg/loghelper"
	terragruntparser "github.com/tofuutils/tenv/versionmanager/semantic/parser/terragrunt"
	tfparser "github.com/tofuutils/tenv/versionmanager/semantic/parser/tf"
	"github.com/tofuutils/tenv/versionmanager/semantic/parser/types"
)

const (
	LatestAllowedKey = "latest-allowed"
	LatestPreKey     = "latest-pre"
	LatestStableKey  = "latest-stable"
	LatestKey        = "latest"
	MinRequiredKey   = "min-required"
)

var (
	TfPredicateReaders = []types.PredicateReader{readTfVersionFromTerragruntFile, readTfFiles} //nolint
	TgPredicateReaders = []types.PredicateReader{readTgVersionFromTerragruntFile}              //nolint
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

func ParsePredicate(behaviourOrConstraint types.DetectionInfo, displayName string, predicateReaders []types.PredicateReader, conf *config.Config) (types.PredicateInfo, error) {
	reverseOrder := true
	recordeds := behaviourOrConstraint.Messages
	switch behaviourOrConstraint.Version {
	case MinRequiredKey:
		reverseOrder = false // start with older

		fallthrough // same predicate retrieving
	case LatestAllowedKey:
		for _, reader := range predicateReaders {
			predicate, msgs, err := reader(conf)
			recordeds = append(recordeds, msgs...)
			if err != nil {
				return types.PredicateInfo{Messages: recordeds}, err
			}
			if predicate != nil {
				return types.PredicateInfo{Predicate: predicate, ReverseOrder: reverseOrder, Messages: recordeds}, nil
			}
		}

		recordeds = append(recordeds, loghelper.RecordedMessage{Message: loghelper.Concat("No ", displayName, " version requirement found in project files, fallback to ", LatestKey, " strategy")})

		fallthrough // fallback to latest
	case LatestKey, LatestStableKey:
		return types.PredicateInfo{Predicate: StableVersion, ReverseOrder: true, Messages: recordeds}, nil
	case LatestPreKey:
		return types.PredicateInfo{Predicate: alwaysTrue, ReverseOrder: true, Messages: recordeds}, nil
	default:
		constraint, err := version.NewConstraint(behaviourOrConstraint.Version)
		if err != nil {
			return types.PredicateInfo{Messages: recordeds}, err
		}

		return types.PredicateInfo{Predicate: predicateFromConstraint(constraint), ReverseOrder: true, Messages: recordeds}, nil
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

func readPredicate(constraintRetriever func(*config.Config) (string, []loghelper.RecordedMessage, error), conf *config.Config) (func(string) bool, []loghelper.RecordedMessage, error) {
	constraintStr, recordeds, err := constraintRetriever(conf)
	if err != nil || constraintStr == "" {
		return nil, recordeds, err
	}

	constraint, err := version.NewConstraint(constraintStr)
	if err != nil {
		return nil, recordeds, err
	}

	return predicateFromConstraint(constraint), recordeds, nil
}

func readTfFiles(conf *config.Config) (func(string) bool, []loghelper.RecordedMessage, error) {
	requireds, recordeds, err := tfparser.GatherRequiredVersion(conf)
	if err != nil {
		return nil, recordeds, err
	}

	var constraint version.Constraints
	for _, required := range requireds {
		temp, err := version.NewConstraint(required)
		if err != nil {
			return nil, recordeds, err
		}
		constraint = append(constraint, temp...)
	}
	if len(constraint) == 0 {
		return nil, recordeds, nil
	}

	return predicateFromConstraint(constraint), recordeds, nil
}

func readTfVersionFromTerragruntFile(conf *config.Config) (func(string) bool, []loghelper.RecordedMessage, error) {
	return readPredicate(terragruntparser.RetrieveTerraformVersionConstraint, conf)
}

func readTgVersionFromTerragruntFile(conf *config.Config) (func(string) bool, []loghelper.RecordedMessage, error) {
	return readPredicate(terragruntparser.RetrieveTerraguntVersionConstraint, conf)
}
