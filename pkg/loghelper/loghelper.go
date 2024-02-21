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

package loghelper

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/hashicorp/go-hclog"
)

const Error = "error"

func BuildDisplayFunc(writer io.Writer, usedColor *color.Color) func(...any) {
	return func(a ...any) {
		usedColor.Fprintln(writer, a...)
	}
}

func LevelWarnOrDebug(debug bool) hclog.Level {
	if debug {
		return hclog.Debug
	}

	return hclog.Warn
}

func MultiDisplay(display func(...any)) func([]string) {
	return func(msgs []string) {
		for _, msg := range msgs {
			display(msg)
		}
	}
}

func MultiLogDebug(logger hclog.Logger) func([]string) {
	return func(msgs []string) {
		for _, msg := range msgs {
			logger.Debug(msg)
		}
	}
}

func NoDisplay(...any) {}

func StdDisplay(a ...any) {
	fmt.Println(a...) //nolint
}
