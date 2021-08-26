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
	"github.com/yndd/ndd-yang/pkg/parser"
)

type Validator interface {
	ValidateLocalleafRef(ctx context.Context, mg resource.Managed) (ValidateLocalleafRefObservation, error)

	ValidateExternalleafRef(ctx context.Context, mg resource.Managed, cfg []byte) (ValidateExternalleafRefObservation, error)

	ValidateParentDependency(ctx context.Context, mg resource.Managed, cfg []byte) (ValidationParentDependencyObservation, error)
}

type ValidatorFn struct {
	ValidateLocalleafRefFn     func(ctx context.Context, mg resource.Managed) (ValidateLocalleafRefObservation, error)
	ValidateExternalleafRefFn  func(ctx context.Context, mg resource.Managed, cfg []byte) (ValidateExternalleafRefObservation, error)
	ValidateParentDependencyFn func(ctx context.Context, mg resource.Managed, cfg []byte) (ValidationParentDependencyObservation, error)
}

func (e ValidatorFn) ValidateLocalleafRef(ctx context.Context, mg resource.Managed) (ValidateLocalleafRefObservation, error) {
	return e.ValidateLocalleafRefFn(ctx, mg)
}

func (e ValidatorFn) ValidateExternalleafRef(ctx context.Context, mg resource.Managed, cfg []byte) (ValidateExternalleafRefObservation, error) {
	return e.ValidateExternalleafRefFn(ctx, mg, cfg)
}

func (e ValidatorFn) ValidateParentDependency(ctx context.Context, mg resource.Managed, cfg []byte) (ValidationParentDependencyObservation, error) {
	return e.ValidateParentDependencyFn(ctx, mg, cfg)
}

type NopValidator struct{}

func (e *NopValidator) ValidateLocalleafRef(ctx context.Context, mg resource.Managed) (ValidateLocalleafRefObservation, error) {
	return ValidateLocalleafRefObservation{}, nil
}

func (e *NopValidator) ValidateExternalleafRef(ctx context.Context, mg resource.Managed, cfg []byte) (ValidateExternalleafRefObservation, error) {
	return ValidateExternalleafRefObservation{}, nil
}

func (e *NopValidator) ValidateParentDependency(ctx context.Context, mg resource.Managed, cfg []byte) (ValidationParentDependencyObservation, error) {
	return ValidationParentDependencyObservation{}, nil
}

type ValidateLocalleafRefObservation struct {
	Success bool

	ResolvedLeafRefs []*parser.ResolvedLeafRef

	Details []*parser.ResultleafRefValidation
}

type ValidateExternalleafRefObservation struct {
	Success bool

	ResolvedLeafRefs []*parser.ResolvedLeafRef

	Details []*parser.ResultleafRefValidation
}

type ValidationParentDependencyObservation struct {
	Success bool

	ResolvedLeafRefs []*parser.ResolvedLeafRef

	Details []*parser.ResultleafRefValidation
}
