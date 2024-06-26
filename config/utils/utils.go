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

package configutils

import (
	"os"
	"strconv"
)

func GetenvBool(defaultValue bool, key string) (bool, error) {
	if valueStr := os.Getenv(key); valueStr != "" {
		return strconv.ParseBool(valueStr)
	}

	return defaultValue, nil
}

func GetenvBoolFallback(defaultValue bool, keys ...string) (bool, error) {
	if valueStr := GetenvFallback(keys...); valueStr != "" {
		return strconv.ParseBool(valueStr)
	}

	return defaultValue, nil
}

func GetenvFallback(keys ...string) string {
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}

	return ""
}
