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

import "context"


// the provider finalizer manages the finalizer for the external leafrefs accross all 
// all managed resources of the provider
type ProviderFinalizer interface {
	AddFinalizer(ctx context.Context, name, namespace string) error

	RemoveFinalizer(ctx context.Context, name, namespace string) error
}

type ProviderFinalizerFn struct {
	AddFinalizerFn func(ctx context.Context, name, namespace string) error

	RemoveFinalizerFn func(ctx context.Context, name, namespace string) error
}

func (e ProviderFinalizerFn) AddFinalizer(ctx context.Context, name, namespace string) error {
	return e.AddFinalizerFn(ctx, name, namespace)
}

func (e ProviderFinalizerFn) RemoveFinalizer(ctx context.Context, name, namespace string) error {
	return e.AddFinalizerFn(ctx, name, namespace)
}

type NopProviderFinalizer struct{}

func (e *NopProviderFinalizer) AddFinalizer(ctx context.Context, name, namespace string) error {
	return nil
}

func (e *NopProviderFinalizer) RemoveFinalizer(ctx context.Context, name, namespace string) error {
	return nil
}
