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

type ResourceIndexes interface {
	SetResourceIndexes(resourceIndexes map[string]string)
	GetResourceIndexes() map[string]string
}

type ExternalLeafRefs interface {
	SetExternalLeafRefs(externalResourceNames []string)
	GetExternalLeafRefs() []string
}

// A Target may have targets
type Target interface {
	SetTarget(target []string)
	GetTarget() []string
}

// A NetworkNodeReferencer may reference a Network Node resource.
type NetworkNodeReferencer interface {
	GetNetworkNodeReference() *nddv1.Reference
	SetNetworkNodeReference(p *nddv1.Reference)
}

// A RequiredNetworkNodeReferencer may reference a network node config resource.
// Unlike NetworkNodeReferencer, the reference is required (i.e. not nil).
type RequiredNetworkNodeReferencer interface {
	GetNetworkNodeReference() nddv1.Reference
	SetNetworkNodeReference(p nddv1.Reference)
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

// A Managed is a Kubernetes object representing a concrete managed
// resource (e.g. a CloudSQL instance).
type Managed interface {
	Object

	NetworkNodeReferencer
	Orphanable
	Active

	Conditioned
	Target
	ExternalLeafRefs
	ResourceIndexes
}

// A ManagedList is a list of managed resources.
type ManagedList interface {
	client.ObjectList

	// GetItems returns the list of managed resources.
	GetItems() []Managed
}

// A NetworkNode configures a NetworkNode Device Driver.
type NetworkNode interface {
	Object

	UserCounter
	Conditioned
}

// A NetworkNodeUsage indicates a usage of a network node config.
type NetworkNodeUsage interface {
	Object

	RequiredNetworkNodeReferencer
	RequiredTypedResourceReferencer
}

// A NetworkNodeUsageList is a list of network node usages.
type NetworkNodeUsageList interface {
	client.ObjectList

	// GetItems returns the list of network node usages.
	GetItems() []NetworkNodeUsage
}

// A Finalizer manages the finalizers on the resource.
type Finalizer interface {
	AddFinalizer(ctx context.Context, obj Object) error
	RemoveFinalizer(ctx context.Context, obj Object) error
	HasOtherFinalizer(ctx context.Context, obj Object) (bool, error)
	AddFinalizerString(ctx context.Context, obj Object, finalizerString string) error
	RemoveFinalizerString(ctx context.Context, obj Object, finalizerString string) error
}
