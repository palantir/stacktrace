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
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/palantir/stacktrace"
)

func TestMessage(t *testing.T) {
	err := startDoing()
	err = PublicObj{}.DoPublic(err)
	err = PublicObj{}.doPrivate(err)
	err = privateObj{}.DoPublic(err)
	err = privateObj{}.doPrivate(err)
	err = (&ptrObj{}).doPtr(err)
	err = doClosure(err)

	expected := strings.Join([]string{
		"so closed",
		" --- at github.com/palantir/stacktrace/functions_for_test.go:51 (doClosure.func1) ---",
		"Caused by: pointedly",
		" --- at github.com/palantir/stacktrace/functions_for_test.go:46 (ptrObj.doPtr) ---",
		" --- at github.com/palantir/stacktrace/functions_for_test.go:42 (privateObj.doPrivate) ---",
		" --- at github.com/palantir/stacktrace/functions_for_test.go:38 (privateObj.DoPublic) ---",
		" --- at github.com/palantir/stacktrace/functions_for_test.go:34 (PublicObj.doPrivate) ---",
		" --- at github.com/palantir/stacktrace/functions_for_test.go:30 (PublicObj.DoPublic) ---",
		"Caused by: failed to start doing",
		" --- at github.com/palantir/stacktrace/functions_for_test.go:26 (startDoing) ---",
	}, "\n")
	stacktrace.DefaultFormat = stacktrace.FormatFull
	assert.Equal(t, expected, err.Error())
	assert.Equal(t, expected, fmt.Sprint(err))
}

func TestGetCode(t *testing.T) {
	for _, test := range []struct {
		originalError error
		originalCode  stacktrace.ErrorCode
	}{
		{
			originalError: errors.New("err"),
			originalCode:  stacktrace.NoCode,
		},
		{
			originalError: stacktrace.NewError("err"),
			originalCode:  stacktrace.NoCode,
		},
		{
			originalError: stacktrace.NewErrorWithCode(EcodeInvalidVillain, "err"),
			originalCode:  EcodeInvalidVillain,
		},
		{
			originalError: stacktrace.NewMessageWithCode(EcodeNoSuchPseudo, "err"),
			originalCode:  EcodeNoSuchPseudo,
		},
	} {
		err := test.originalError
		assert.Equal(t, test.originalCode, stacktrace.GetCode(err))

		err = stacktrace.Propagate(err, "")
		assert.Equal(t, test.originalCode, stacktrace.GetCode(err))

		err = stacktrace.PropagateWithCode(err, EcodeNotFastEnough, "")
		assert.Equal(t, EcodeNotFastEnough, stacktrace.GetCode(err))

		err = stacktrace.PropagateWithCode(err, EcodeTimeIsIllusion, "")
		assert.Equal(t, EcodeTimeIsIllusion, stacktrace.GetCode(err))
	}
}

func TestPropagateNil(t *testing.T) {
	var err error

	err = stacktrace.Propagate(err, "")
	assert.Nil(t, err)

	err = stacktrace.PropagateWithCode(err, EcodeNotImplemented, "")
	assert.Nil(t, err)

	assert.Equal(t, stacktrace.NoCode, stacktrace.GetCode(err))
}
