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
	"strings"

	"github.com/fatih/color"
	"github.com/hashicorp/go-hclog"
)

const Error = "error"

type RecordedMessage struct {
	Level   hclog.Level
	Message string
	Args    []any
}

func BuildDisplayFunc(writer io.Writer, usedColor *color.Color) func(string) {
	return func(msg string) {
		usedColor.Fprintln(writer, msg)
	}
}

func LevelWarnOrDebug(debug bool) hclog.Level {
	if debug {
		return hclog.Debug
	}

	return hclog.Warn
}

func MultiDisplay(logger hclog.Logger, display func(string)) func([]RecordedMessage) {
	return func(recordeds []RecordedMessage) {
		for _, recorded := range recordeds {
			if recorded.Level == hclog.NoLevel {
				display(recorded.Message)
			} else {
				logger.Log(recorded.Level, recorded.Message, recorded.Args...)
			}
		}
	}
}

func MultiLog(logger hclog.Logger) func([]RecordedMessage) {
	return func(recordeds []RecordedMessage) {
		for _, recorded := range recordeds {
			if recorded.Level == hclog.NoLevel {
				logger.Debug(recorded.Message, recorded.Args...)
			} else {
				logger.Log(recorded.Level, recorded.Message, recorded.Args...)
			}
		}
	}
}

func NoDisplay(string) {}

func StdDisplay(msg string) {
	fmt.Println(msg) //nolint
}

func Concat(parts ...string) string {
	var builder strings.Builder
	for _, part := range parts {
		builder.WriteString(part)
	}

	return builder.String()
}
