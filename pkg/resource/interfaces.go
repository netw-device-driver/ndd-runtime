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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	nddv1 "github.com/netw-device-driver/ndd-runtime/apis/common/v1"
)

// A Conditioned may have conditions set or retrieved. Conditions are typically
// indicate the status of both a resource and its reconciliation process.
type Conditioned interface {
	SetConditions(c ...nddv1.Condition)
	GetCondition(ck nddv1.ConditionKind) nddv1.Condition
}

// A TargetConditioned may have conditions set or retrieved. TargetConditioned are typically
// indicate the status of both a resource and its reconciliation process.
type TargetConditioned interface {
	SetTargetConditions(target string, c ...nddv1.Condition)
	GetTargetCondition(target string, ck nddv1.ConditionKind) nddv1.Condition
}

// A TargetConfigReferencer may reference a Target config resource.
type TargetConfigReferencer interface {
	GetTargetConfigReference() *nddv1.Reference
	SetTargetConfigReference(p *nddv1.Reference)
}

// A RequiredTargetConfigReferencer may reference a Target config resource.
// Unlike TargetConfigReferencer, the reference is required (i.e. not nil).
type RequiredTargetConfigReferencer interface {
	GetTargetConfigReference() nddv1.Reference
	SetTargetConfigReference(p nddv1.Reference)
}

// A RequiredTypedResourceReferencer can reference a resource.
type RequiredTypedResourceReferencer interface {
	SetResourceReference(r nddv1.TypedReference)
	GetResourceReference() nddv1.TypedReference
}

// An Orphanable resource may specify a DeletionPolicy.
type Orphanable interface {
	SetDeletionPolicy(p nddv1.DeletionPolicy)
	GetDeletionPolicy() nddv1.DeletionPolicy
}

// An Orphanable resource may specify a DeletionPolicy.
type Active interface {
	SetActive(p bool)
	GetActive() bool
}

// A UserCounter can count how many users it has.
type UserCounter interface {
	SetUsers(i int64)
	GetUsers() int64
}

// An Object is a Kubernetes object.
type Object interface {
	metav1.Object
	runtime.Object
}

// A TargetConfig configures a Network Device Driver Target.
type TargetConfig interface {
	Object

	UserCounter
	Conditioned
}

// A TargetConfigUsage indicates a usage of a ndd target config.
type TargetConfigUsage interface {
	Object

	RequiredTargetConfigReferencer
	RequiredTypedResourceReferencer
}

// A TargetConfigUsageList is a list of Target config usages.
type TargetConfigUsageList interface {
	client.ObjectList

	// GetItems returns the list of Target config usages.
	GetItems() []TargetConfigUsage
}

// A Finalizer manages the finalizers on the resource.
type Finalizer interface {
	AddFinalizer(ctx context.Context, obj Object) error
	RemoveFinalizer(ctx context.Context, obj Object) error
}

// A Managed is a Kubernetes object representing a concrete managed
// resource (e.g. a CloudSQL instance).
type Managed interface {
	Object

	TargetConfigReferencer
	Orphanable
	Active

	Conditioned
}

// A ManagedList is a list of managed resources.
type ManagedList interface {
	client.ObjectList

	// GetItems returns the list of managed resources.
	GetItems() []Managed
}
