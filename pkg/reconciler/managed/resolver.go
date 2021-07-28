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

// A ReferenceResolver resolves references to other managed resources.
type ReferenceResolver interface {
	// ResolveReferences resolves all fields in the supplied managed resource
	// that are references to other managed resources by updating corresponding
	// fields, for example setting spec.network to the Network resource
	// specified by spec.networkRef.name.
	ResolveReferences(ctx context.Context, mg resource.Managed) error
}

// A ReferenceResolverFn is a function that satisfies the
// ReferenceResolver interface.
type ReferenceResolverFn func(context.Context, resource.Managed) error

// ResolveReferences calls ReferenceResolverFn function
func (m ReferenceResolverFn) ResolveReferences(ctx context.Context, mg resource.Managed) error {
	return m(ctx, mg)
}
