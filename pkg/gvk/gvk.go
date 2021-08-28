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

package gvk

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

const (
	errJSONMarshal   = "cannot marshal JSON object"
	errJSONUnMarshal = "cannot unmarshal JSON object"
)

// GVK identifies a resource with Group - Version - Kind - Name - NameSpace
type GVK struct {
	Group     string `json:"group,omitempty"`
	Version   string `json:"version,omitempty"`
	Kind      string `json:"kind,omitempty"`
	Name      string `json:"Name,omitempty"`
	NameSpace string `json:"Namespace,omitempty"`
}

func (gvk *GVK) String() (string, error) {
	d, err := json.Marshal(&gvk)
	if err != nil {
		return "", errors.Wrap(err, errJSONMarshal)
	}
	return string(d), nil
}

func (gvk *GVK) GetGroup() string {
	return gvk.Group
}

func (gvk *GVK) GetVersion() string {
	return gvk.Version
}

func (gvk *GVK) GetKind() string {
	return gvk.Kind
}

func (gvk *GVK) GetName() string {
	return gvk.Name
}

func (gvk *GVK) GetNameSpace() string {
	return gvk.NameSpace
}

func (gvk *GVK) SetGroup(s string) {
	gvk.Group = s
}

func (gvk *GVK) SetVersion(s string) {
	gvk.Version = s
}

func (gvk *GVK) SetKind(s string) {
	gvk.Kind = s
}

func (gvk *GVK) SetName(s string) {
	gvk.Name = s
}

func (gvk *GVK) SetNameSpace(s string) {
	gvk.NameSpace = s
}

func String2GVK(s string) (*GVK, error) {
	var gvk GVK
	fmt.Printf("String2GVK: %s\n", s)
	if err := json.Unmarshal([]byte(s), &gvk); err != nil {
		return nil, errors.Wrap(err, errJSONUnMarshal)
	}
	fmt.Printf("GVK: %v\n", gvk)
	return &gvk, nil
}
