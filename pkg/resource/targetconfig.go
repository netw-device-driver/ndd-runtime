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

	nddv1 "github.com/netw-device-driver/ndd-runtime/apis/common/v1"
	"github.com/netw-device-driver/ndd-runtime/pkg/meta"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	errMissingTCRef = "managed resource does not reference a TargetConfig"
	errApplyTCU     = "cannot apply TargetConfigUsage"
)

type errMissingRef struct{ error }

func (m errMissingRef) MissingReference() bool { return true }

// A Tracker tracks managed resources.
type Tracker interface {
	// Track the supplied managed resource.
	Track(ctx context.Context, mg Managed) error
}

// A TrackerFn is a function that tracks managed resources.
type TrackerFn func(ctx context.Context, mg Managed) error

// Track the supplied managed resource.
func (fn TrackerFn) Track(ctx context.Context, mg Managed) error {
	return fn(ctx, mg)
}

// A TargetConfigUsageTracker tracks usages of a TargetConfig by creating or
// updating the appropriate TargetConfigUsage.
type TargetConfigUsageTracker struct {
	c  Applicator
	of TargetConfigUsage
}

// NewTargetConfigUsageTracker creates a TargetConfigUsageTracker.
func NewTargetConfigUsageTracker(c client.Client, of TargetConfigUsage) *TargetConfigUsageTracker {
	return &TargetConfigUsageTracker{c: NewAPIUpdatingApplicator(c), of: of}
}

// Track that the supplied Managed resource is using the TargetConfig it
// references by creating or updating a TargetConfigUsage. Track should be
// called _before_ attempting to use the TargetConfig. This ensures the
// managed resource's usage is updated if the managed resource is updated to
// reference a misconfigured TargetConfig.
func (u *TargetConfigUsageTracker) Track(ctx context.Context, mg Managed) error {
	pcu := u.of.DeepCopyObject().(TargetConfigUsage)
	gvk := mg.GetObjectKind().GroupVersionKind()
	ref := mg.GetTargetConfigReference()
	if ref == nil {
		return errMissingRef{errors.New(errMissingTCRef)}
	}

	pcu.SetName(string(mg.GetUID()))
	pcu.SetLabels(map[string]string{nddv1.LabelKeyTargetName: ref.Name})
	pcu.SetOwnerReferences([]metav1.OwnerReference{meta.AsController(meta.TypedReferenceTo(mg, gvk))})
	pcu.SetTargetConfigReference(nddv1.Reference{Name: ref.Name})
	pcu.SetResourceReference(nddv1.TypedReference{
		APIVersion: gvk.GroupVersion().String(),
		Kind:       gvk.Kind,
		Name:       mg.GetName(),
	})

	err := u.c.Apply(ctx, pcu,
		MustBeControllableBy(mg.GetUID()),
		AllowUpdateIf(func(current, _ runtime.Object) bool {
			return current.(TargetConfigUsage).GetTargetConfigReference() != pcu.GetTargetConfigReference()
		}),
	)
	return errors.Wrap(Ignore(IsNotAllowed, err), errApplyTCU)
}
