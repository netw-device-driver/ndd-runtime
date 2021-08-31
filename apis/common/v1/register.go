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

// Registrations defines the Registrations the device driver subscribes to for config change notifications
type Register struct {
	DeviceType    DeviceType `json:"deviceType,omitempty"`
	DeviceMatch   string     `json:"deviceMatch,omitempty"`
	Subscriptions []string   `json:"subscriptions,omitempty"`
}

func (r *Register) GetDeviceType() DeviceType {
	return r.DeviceType
}

func (r *Register) GetDeviceMatch() string {
	return r.DeviceMatch
}

func (r *Register) GetSubscriptions() []string {
	return r.Subscriptions
}

func (r *Register) SetDeviceType(s DeviceType) {
	r.DeviceType = s
}

func (r *Register) SetDeviceMatch(s string) {
	r.DeviceMatch = s
}

func (r *Register) SetSubscriptions(s []string) {
	r.Subscriptions = s
}
