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
	errMissingNNRef = "managed resource does not reference a NetworkNode"
	errApplyNNU     = "cannot apply NetworkNodeUsage"
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

// A NetworkNodeUsageTracker tracks usages of a NetworkNode by creating or
// updating the appropriate NetworkNodeUsage.
type NetworkNodeUsageTracker struct {
	c  Applicator
	of NetworkNodeUsage
}

// NewNetworkNodeUsageTracker creates a NetworkNodeUsageTracker.
func NewNetworkNodeUsageTracker(c client.Client, of NetworkNodeUsage) *NetworkNodeUsageTracker {
	return &NetworkNodeUsageTracker{c: NewAPIUpdatingApplicator(c), of: of}
}

// Track that the supplied Managed resource is using the NetworkNode it
// references by creating or updating a NetworkNodeUsage. Track should be
// called _before_ attempting to use the NetworkNode. This ensures the
// managed resource's usage is updated if the managed resource is updated to
// reference a misconfigured NetworkNode.
func (u *NetworkNodeUsageTracker) Track(ctx context.Context, mg Managed) error {
	pcu := u.of.DeepCopyObject().(NetworkNodeUsage)
	gvk := mg.GetObjectKind().GroupVersionKind()
	ref := mg.GetNetworkNodeReference()
	if ref == nil {
		return errMissingRef{errors.New(errMissingNNRef)}
	}

	pcu.SetName(string(mg.GetUID()))
	pcu.SetLabels(map[string]string{nddv1.LabelKeyNetworkNodeName: ref.Name})
	pcu.SetOwnerReferences([]metav1.OwnerReference{meta.AsController(meta.TypedReferenceTo(mg, gvk))})
	pcu.SetNetworkNodeReference(nddv1.Reference{Name: ref.Name})
	pcu.SetResourceReference(nddv1.TypedReference{
		APIVersion: gvk.GroupVersion().String(),
		Kind:       gvk.Kind,
		Name:       mg.GetName(),
	})

	err := u.c.Apply(ctx, pcu,
		MustBeControllableBy(mg.GetUID()),
		AllowUpdateIf(func(current, _ runtime.Object) bool {
			return current.(NetworkNodeUsage).GetNetworkNodeReference() != pcu.GetNetworkNodeReference()
		}),
	)
	return errors.Wrap(Ignore(IsNotAllowed, err), errApplyNNU)
}
