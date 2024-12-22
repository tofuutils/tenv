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
	"errors"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-version"

	"github.com/tofuutils/tenv/v4/pkg/loghelper"
	"github.com/tofuutils/tenv/v4/versionmanager/lastuse"
)

const (
	allKey  = "all"
	butLast = "but-last"

	notUsedForPrefix   = "not-used-for:"
	notUsedSincePrefix = "not-used-since:"

	notUsedForPrefixLen   = len(notUsedForPrefix)
	notUsedSincePrefixLen = len(notUsedSincePrefix)
)

var errDurationParsing = errors.New("unrecognized duration format")

// versions must be sorted in descending order.
func SelectVersionsToUninstall(behaviourOrConstraint string, installPath string, versions []string, displayer loghelper.Displayer) ([]string, error) {
	switch {
	case behaviourOrConstraint == allKey:
		return versions, nil
	case behaviourOrConstraint == butLast:
		if len(versions) == 0 {
			return nil, nil
		}

		return versions[1:], nil // allowed by descending order
	case strings.HasPrefix(behaviourOrConstraint, notUsedForPrefix):
		forStr := behaviourOrConstraint[notUsedForPrefixLen:]

		var err error
		daysInt, monthsInt := 0, 0
		lastIndex := len(forStr) - 1
		switch forStr[lastIndex] {
		case 'd', 'D':
			daysInt, err = strconv.Atoi(forStr[:lastIndex])
			if err != nil {
				return nil, err
			}
		case 'm', 'M':
			monthsInt, err = strconv.Atoi(forStr[:lastIndex])
			if err != nil {
				return nil, err
			}
		default:
			return nil, errDurationParsing
		}

		beforeDate := time.Now().AddDate(0, -monthsInt, -daysInt)
		pred := predicateBeforeDate(installPath, beforeDate, displayer)

		return filterStrings(versions, pred), nil
	case strings.HasPrefix(behaviourOrConstraint, notUsedSincePrefix):
		dateStr := behaviourOrConstraint[notUsedSincePrefixLen:]

		beforeDate, err := time.Parse(time.DateOnly, dateStr)
		if err != nil {
			return nil, err
		}
		pred := predicateBeforeDate(installPath, beforeDate, displayer)

		return filterStrings(versions, pred), nil
	default:
		constraint, err := version.NewConstraint(behaviourOrConstraint)
		if err != nil {
			return nil, err
		}
		pred := predicateFromConstraint(constraint)

		return filterStrings(versions, pred), nil
	}
}

func filterStrings(stringSlice []string, pred func(string) bool) []string {
	selected := make([]string, 0, len(stringSlice))
	for _, str := range stringSlice {
		if pred(str) {
			selected = append(selected, str)
		}
	}

	return selected
}

func predicateBeforeDate(installPath string, beforeDate time.Time, displayer loghelper.Displayer) func(string) bool {
	return func(versionStr string) bool {
		useDate := lastuse.Read(filepath.Join(installPath, versionStr), displayer)

		return useDate.Before(beforeDate)
	}
}
