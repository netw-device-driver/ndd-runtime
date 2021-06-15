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

package meta

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	nddv1 "github.com/netw-device-driver/ndd-runtime/apis/common/v1"
	"github.com/pkg/errors"
)

// TypedReferenceTo returns a typed object reference to the supplied object,
// presumed to be of the supplied group, version, and kind.
func TypedReferenceTo(o metav1.Object, of schema.GroupVersionKind) *nddv1.TypedReference {
	v, k := of.ToAPIVersionAndKind()
	return &nddv1.TypedReference{
		APIVersion: v,
		Kind:       k,
		Name:       o.GetName(),
		UID:        o.GetUID(),
	}
}

// AsOwner converts the supplied object reference to an owner reference.
func AsOwner(r *nddv1.TypedReference) metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: r.APIVersion,
		Kind:       r.Kind,
		Name:       r.Name,
		UID:        r.UID,
	}
}

// AsController converts the supplied object reference to a controller
// reference. You may also consider using metav1.NewControllerRef.
func AsController(r *nddv1.TypedReference) metav1.OwnerReference {
	c := true
	ref := AsOwner(r)
	ref.Controller = &c
	return ref
}

// AddOwnerReference to the supplied object' metadata. Any existing owner with
// the same UID as the supplied reference will be replaced.
func AddOwnerReference(o metav1.Object, r metav1.OwnerReference) {
	refs := o.GetOwnerReferences()
	for i := range refs {
		if refs[i].UID == r.UID {
			refs[i] = r
			o.SetOwnerReferences(refs)
			return
		}
	}
	o.SetOwnerReferences(append(refs, r))
}

// AddControllerReference to the supplied object's metadata. Any existing owner
// with the same UID as the supplied reference will be replaced. Returns an
// error if the supplied object is already controlled by a different owner.
func AddControllerReference(o metav1.Object, r metav1.OwnerReference) error {
	if c := metav1.GetControllerOf(o); c != nil && c.UID != r.UID {
		return errors.Errorf("%s is already controlled by %s %s (UID %s)", o.GetName(), c.Kind, c.Name, c.UID)
	}

	AddOwnerReference(o, r)
	return nil
}

// AddFinalizer to the supplied Kubernetes object's metadata.
func AddFinalizer(o metav1.Object, finalizer string) {
	f := o.GetFinalizers()
	for _, e := range f {
		if e == finalizer {
			return
		}
	}
	o.SetFinalizers(append(f, finalizer))
}

// RemoveFinalizer from the supplied Kubernetes object's metadata.
func RemoveFinalizer(o metav1.Object, finalizer string) {
	f := o.GetFinalizers()
	for i, e := range f {
		if e == finalizer {
			f = append(f[:i], f[i+1:]...)
		}
	}
	o.SetFinalizers(f)
}

// FinalizerExists checks whether given finalizer is already set.
func FinalizerExists(o metav1.Object, finalizer string) bool {
	f := o.GetFinalizers()
	for _, e := range f {
		if e == finalizer {
			return true
		}
	}
	return false
}
