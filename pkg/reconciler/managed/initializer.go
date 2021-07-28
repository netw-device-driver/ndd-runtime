/*
Copyright 2021 Wim Henderickx.

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

package managed

import (
	"context"

	"github.com/netw-device-driver/ndd-runtime/pkg/resource"
)

// A Initializer establishes ownership of the supplied Managed resource.
// This typically involves the operations that are run before calling any
// ExternalClient methods.
type Initializer interface {
	Initialize(ctx context.Context, mg resource.Managed) error
}

// A InitializerChain chains multiple managed initializers.
type InitializerChain []Initializer

// Initialize calls each Initializer serially. It returns the first
// error it encounters, if any.
func (cc InitializerChain) Initialize(ctx context.Context, mg resource.Managed) error {
	for _, c := range cc {
		if err := c.Initialize(ctx, mg); err != nil {
			return err
		}
	}
	return nil
}

// A InitializerFn is a function that satisfies the Initializer
// interface.
type InitializerFn func(ctx context.Context, mg resource.Managed) error

// Initialize calls InitializerFn function.
func (m InitializerFn) Initialize(ctx context.Context, mg resource.Managed) error {
	return m(ctx, mg)
}
