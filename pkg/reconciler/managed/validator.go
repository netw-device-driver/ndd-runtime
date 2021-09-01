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
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/yndd/ndd-yang/pkg/parser"
)

type Validator interface {
	ValidateLocalleafRef(ctx context.Context, mg resource.Managed) (ValidateLocalleafRefObservation, error)

	ValidateExternalleafRef(ctx context.Context, mg resource.Managed, cfg []byte) (ValidateExternalleafRefObservation, error)

	ValidateParentDependency(ctx context.Context, mg resource.Managed, cfg []byte) (ValidateParentDependencyObservation, error)

	ValidateResourceIndexes(ctx context.Context, mg resource.Managed) (ValidateResourceIndexesObservation, error)
}

type ValidatorFn struct {
	ValidateLocalleafRefFn     func(ctx context.Context, mg resource.Managed) (ValidateLocalleafRefObservation, error)
	ValidateExternalleafRefFn  func(ctx context.Context, mg resource.Managed, cfg []byte) (ValidateExternalleafRefObservation, error)
	ValidateParentDependencyFn func(ctx context.Context, mg resource.Managed, cfg []byte) (ValidateParentDependencyObservation, error)
	ValidateResourceIndexesFn  func(ctx context.Context, mg resource.Managed) (ValidateResourceIndexesObservation, error)
}

func (e ValidatorFn) ValidateLocalleafRef(ctx context.Context, mg resource.Managed) (ValidateLocalleafRefObservation, error) {
	return e.ValidateLocalleafRefFn(ctx, mg)
}

func (e ValidatorFn) ValidateExternalleafRef(ctx context.Context, mg resource.Managed, cfg []byte) (ValidateExternalleafRefObservation, error) {
	return e.ValidateExternalleafRefFn(ctx, mg, cfg)
}

func (e ValidatorFn) ValidateParentDependency(ctx context.Context, mg resource.Managed, cfg []byte) (ValidateParentDependencyObservation, error) {
	return e.ValidateParentDependencyFn(ctx, mg, cfg)
}

func (e ValidatorFn) ValidateResourceIndexes(ctx context.Context, mg resource.Managed) (ValidateResourceIndexesObservation, error) {
	return e.ValidateResourceIndexesFn(ctx, mg)
}

type NopValidator struct{}

func (e *NopValidator) ValidateLocalleafRef(ctx context.Context, mg resource.Managed) (ValidateLocalleafRefObservation, error) {
	return ValidateLocalleafRefObservation{}, nil
}

func (e *NopValidator) ValidateExternalleafRef(ctx context.Context, mg resource.Managed, cfg []byte) (ValidateExternalleafRefObservation, error) {
	return ValidateExternalleafRefObservation{}, nil
}

func (e *NopValidator) ValidateParentDependency(ctx context.Context, mg resource.Managed, cfg []byte) (ValidateParentDependencyObservation, error) {
	return ValidateParentDependencyObservation{}, nil
}

func (e *NopValidator) ValidateResourceIndexes(ctx context.Context, mg resource.Managed) (ValidateResourceIndexesObservation, error) {
	return ValidateResourceIndexesObservation{}, nil
}

type ValidateLocalleafRefObservation struct {
	Success bool

	ResolvedLeafRefs []*parser.ResolvedLeafRefGnmi
}

type ValidateExternalleafRefObservation struct {
	Success bool

	ResolvedLeafRefs []*parser.ResolvedLeafRefGnmi
}

type ValidateParentDependencyObservation struct {
	Success bool

	ResolvedLeafRefs []*parser.ResolvedLeafRefGnmi
}

type ValidateResourceIndexesObservation struct {
	Changed         bool
	ResourceDeletes []*gnmi.Path
	ResourceIndexes map[string]string
}
