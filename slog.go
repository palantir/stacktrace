// Copyright 2016 Palantir Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stacktrace

import (
	"fmt"
	"log/slog"
)

// LogValue implements the slog.LogValuer interface so that stacktrace errors
// appear as structured groups in structured log output rather than as a flat
// string.
//
// The returned group contains:
//   - "message": the full human-readable error string (same as err.Error())
//   - "frames": a list of slog.Attr values, one per stack frame in the chain,
//     each being a group with "function", "file", and "line" keys.
//
// Example slog output (JSON handler):
//
//	{"error":{"message":"outer: inner","frames":[{"function":"pkg.Fn","file":"pkg/fn.go","line":42}]}}
func (st *stacktrace) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("message", st.Error()),
		slog.Any("frames", collectFrames(st)),
	)
}

// collectFrames walks the stacktrace chain and returns a slice of slog.Value,
// one per frame that has file/line information.
func collectFrames(st *stacktrace) []slog.Value {
	var frames []slog.Value
	for curr, ok := st, true; ok; curr, ok = curr.cause.(*stacktrace) {
		if curr.file == "" {
			continue
		}
		attrs := []slog.Attr{
			slog.String("file", fmt.Sprintf("%s:%d", curr.file, curr.line)),
		}
		if curr.function != "" {
			attrs = append([]slog.Attr{slog.String("function", curr.function)}, attrs...)
		}
		if curr.message != "" {
			attrs = append([]slog.Attr{slog.String("message", curr.message)}, attrs...)
		}
		frames = append(frames, slog.GroupValue(attrs...))
	}
	return frames
}

// Ensure *stacktrace implements slog.LogValuer at compile time.
var _ slog.LogValuer = (*stacktrace)(nil)
