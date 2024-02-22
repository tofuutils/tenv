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

type Displayer interface {
	Display(msg string)
	IsDebug() bool
	Log(level hclog.Level, msg string, args ...any)
	Flush(logMode bool)
}

type basicDisplayer struct {
	display func(string)
	logger  hclog.Logger
}

func NewBasicDisplayer(logger hclog.Logger, display func(string)) Displayer {
	return basicDisplayer{display: display, logger: logger}
}

func (bd basicDisplayer) Display(msg string) {
	bd.display(msg)
}

func (bd basicDisplayer) IsDebug() bool {
	return bd.logger.IsDebug()
}

func (bd basicDisplayer) Log(level hclog.Level, msg string, args ...any) {
	bd.logger.Log(level, msg, args...)
}

func (bd basicDisplayer) Flush(bool) {
}

type logWrapper struct {
	Displayer
}

func (lw logWrapper) Display(msg string) {
	lw.Displayer.Log(hclog.Debug, msg)
}

type recordedMessage struct {
	Level   hclog.Level
	Message string
	Args    []any
}

type recordingWrapper struct {
	Displayer
	recordeds []recordedMessage
}

func (rw *recordingWrapper) Display(msg string) {
	rw.recordeds = append(rw.recordeds, recordedMessage{Message: msg})
}

func (rw *recordingWrapper) Log(level hclog.Level, msg string, args ...any) {
	rw.recordeds = append(rw.recordeds, recordedMessage{Level: level, Message: msg, Args: args})
}

func (rw *recordingWrapper) Flush(logMode bool) {
	if logMode {
		rw.Displayer = logWrapper{Displayer: rw.Displayer}
	}
	for _, recorded := range rw.recordeds {
		if recorded.Level == hclog.NoLevel {
			rw.Displayer.Display(recorded.Message)
		} else {
			rw.Displayer.Log(recorded.Level, recorded.Message, recorded.Args...)
		}
	}
}

type stateWrapper struct {
	Displayer
}

func NewRecordingDisplayer(displayer Displayer) Displayer {
	return &stateWrapper{Displayer: &recordingWrapper{Displayer: displayer}}
}

func (sw *stateWrapper) Flush(logMode bool) {
	sw.Displayer.Flush(logMode)
	if rw, ok := sw.Displayer.(*recordingWrapper); ok {
		sw.Displayer = rw.Displayer // following call will be direct
	}
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
