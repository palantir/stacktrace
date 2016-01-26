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
	"github.com/palantir/stacktrace"
)

type PublicObj struct{}
type privateObj struct{}
type ptrObj struct{}

func startDoing() error {
	return stacktrace.NewError("%s %s %s %s", "failed", "to", "start", "doing")
}

func (PublicObj) DoPublic(err error) error {
	return stacktrace.Propagate(err, "")
}

func (PublicObj) doPrivate(err error) error {
	return stacktrace.Propagate(err, "")
}

func (privateObj) DoPublic(err error) error {
	return stacktrace.Propagate(err, "")
}

func (privateObj) doPrivate(err error) error {
	return stacktrace.Propagate(err, "")
}

func (*ptrObj) doPtr(err error) error {
	return stacktrace.Propagate(err, "pointedly")
}

func doClosure(err error) error {
	return func() error {
		return stacktrace.Propagate(err, "so closed")
	}()
}
