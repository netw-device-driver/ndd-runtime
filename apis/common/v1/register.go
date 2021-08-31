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

package v1

import "github.com/openconfig/gnmi/proto/gnmi"

var RegisterPath = &gnmi.Path{
	Elem: []*gnmi.PathElem{
		{Name: RegisterPathElemName, Key: map[string]string{RegisterPathElemKey: ""}},
	},
}

var RegisterPathElemName = "ndd-registration"
var RegisterPathElemKey = "name"

// RegistrationParameters defines the Registrations the device driver subscribes to for config change notifications
type Register struct {
	// Registrations defines the Registrations the device driver subscribes to for config change notifications
	// +optional
	Subscriptions []string `json:"subscriptions,omitempty"`

	// ExceptionPaths defines the exception paths that should be ignored during change notifications
	// if the xpath contains the exception path it is considered a match
	// +optional
	ExceptionPaths []string `json:"exceptionPaths,omitempty"`

	// ExplicitExceptionPaths defines the exception paths that should be ignored during change notifications
	// the match should be exact to condider this xpath
	// +optional
	ExplicitExceptionPaths []string `json:"explicitExceptionPaths,omitempty"`
}

func (r *Register) GetSubscriptions() []string {
	return r.Subscriptions
}

func (r *Register) SetSubscriptions(s []string) {
	r.Subscriptions = s
}
