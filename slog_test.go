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

package stacktrace_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/stacktrace"
)

// TestLogValue verifies that stacktrace errors implement slog.LogValuer and
// produce a structured group containing at least a "message" key and a
// "frames" key when logged.
func TestLogValue(t *testing.T) {
	err := stacktrace.NewError("something went wrong")

	// Cast to slog.LogValuer — if the interface is not implemented this will panic.
	lv, ok := err.(slog.LogValuer)
	require.True(t, ok, "stacktrace error should implement slog.LogValuer")

	val := lv.LogValue()
	assert.Equal(t, slog.KindGroup, val.Kind(), "LogValue should return a group value")

	attrs := val.Group()
	keys := make([]string, 0, len(attrs))
	for _, a := range attrs {
		keys = append(keys, a.Key)
	}
	assert.Contains(t, keys, "message", "group should contain a 'message' key")
	assert.Contains(t, keys, "frames", "group should contain a 'frames' key")

	// Verify the message value equals the error string.
	for _, a := range attrs {
		if a.Key == "message" {
			assert.Equal(t, err.Error(), a.Value.String())
		}
	}
}

// TestLogValuePropagated verifies that a propagated (wrapped) stacktrace error
// also implements slog.LogValuer and that all frames are captured.
func TestLogValuePropagated(t *testing.T) {
	inner := stacktrace.NewError("inner error")
	outer := stacktrace.Propagate(inner, "outer context")

	lv, ok := outer.(slog.LogValuer)
	require.True(t, ok, "propagated stacktrace error should implement slog.LogValuer")

	val := lv.LogValue()
	assert.Equal(t, slog.KindGroup, val.Kind())

	// The message should reflect the full chain.
	attrs := val.Group()
	for _, a := range attrs {
		if a.Key == "message" {
			msg := a.Value.String()
			assert.Contains(t, msg, "outer context", "message should include outer context")
			assert.Contains(t, msg, "inner error", "message should include inner error text")
		}
	}
}

// TestLogValueInSlogHandler verifies that a stacktrace error is emitted as
// structured fields (not a flat string) when passed to a slog.Logger with a
// JSON handler.
func TestLogValueInSlogHandler(t *testing.T) {
	err := stacktrace.NewError("disk full")

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	logger.Error("operation failed", "error", err)

	output := buf.String()
	// The JSON output should contain the "message" key inside the error group,
	// not just a flat string value for the error key.
	assert.True(t,
		strings.Contains(output, `"message"`) && strings.Contains(output, "disk full"),
		"JSON log output should contain structured 'message' field; got: %s", output,
	)
}
