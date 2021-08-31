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

/*
// A Resolver resolves references to other managed resources.
type Resolver interface {
	GetManagedResource(ctx context.Context, resourceName string) (resource.Managed, error)
}

type ResolverFn struct {
	// A GetManagedResourceFn is a function that satisfies the GetManagedResource interface.
	GetManagedResourceFn func(ctx context.Context, resourceName string) (resource.Managed, error)
}

// GetManagedResource calls ReferenceResolverFn function
func (r ResolverFn) GetManagedResource(ctx context.Context, resourceName string) (resource.Managed, error) {
	return r.GetManagedResourceFn(ctx, resourceName)
}

type NopResolver struct{}

func (e *NopResolver) GetManagedResource(ctx context.Context, resourceName string) (resource.Managed, error) {
	return nil, nil
}
*/
