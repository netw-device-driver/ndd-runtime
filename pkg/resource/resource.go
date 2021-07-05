/*
Copyright 2020 Wim Henderickx.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resource

import (
	"context"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// A ClientApplicator may be used to build a single 'client' that satisfies both
// client.Client and Applicator.
type ClientApplicator struct {
	client.Client
	Applicator
}

// An ApplyFn is a function that satisfies the Applicator interface.
type ApplyFn func(context.Context, client.Object, ...ApplyOption) error

// Apply changes to the supplied object.
func (fn ApplyFn) Apply(ctx context.Context, o client.Object, ao ...ApplyOption) error {
	return fn(ctx, o, ao...)
}

// An Applicator applies changes to an object.
type Applicator interface {
	Apply(context.Context, client.Object, ...ApplyOption) error
}

// An ApplyOption is called before patching the current object to match the
// desired object. ApplyOptions are not called if no current object exists.
type ApplyOption func(ctx context.Context, current, desired runtime.Object) error

// UpdateFn returns an ApplyOption that is used to modify the current object to
// match fields of the desired.
func UpdateFn(fn func(current, desired runtime.Object)) ApplyOption {
	return func(_ context.Context, c, d runtime.Object) error {
		fn(c, d)
		return nil
	}
}

type errNotControllable struct{ error }

func (e errNotControllable) NotControllable() bool {
	return true
}

// MustBeControllableBy requires that the current object is controllable by an
// object with the supplied UID. An object is controllable if its controller
// reference matches the supplied UID, or it has no controller reference. An
// error that satisfies IsNotControllable will be returned if the current object
// cannot be controlled by the supplied UID.
func MustBeControllableBy(u types.UID) ApplyOption {
	return func(_ context.Context, current, _ runtime.Object) error {
		c := metav1.GetControllerOf(current.(metav1.Object))
		if c == nil {
			return nil
		}

		if c.UID != u {
			return errNotControllable{errors.Errorf("existing object is not controlled by UID %q", u)}

		}
		return nil
	}
}

type errNotAllowed struct{ error }

func (e errNotAllowed) NotAllowed() bool {
	return true
}

// IsNotAllowed returns true if the supplied error indicates that an operation
// was not allowed.
func IsNotAllowed(err error) bool {
	_, ok := err.(interface {
		NotAllowed() bool
	})
	return ok
}

// AllowUpdateIf will only update the current object if the supplied fn returns
// true. An error that satisfies IsNotAllowed will be returned if the supplied
// function returns false. Creation of a desired object that does not currently
// exist is always allowed.
func AllowUpdateIf(fn func(current, desired runtime.Object) bool) ApplyOption {
	return func(_ context.Context, current, desired runtime.Object) error {
		if fn(current, desired) {
			return nil
		}
		return errNotAllowed{errors.New("update not allowed")}
	}
}
