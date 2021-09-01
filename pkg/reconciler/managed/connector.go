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
)

// ConnectionDetails created or updated during an operation on an external
// resource, for example usernames, passwords, endpoints, ports, etc.
type ConnectionDetails map[string][]byte

// An ExternalConnecter produces a new ExternalClient given the supplied
// Managed resource.
type ExternalConnecter interface {
	// Connect to the provider specified by the supplied managed resource and
	// produce an ExternalClient.
	Connect(ctx context.Context, mg resource.Managed) (ExternalClient, error)
}

// An ExternalConnectorFn is a function that satisfies the ExternalConnecter
// interface.
type ExternalConnectorFn func(ctx context.Context, mg resource.Managed) (ExternalClient, error)

// Connect to the provider specified by the supplied managed resource and
// produce an ExternalClient.
func (ec ExternalConnectorFn) Connect(ctx context.Context, mg resource.Managed) (ExternalClient, error) {
	return ec(ctx, mg)
}

// An ExternalClient manages the lifecycle of an external resource.
// None of the calls here should be blocking. All of the calls should be
// idempotent. For example, Create call should not return AlreadyExists error
// if it's called again with the same parameters or Delete call should not
// return error if there is an ongoing deletion or resource does not exist.
type ExternalClient interface {
	// Observe the external resource the supplied Managed resource represents,
	// if any. Observe implementations must not modify the external resource,
	// but may update the supplied Managed resource to reflect the state of the
	// external resource.
	Observe(ctx context.Context, mg resource.Managed) (ExternalObservation, error)

	// Create an external resource per the specifications of the supplied
	// Managed resource. Called when Observe reports that the associated
	// external resource does not exist.
	Create(ctx context.Context, mg resource.Managed) (ExternalCreation, error)

	// Update the external resource represented by the supplied Managed
	// resource, if necessary. Called unless Observe reports that the
	// associated external resource is up to date.
	Update(ctx context.Context, mg resource.Managed, obs ExternalObservation) (ExternalUpdate, error)

	// Delete the external resource upon deletion of its associated Managed
	// resource. Called when the managed resource has been deleted.
	Delete(ctx context.Context, mg resource.Managed) error

	// GetTarget returns the targets the resource is assigned assigned to
	GetTarget() []string

	// GetConfig returns the full configuration of the network node
	GetConfig(ctx context.Context) ([]byte, error)

	// GetResourceName returns the resource that matches the path
	GetResourceName(ctx context.Context, path []*gnmi.Path) (string, error)
}

// ExternalClientFns are a series of functions that satisfy the ExternalClient
// interface.
type ExternalClientFns struct {
	ObserveFn         func(ctx context.Context, mg resource.Managed) (ExternalObservation, error)
	CreateFn          func(ctx context.Context, mg resource.Managed) (ExternalCreation, error)
	UpdateFn          func(ctx context.Context, mg resource.Managed) (ExternalUpdate, error)
	DeleteFn          func(ctx context.Context, mg resource.Managed) error
	GetTargetFn       func() []string
	GetConfigFn       func(ctx context.Context) ([]byte, error)
	GetResourceNameFn func(ctx context.Context, path []*gnmi.Path) (string, error)
}

// Observe the external resource the supplied Managed resource represents, if
// any.
func (e ExternalClientFns) Observe(ctx context.Context, mg resource.Managed) (ExternalObservation, error) {
	return e.ObserveFn(ctx, mg)
}

// Create an external resource per the specifications of the supplied Managed
// resource.
func (e ExternalClientFns) Create(ctx context.Context, mg resource.Managed) (ExternalCreation, error) {
	return e.CreateFn(ctx, mg)
}

// Update the external resource represented by the supplied Managed resource, if
// necessary.
func (e ExternalClientFns) Update(ctx context.Context, mg resource.Managed) (ExternalUpdate, error) {
	return e.UpdateFn(ctx, mg)
}

// Delete the external resource upon deletion of its associated Managed
// resource.
func (e ExternalClientFns) Delete(ctx context.Context, mg resource.Managed) error {
	return e.DeleteFn(ctx, mg)
}

// GetTarget return the real target for the external resource
func (e ExternalClientFns) GetTarget() []string {
	return e.GetTargetFn()
}

// GetConfig returns the full configuration of the network node
func (e ExternalClientFns) GetConfig(ctx context.Context) ([]byte, error) {
	return e.GetConfigFn(ctx)
}

// GetResourceName returns the resource matching the path
func (e ExternalClientFns) GetResourceName(ctx context.Context, path []*gnmi.Path) (string, error) {
	return e.GetResourceNameFn(ctx, path)
}

// A NopConnecter does nothing.
type NopConnecter struct{}

// Connect returns a NopClient. It never returns an error.
func (c *NopConnecter) Connect(_ context.Context, _ resource.Managed) (ExternalClient, error) {
	return &NopClient{}, nil
}

// A NopClient does nothing.
type NopClient struct{}

// Observe does nothing. It returns an empty ExternalObservation and no error.
func (c *NopClient) Observe(ctx context.Context, mg resource.Managed) (ExternalObservation, error) {
	return ExternalObservation{}, nil
}

// Create does nothing. It returns an empty ExternalCreation and no error.
func (c *NopClient) Create(ctx context.Context, mg resource.Managed) (ExternalCreation, error) {
	return ExternalCreation{}, nil
}

// Update does nothing. It returns an empty ExternalUpdate and no error.
func (c *NopClient) Update(ctx context.Context, mg resource.Managed, obs ExternalObservation) (ExternalUpdate, error) {
	return ExternalUpdate{}, nil
}

// Delete does nothing. It never returns an error.
func (c *NopClient) Delete(ctx context.Context, mg resource.Managed) error { return nil }

// GetTarget return on empty string list
func (c *NopClient) GetTarget() []string { return make([]string, 0) }

// GetConfig returns the full configuration of the network node
func (c *NopClient) GetConfig(ctx context.Context) ([]byte, error) {
	return make([]byte, 0), nil
}

// GetResourceName returns the resource matching the path
func (c *NopClient) GetResourceName(ctx context.Context, path []*gnmi.Path) (string, error) {
	return "", nil
}

// An ExternalObservation is the result of an observation of an external
// resource.
type ExternalObservation struct {
	// ResourceExists must be true if a corresponding external resource exists
	// for the managed resource. This is checked using the resourceKey in the
	// device driver
	ResourceExists bool

	// ResourceHasData can be true when a managed resource is created, but the
	// device had already data in that resource. The data needs to get aligned
	// with the intended resource data
	ResourceHasData bool

	// ResourceUpToDate should be true if the corresponding external resource
	// appears to be up-to-date
	ResourceUpToDate bool

	ResourceDeletes []*gnmi.Path
	ResourceUpdates []*gnmi.Update
}

// An ExternalCreation is the result of the creation of an external resource.

type ExternalCreation struct {
}

// An ExternalUpdate is the result of an update to an external resource.
type ExternalUpdate struct {
}
