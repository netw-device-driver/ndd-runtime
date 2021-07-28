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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type adder interface {
	Add(item interface{})
}

// EnqueueRequestForNetworkNode enqueues a reconcile.Request for a referenced
// NetworkNode.
type EnqueueRequestForNetworkNode struct{}

// Create adds a NamespacedName for the supplied CreateEvent if its Object is a
// NetworkNodeReferencer.
func (e *EnqueueRequestForNetworkNode) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	addNetworkNode(evt.Object, q)
}

// Update adds a NamespacedName for the supplied UpdateEvent if its Objects are
// a NetworkNodeReferencer.
func (e *EnqueueRequestForNetworkNode) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	addNetworkNode(evt.ObjectOld, q)
	addNetworkNode(evt.ObjectNew, q)
}

// Delete adds a NamespacedName for the supplied DeleteEvent if its Object is a
// NetworkNodeReferencer.
func (e *EnqueueRequestForNetworkNode) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	addNetworkNode(evt.Object, q)
}

// Generic adds a NamespacedName for the supplied GenericEvent if its Object is
// a NetworkNodeReferencer.
func (e *EnqueueRequestForNetworkNode) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {
	addNetworkNode(evt.Object, q)
}

func addNetworkNode(obj runtime.Object, queue adder) {
	pcr, ok := obj.(RequiredNetworkNodeReferencer)
	if !ok {
		return
	}

	queue.Add(reconcile.Request{NamespacedName: types.NamespacedName{Name: pcr.GetNetworkNodeReference().Name}})
}
